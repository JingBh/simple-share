<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import type { AxiosError } from 'axios'

import { useAxiosInstance } from '../lib/axios.ts'
import LayoutDashboard from '../layouts/LayoutDashboard.vue'
import FileUpload from '../components/FileUpload.vue'
import FlowbiteSpinner from '../components/FlowbiteSpinner.vue'
import type { GenerateShareSettings, ShareFileUpload } from '../types/share.ts'

import BiFolder2 from 'bootstrap-icons/icons/folder2.svg?component'
import BiFonts from 'bootstrap-icons/icons/fonts.svg?component'
import BiLink45deg from 'bootstrap-icons/icons/link-45deg.svg?component'

const router = useRouter()

const files = ref<ShareFileUpload[]>([])
const contents = ref('')
const settings = ref<GenerateShareSettings>({
  type: 'file',
  name: '',
  displayName: '',
  nameRandom: true,
  nameRandomLength: 6,
  password: '',
  expiry: 0
})

const isSubmitting = ref(false)
const errors = ref<Record<string, Record<string, string>>>({})
const errorMessage = ref('')
const submit = () => {
  if (isSubmitting.value) {
    return
  }
  errors.value = {}
  errorMessage.value = ''
  isSubmitting.value = true
  useAxiosInstance().post<{
    name: string
  }>('/api/shares', Object.assign({
    text: contents.value,
    files: files.value
  }, settings.value)).then(({ data }) => {
    router.replace(`/shares/${data.name}`)
  }).catch((e: AxiosError | any) => {
    if (e.response?.status === 422) {
      errorMessage.value = 'There are one or more errors in the form:'
      errors.value = e.response.data
    } else if (e.response?.data?.message) {
      errorMessage.value = 'Request failed: ' + e.response.data.message
    } else {
      console.error(e)
      errorMessage.value = 'Request failed.'
    }
  }).finally(() => {
    isSubmitting.value = false
  })
}
</script>

<template>
  <layout-dashboard>
    <div class="flex flex-col items-stretch gap-4 md:gap-6">
      <div class="flex flex-col sm:flex-row sm:items-center gap-2">
        <label class="font-medium text-neutral-900 dark:text-white">
          I want to share
        </label>
        <div class="inline-flex rounded-md shadow-sm">
          <button
            type="button"
            class="flex-1 inline-flex whitespace-nowrap justify-center items-center gap-2 px-4 py-2 text-sm font-medium bg-white border border-neutral-200 rounded-s-lg hover:bg-neutral-100 focus:z-10 focus:ring-2 focus:ring-blue-700 focus:text-blue-700 dark:bg-neutral-800 dark:border-neutral-700 dark:hover:bg-neutral-700 dark:focus:ring-blue-500 dark:focus:text-blue-500"
            :class="{ 'text-blue-700 dark:text-blue-500': settings.type === 'file', 'text-neutral-900 dark:text-neutral-100 hover:text-blue-700 dark:hover:text-blue-500': settings.type !== 'file' }"
            @click="settings.type = 'file'"
          >
            <bi-folder2 class="w-4 h-4" />
            <span>Files</span>
          </button>
          <button
            type="button"
            class="flex-1 inline-flex whitespace-nowrap justify-center items-center gap-2 px-4 py-2 text-sm font-medium bg-white border-t border-b border-neutral-200 hover:bg-neutral-100 focus:z-10 focus:ring-2 focus:ring-blue-700 focus:text-blue-700 dark:bg-neutral-800 dark:border-neutral-700 dark:hover:bg-neutral-700 dark:focus:ring-blue-500 dark:focus:text-blue-500"
            :class="{ 'text-blue-700 dark:text-blue-500': settings.type === 'text', 'text-neutral-900 dark:text-neutral-100 hover:text-blue-700 dark:hover:text-blue-500': settings.type !== 'text' }"
            @click="settings.type = 'text'"
          >
            <bi-fonts class="w-4 h-4" />
            <span>Text</span>
          </button>
          <button
            type="button"
            class="flex-1 inline-flex whitespace-nowrap justify-center items-center gap-2 px-4 py-2 text-sm font-medium bg-white border border-neutral-200 rounded-e-lg hover:bg-neutral-100 focus:z-10 focus:ring-2 focus:ring-blue-700 focus:text-blue-700 dark:bg-neutral-800 dark:border-neutral-700 dark:hover:bg-neutral-700 dark:focus:ring-blue-500 dark:focus:text-blue-500"
            :class="{ 'text-blue-700 dark:text-blue-500': settings.type === 'url', 'text-neutral-900 dark:text-neutral-100 hover:text-blue-700 dark:hover:text-blue-500': settings.type !== 'url' }"
            @click="settings.type = 'url'"
          >
            <bi-link45deg class="w-4 h-4" />
            <span>URL</span>
          </button>
        </div>
      </div>
      <keep-alive>
        <file-upload
          v-if="settings.type === 'file'"
          v-model="files"
          drop
        />
        <textarea
          v-else-if="settings.type === 'text'"
          v-model="contents"
          class="block p-2.5 w-full min-h-32 text-sm font-mono text-neutral-900 bg-neutral-50 rounded-lg border border-neutral-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-neutral-700 dark:border-neutral-600 dark:placeholder-neutral-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          placeholder="Input here the text you want to share"
          rows="12"
          autocomplete="off"
        />
        <div
          v-else-if="settings.type === 'url'"
          class="flex flex-col gap-2"
        >
          <textarea
            v-model="contents"
            class="block p-2.5 w-full resize-none text-sm font-mono text-neutral-900 bg-neutral-50 rounded-lg border border-neutral-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-neutral-700 dark:border-neutral-600 dark:placeholder-neutral-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="Input here the URL you want to share"
            rows="4"
            autocomplete="off"
          />
          <p class="text-xs text-neutral-500 dark:text-neutral-400">
            URL shares works like an URL shortener, will immediately redirect when accessed.
          </p>
        </div>
      </keep-alive>
      <hr class="my-4 border-neutral-300 dark:border-neutral-700" />
      <label class="flex items-center cursor-pointer">
        <input
          v-model="settings.nameRandom"
          type="checkbox"
          class="sr-only peer"
          autocomplete="off"
        />
        <span class="relative w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-neutral-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:w-5 after:h-5 after:transition-all dark:border-neutral-600 peer-checked:bg-blue-600" />
        <span class="ms-3 font-medium text-neutral-700 dark:text-neutral-300">
          Random link
        </span>
      </label>
      <label
        v-if="settings.nameRandom"
        class="flex items-center gap-2"
      >
        <span class="font-medium text-neutral-700 dark:text-neutral-300">
          Link length
        </span>
        <input
          v-model="settings.nameRandomLength"
          type="number"
          class="bg-neutral-50 border border-neutral-300 text-neutral-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 py-1.5 px-2.5 dark:bg-neutral-700 dark:border-neutral-600 dark:placeholder-neutral-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          min="4"
          max="32"
          step="1"
          autocomplete="off"
        />
      </label>
      <label
        v-else
        class="flex flex-col sm:flex-row sm:items-center gap-2"
      >
        <span class="font-medium text-neutral-700 dark:text-neutral-300">
          Link
        </span>
        <input
          v-model="settings.name"
          type="text"
          class="bg-neutral-50 border border-neutral-300 text-neutral-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 py-1.5 px-2.5 dark:bg-neutral-700 dark:border-neutral-600 dark:placeholder-neutral-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          placeholder="Custom link name"
          autocomplete="off"
        />
      </label>
      <label class="flex flex-col sm:flex-row sm:items-center gap-2">
        <span class="font-medium text-neutral-700 dark:text-neutral-300">
          Title
        </span>
        <input
          v-model="settings.displayName"
          type="text"
          class="bg-neutral-50 border border-neutral-300 text-neutral-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 py-1.5 px-2.5 dark:bg-neutral-700 dark:border-neutral-600 dark:placeholder-neutral-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          placeholder="Display name"
          autocomplete="off"
        />
      </label>
      <label class="flex items-center gap-2">
        <span class="font-medium text-neutral-700 dark:text-neutral-300">
          Password
        </span>
        <input
          v-model="settings.password"
          type="password"
          class="bg-neutral-50 border border-neutral-300 text-neutral-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 py-1.5 px-2.5 dark:bg-neutral-700 dark:border-neutral-600 dark:placeholder-neutral-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          placeholder="Publicly available"
          autocomplete="off"
        />
      </label>
      <div class="flex flex-col sm:flex-row sm:items-center gap-2">
        <label class="font-medium text-neutral-700 dark:text-neutral-300">
          Expires in
        </label>
        <div class="inline-flex rounded-md shadow-sm">
          <button
            type="button"
            class="inline-flex flex-1 whitespace-nowrap justify-center items-center gap-2 px-3 py-1.5 text-sm font-medium bg-white border border-neutral-200 rounded-s-lg hover:bg-neutral-100 focus:z-10 focus:ring-2 focus:ring-blue-700 focus:text-blue-700 dark:bg-neutral-800 dark:border-neutral-700 dark:hover:bg-neutral-700 dark:focus:ring-blue-500 dark:focus:text-blue-500"
            :class="{ 'text-blue-700 dark:text-blue-500': settings.expiry === 0, 'text-neutral-900 dark:text-neutral-100 hover:text-blue-700 dark:hover:text-blue-500': settings.expiry !== 0 }"
            @click="settings.expiry = 0"
          >
            Never
          </button>
          <button
            type="button"
            class="inline-flex flex-1 whitespace-nowrap justify-center items-center gap-2 px-3 py-1.5 text-sm font-medium bg-white border border-l-0 border-neutral-200 hover:bg-neutral-100 focus:z-10 focus:ring-2 focus:ring-blue-700 focus:text-blue-700 dark:bg-neutral-800 dark:border-neutral-700 dark:hover:bg-neutral-700 dark:focus:ring-blue-500 dark:focus:text-blue-500"
            :class="{ 'text-blue-700 dark:text-blue-500': settings.expiry === 1, 'text-neutral-900 dark:text-neutral-100 hover:text-blue-700 dark:hover:text-blue-500': settings.expiry !== 1 }"
            @click="settings.expiry = 1"
          >
            1 day
          </button>
          <button
            type="button"
            class="inline-flex flex-1 whitespace-nowrap justify-center items-center gap-2 px-3 py-1.5 text-sm font-medium bg-white border border-l-0 border-neutral-200 hover:bg-neutral-100 focus:z-10 focus:ring-2 focus:ring-blue-700 focus:text-blue-700 dark:bg-neutral-800 dark:border-neutral-700 dark:hover:bg-neutral-700 dark:focus:ring-blue-500 dark:focus:text-blue-500"
            :class="{ 'text-blue-700 dark:text-blue-500': settings.expiry === 3, 'text-neutral-900 dark:text-neutral-100 hover:text-blue-700 dark:hover:text-blue-500': settings.expiry !== 3 }"
            @click="settings.expiry = 3"
          >
            3 days
          </button>
          <button
            type="button"
            class="inline-flex flex-1 whitespace-nowrap justify-center items-center gap-2 px-3 py-1.5 text-sm font-medium bg-white border border-l-0 border-neutral-200 rounded-e-lg hover:bg-neutral-100 focus:z-10 focus:ring-2 focus:ring-blue-700 focus:text-blue-700 dark:bg-neutral-800 dark:border-neutral-700 dark:hover:bg-neutral-700 dark:focus:ring-blue-500 dark:focus:text-blue-500"
            :class="{ 'text-blue-700 dark:text-blue-500': settings.expiry === 7, 'text-neutral-900 dark:text-neutral-100 hover:text-blue-700 dark:hover:text-blue-500': settings.expiry !== 7 }"
            @click="settings.expiry = 7"
          >
            7 days
          </button>
        </div>
      </div>
      <button
        type="button"
        class="sm:max-w-64 px-5 py-2.5 text-base font-medium text-white bg-blue-700 focus:ring-4 focus:outline-none focus:ring-blue-300 rounded-lg text-center dark:bg-blue-600 dark:focus:ring-blue-800"
        :class="{ 'cursor-pointer hover:bg-blue-800 dark:hover:bg-blue-700': !isSubmitting, 'cursor-default': isSubmitting }"
        :disabled="isSubmitting"
        @click="submit"
      >
        <span
          v-if="isSubmitting"
          class="inline-flex items-center gap-2"
        >
          <flowbite-spinner class="w-5 h-5" inverse />
          <span>Loading...</span>
        </span>
        <span v-else>Submit</span>
      </button>
      <div class="space-y-2 text-sm text-red-600 dark:text-red-400">
        <p
          v-if="errorMessage"
          class="font-medium"
          v-text="errorMessage"
        />
        <ul
          v-if="Object.keys(errors).length"
          class="list-disc list-inside space-y-1"
        >
          <template v-for="(fieldErrors, field) in errors">
            <li
              v-for="(message, validation) in fieldErrors"
              :key="field + '-' + validation"
              v-text="message"
            />
          </template>
        </ul>
      </div>
    </div>
  </layout-dashboard>
</template>
