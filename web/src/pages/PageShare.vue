<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAxios } from '@vueuse/integrations/useAxios'
import type { AxiosError } from 'axios'

import { useAxiosInstance } from '../lib/axios.ts'
import { useStore } from '../lib/store.ts'
import { isoToRelative } from '../utils/datetime.ts'
import { formatSize } from '../utils/filesize.ts'
import { useSharePasswords } from '../utils/share-passwords.ts'
import { getContentUrl } from '../utils/share-url.ts'
import FlowbiteSpinner from '../components/FlowbiteSpinner.vue'
import ShareIcon from '../components/ShareIcon.vue'
import LayoutDashboard from '../layouts/LayoutDashboard.vue'
import type { Share, ShareFile } from '../types/share.ts'

import BiBoxArrowInUpRight from 'bootstrap-icons/icons/box-arrow-in-up-right.svg?component'
import BiCopy from 'bootstrap-icons/icons/copy.svg?component'
import BiDot from 'bootstrap-icons/icons/dot.svg?component'
import BiTrash3 from 'bootstrap-icons/icons/trash3.svg?component'
import ShareFileTree from '../components/ShareFileTree.vue'

const store = useStore()

const route = useRoute()

const router = useRouter()

const name = computed<string>(() => {
  if (typeof route.params.name === 'object') {
    return route.params.name[0]
  }
  return route.params.name
})

const passwords = useSharePasswords()

const error = ref('')

const locked = ref(false)

const {
  data: share,
  isLoading,
  execute: _loadData
} = useAxios<Share>('', useAxiosInstance(), {
  immediate: false,
  resetOnExecute: true,
  onSuccess: (data: Share) => {
    if (data.type === 'text' || data.type === 'url') {
      loadContent()
    }
  },
  onError: (e: AxiosError | any) => {
    if (e.response?.status === 404) {
      error.value = 'The share you are looking for is missing'
    } else if (e.response?.status === 401) {
      locked.value = true
    } else {
      console.error(e)
      error.value = 'Failed to load share content'
    }
  }
})

const loadData = () => {
  error.value = ''
  locked.value = false
  textContent.value = ''
  _loadData(`/api/shares/${name.value}`, {
    headers: {
      'X-Share-Password': passwords.value[name.value]
    }
  })
}

const textContent = ref('')

const loadContent = (fileId?: string) => {
  if (!share.value) {
    return
  }

  useAxiosInstance().get(getContentUrl(share.value, fileId))
    .then(({ data }) => {
      if (share.value?.type === 'url') {
        textContent.value = data
        if (!isOwner.value) {
          location.href = data
        }
      } else if (share.value?.type === 'text') {
        textContent.value = data
      }
    })
}

const displayName = computed(() => {
  if (share.value) {
    if (share.value.displayName) {
      return share.value.displayName
    }
    return share.value.name
  }
  return name.value
})

const isOwner = computed(() => {
  return store.loggedIn && store.userinfo?.subject === share.value?.creator?.subject
})

const onCopyLink = () => {
  const url = new URL(location.origin)
  url.pathname = `/s/${share.value?.name ?? name.value}`
  navigator.clipboard.writeText(url.toString())
  alert('Link to this share copied!')
}

const onCopyContent = () => {
  navigator.clipboard.writeText(textContent.value)
  alert('Content copied!')
}

const onRedirect = () => {
  location.href = textContent.value
}

const onDownloadFile = (file: ShareFile) => {
  if (!share.value) {
    return
  }

  const a = document.createElement('a')
  a.href = getContentUrl(share.value, file.id)
  a.target = '_blank'
  {
    const parts = file.path.split('/')
    if (parts.length) {
      a.download = parts[parts.length - 1]
    }
  }
  a.click()
}

const onDelete = () => {
  if (confirm('Are you sure you want to delete this share?')) {
    useAxiosInstance().delete(`/api/shares/${name.value}`)
      .then(() => {
        router.replace('/shares')
      })
      .catch((e) => {
        console.error(e)
        alert('Failed to delete share')
      })
  }
}

watch(name, () => {
  loadData()
}, {
  immediate: true
})
</script>

<template>
  <layout-dashboard>
    <div class="py-4 border-b border-gray-200 dark:border-neutral-800">
      <h5
        class="flex items-center gap-3 text-3xl font-name font-medium select-none cursor-pointer"
        @click="onCopyLink"
      >
        <flowbite-spinner
          v-if="isLoading"
          class="flex-shrink-0 w-6 h-6"
        />
        <share-icon
          v-else
          :share-type="share?.type"
          class="flex-shrink-0 w-6 h-6"
        />
        <span v-text="displayName" />
      </h5>
      <p class="flex items-center flex-wrap gap-y-1 mt-2 text-xs sm:text-sm text-gray-500 dark:text-neutral-400">
        <template v-if="isOwner">
          <button
            class="inline-flex items-center gap-1 hover:text-red-500 hover:underline underline-offset-2"
            @click="onDelete"
          >
            <bi-trash3 class="w-3 h-3" />
            <span>Delete</span>
          </button>
          <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
        </template>
        <template v-else-if="share?.creator?.username">
          <span>By <strong>{{ share.creator.username }}</strong></span>
          <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
        </template>
        <template v-if="share?.files && share.files.length > 1">
          <span>{{ share.files.length }} files</span>
          <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
        </template>
        <template v-if="share?.files && share?.size">
          <span v-text="formatSize(share.size)" />
          <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
        </template>
        <template v-if="share?.createdAt">
          <span>created {{ isoToRelative(share.createdAt) }}</span>
          <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
        </template>
        <template v-if="share?.expiresAt">
          <span>expires {{ isoToRelative(share.expiresAt) }}</span>
          <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
        </template>
        <template v-else-if="share?.expiry">
          <span>{{ share.expiry }} days remaining</span>
          <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
        </template>
        <template v-if="textContent">
          <button
            class="inline-flex items-center gap-1 hover:text-gray-600 dark:hover:text-neutral-300 hover:underline underline-offset-2"
            @click="onCopyContent"
          >
            <bi-copy class="w-3 h-3" />
            <span>Copy</span>
          </button>
          <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
          <template v-if="share?.type === 'url'">
            <a
              href="#"
              class="inline-flex items-center gap-1 hover:text-gray-600 dark:hover:text-neutral-300 hover:underline underline-offset-2"
              @click.prevent="onRedirect"
            >
              <bi-box-arrow-in-up-right class="w-3 h-3" />
              <span>Redirect</span>
            </a>
            <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
          </template>
        </template>
      </p>
    </div>
    <div
      v-if="locked || error"
      class="py-4 sm:py-32 flex justify-center"
    >
      <div class="w-full sm:w-auto bg-white dark:bg-neutral-800 flex flex-col items-stretch gap-3 p-4 sm:p-6 rounded-lg">
        <h1 class="text-2xl text-center font-semibold text-gray-800 dark:text-neutral-100 sm:px-4">
          {{ locked ? 'Password Required' : 'Oops!' }}
        </h1>
        <p
          v-if="error"
          class="text-sm text-gray-500 dark:text-neutral-400"
          v-text="error"
        />
        <div
          v-if="locked"
          class="flex items-center gap-3 mt-2"
        >
          <input
            v-model="passwords[name]"
            type="password"
            class="flex-1 bg-neutral-50 text-neutral-900 text-sm rounded-lg focus:ring-blue-500 py-1.5 px-2.5 dark:bg-neutral-700 dark:placeholder-neutral-400 dark:text-white dark:focus:ring-blue-500"
            placeholder="Password"
            autocomplete="off"
            @keydown.enter="loadData"
          />
          <button
            class="flex-shrink-0 inline-flex items-center gap-1 text-gray-200 bg-black hover:bg-neutral-900 focus:outline-none focus:ring-2 focus:ring-gray-400 text-sm font-medium rounded-lg px-4 py-1.5"
            @click="loadData"
          >
            Go
          </button>
        </div>
      </div>
    </div>
    <div
      v-else-if="!share || (share.type === 'text' && !textContent)"
      class="py-4 sm:py-32 flex items-center justify-center gap-2"
    >
      <flowbite-spinner class="w-6 h-6" />
      <p class="text-lg text-gray-500 dark:text-neutral-400">
        Loading content...
      </p>
    </div>
    <div
      v-else-if="share.type === 'text' || (share.type === 'url' && isOwner)"
      class="w-full py-4 text-gray-600 dark:text-neutral-300 text-xs font-mono whitespace-pre-wrap overflow-x-auto"
      v-text="textContent"
    />
    <div
      v-else-if="share.type === 'url'"
      class="py-4 sm:py-32 flex items-center justify-center gap-2"
    >
      <flowbite-spinner class="w-6 h-6" />
      <p class="text-lg text-gray-500 dark:text-neutral-400">
        Redirecting...
      </p>
    </div>
    <share-file-tree
      v-else-if="share.files"
      :files="share.files"
      @download="onDownloadFile"
    />
  </layout-dashboard>
</template>
