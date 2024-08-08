<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'

import { formatSize } from '../utils/filesize.ts'
import ShareFilePreview from './ShareFilePreview.vue'
import type { Share, ShareFile } from '../types/share.ts'

import BiCaretRightFill from 'bootstrap-icons/icons/caret-right-fill.svg?component'
import BiCloudDownload from 'bootstrap-icons/icons/cloud-download.svg?component'
import BiFileEarmarkText from 'bootstrap-icons/icons/file-earmark-text.svg?component'
import BiFolder2 from 'bootstrap-icons/icons/folder2.svg?component'
import BiFolder2Open from 'bootstrap-icons/icons/folder2-open.svg?component'

interface EntryDirectory {
  name: string
  size: number
}

interface EntryFile extends ShareFile {
  name: string
}

interface Entries {
  directories: EntryDirectory[]
  files: EntryFile[]
  isFile: boolean
}

const props = defineProps<{
  share: Share
  files: ShareFile[]
}>()

const emit = defineEmits<{
  download: [value: ShareFile]
}>()

// Note: There are two similar values to distinguish:
// - `isSingleFile`: Whether the share has only one file
// - `entries.isFile`: Whether the current path is a file
const isSingleFile = computed(() => {
  return props.files.length === 1
})

const currentPath = ref('')

const entries = computed<Entries>(() => {
  const directories = {} as Record<string, number>
  const files = [] as EntryFile[]
  for (const file of props.files) {
    if (file.path === currentPath.value) {
      return {
        directories: [],
        files: [Object.assign({}, file, {
          name: file.path
        })],
        isFile: true
      }
    }
    if (file.path.startsWith(currentPath.value)) {
      const subPath = file.path.slice(currentPath.value.length)
      const parts = subPath.split('/')
      if (parts.length > 1) {
        const directory = parts[0]
        if (!directories[directory]) {
          directories[directory] = 0
        }
        if (file.size) {
          directories[directory] += file.size
        }
      } else {
        files.push(Object.assign({}, file, {
          name: subPath
        }))
      }
    }
  }
  return {
    directories: Object.keys(directories).map(name => ({
      name,
      size: directories[name]
    })).sort((a, b) => a.name.localeCompare(b.name)),
    files: files.sort((a, b) => a.name.localeCompare(b.name)),
    isFile: false
  }
})

const breadcrumbsContainer = ref<HTMLDivElement | null>(null)

const breadcrumbs = computed<string[]>(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  if (!isSingleFile.value) {
    parts.unshift('Root')
  }
  return parts
})

const onBreadcrumbNavigate = (i: number) => {
  if (entries.value.isFile && i === breadcrumbs.value.length - 1) {
    return
  }

  if (i === 0) {
    currentPath.value = ''
  } else {
    currentPath.value = breadcrumbs.value.slice(1, i + 1).join('/') + '/'
  }
}

watch(() => props.files, () => {
  if (isSingleFile.value) {
    currentPath.value += props.files[0].path
  } else if (entries.value.directories.length === 1 && entries.value.files.length === 0) {
    currentPath.value += entries.value.directories[0].name + '/'
  }
}, {
  immediate: true
})

watch(breadcrumbs, () => {
  nextTick(() => {
    if (breadcrumbsContainer.value) {
      breadcrumbsContainer.value.scrollTo({
        left: breadcrumbsContainer.value.scrollWidth,
        behavior: 'smooth'
      })
    }
  })
}, {
  immediate: true
})
</script>

<template>
  <div>
    <div class="flex items-stretch gap-3 bg-gray-200 dark:bg-neutral-800 text-xs sm:text-sm">
      <div
        ref="breadcrumbsContainer"
        class="flex flex-1 items-center overflow-x-auto font-name"
      >
        <template
          v-for="(part, i) in breadcrumbs"
          :key="i"
        >
          <button
            class="flex items-center gap-1 p-2 text-gray-500 dark:text-neutral-400 select-none underline-offset-2 whitespace-nowrap"
            :class="(entries.isFile && i === breadcrumbs.length - 1) ? 'cursor-default' : 'cursor-pointer hover:text-gray-600 hover:dark:text-neutral-300 hover:underline'"
            :tabindex="(entries.isFile && i === breadcrumbs.length - 1) ? -1 : 0"
            @click="onBreadcrumbNavigate(i)"
          >
            <bi-file-earmark-text
              v-if="entries.isFile && i === breadcrumbs.length - 1"
              class="flex-shrink-0 w-3 h-3"
            />
            <bi-folder2-open
              v-else-if="i > 0"
              class="flex-shrink-0 w-3 h-3"
            />
            <span v-text="part" />
          </button>
          <bi-caret-right-fill class="flex-shrink-0 w-3 h-3 text-gray-400 dark:text-neutral-500 last:hidden" />
        </template>
      </div>
      <div class="flex items-center">
        <button
          v-if="entries.isFile"
          class="flex items-center gap-1 p-2 text-gray-500 dark:text-neutral-400 select-none cursor-pointer hover:text-gray-600 hover:dark:text-neutral-300 hover:underline underline-offset-2"
          @click="emit('download', entries.files[0])"
        >
          <bi-cloud-download class="flex-shrink-0 w-4 h-4" />
          <span>Download</span>
        </button>
      </div>
    </div>
    <share-file-preview
      v-if="entries.isFile"
      :share="share"
      :file="entries.files[0]"
      @download="emit('download', entries.files[0])"
    />
    <div
      v-else
      class="divide-y divide-gray-200 dark:divide-neutral-800"
    >
      <button
        v-if="breadcrumbs.length > 1"
        class="w-full flex items-center gap-3 p-2 sm:px-4 hover:bg-gray-300 dark:hover:bg-neutral-700 select-none cursor-pointer"
        @click="onBreadcrumbNavigate(breadcrumbs.length - 2)"
      >
        <span class="flex-1 flex items-center gap-1.5 sm:gap-2 font-name text-sm sm:text-base font-medium overflow-x-hidden whitespace-nowrap">
          <bi-folder2 class="flex-shrink-0 w-4 h-4" />
          <span>..</span>
        </span>
      </button>
      <button
        v-for="directory of entries.directories"
        :key="'directory-' + currentPath + directory.name"
        class="w-full flex items-center gap-3 p-2 sm:px-4 hover:bg-gray-300 dark:hover:bg-neutral-700 select-none cursor-pointer"
        @click="currentPath += directory.name + '/'"
      >
        <span class="flex-1 flex items-center gap-1.5 sm:gap-2 font-name text-sm sm:text-base font-medium overflow-x-hidden whitespace-nowrap">
          <bi-folder2 class="flex-shrink-0 w-4 h-4" />
          <span v-text="directory.name" />
        </span>
        <span
          class="flex-shrink-0 text-xs sm:text-sm text-gray-500 dark:text-neutral-400"
          v-text="formatSize(directory.size)"
        />
      </button>
      <button
        v-for="file of entries.files"
        :key="file.id"
        class="w-full flex items-center gap-3 p-2 sm:px-4 hover:bg-gray-300 dark:hover:bg-neutral-700 select-none cursor-pointer"
        @click="currentPath = file.path"
      >
        <span class="flex-1 flex items-center gap-1.5 sm:gap-2 font-name text-sm sm:text-base font-medium overflow-x-hidden whitespace-nowrap">
          <bi-file-earmark-text class="flex-shrink-0 w-4 h-4" />
          <span v-text="file.name" />
        </span>
        <span
          v-if="file.size"
          class="flex-shrink-0 text-xs sm:text-sm text-gray-500 dark:text-neutral-400"
          v-text="formatSize(file.size)"
        />
      </button>
    </div>
  </div>
</template>
