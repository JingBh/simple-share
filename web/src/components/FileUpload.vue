<script setup lang="ts">
import { computed, onActivated, onDeactivated, onUnmounted, ref } from 'vue'
import { useDropZone } from '@vueuse/core'
import type { AxiosError } from 'axios'

import { useAxiosInstance } from '../lib/axios.ts'
import { formatSize } from '../utils/filesize.ts'
import type { ShareFileUpload } from '../types/share.ts'

import BiCloudUpload from 'bootstrap-icons/icons/cloud-upload.svg?component'

type QueuedFile = { file: File, path: string }

const active = ref(false)

const props = defineProps<{
  modelValue: ShareFileUpload[]
  drop?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: ShareFileUpload[]]
}>()

let abortController = null as AbortController | null

const uploading = ref<QueuedFile | null>(null)

const uploadingDone = ref<number>(0)

const uploadingTotal = computed<number>(() => {
  if (uploading.value) {
    return uploading.value.file.size
  }
  return 0
})

const uploadingProgress = computed<number>(() => {
  if (uploading.value) {
    return uploadingDone.value / uploadingTotal.value
  }
  return 0
})

const queue = ref<QueuedFile[]>([])

const failed = ref<QueuedFile[]>([])

const allPaths = computed(() => {
  const paths = new Set<string>()
  for (const file of props.modelValue) {
    paths.add(file.path)
  }
  for (const file of queue.value) {
    paths.add(file.path)
  }
  return paths
})

const checkLimit = (files: any[]): boolean => {
  if (props.modelValue.length + queue.value.length + files.length >= 100) {
    alert('Currently, we only support uploading up to 100 files at a time.')
    return false
  }
  return true
}

const scanEntry = async (entry: FileSystemEntry, basePath: string): Promise<QueuedFile[]> => {
  if (!entry.name) {
    return []
  }

  const files = [] as QueuedFile[]
  const promises = [] as Promise<void>[]
  if (entry.isFile) {
    promises.push(new Promise<void>((resolve) => {
      (entry as FileSystemFileEntry).file((file) => {
        const path = `${basePath}${entry.name}`
        if (!allPaths.value.has(path)) {
          files.push({ file, path })
        }
        resolve()
      }, () => resolve())
    }))
  } else if (entry.isDirectory) {
    promises.push(new Promise<void>((resolve) => {
      const reader = (entry as FileSystemDirectoryEntry).createReader()
      reader.readEntries((subEntries) => {
        const subPromises = [] as Promise<void>[]
        for (const subEntry of subEntries) {
          subPromises.push(new Promise<void>((subResolve) => {
            scanEntry(subEntry, `${basePath}${entry.name}/`)
              .then((subFiles) => files.push(...subFiles))
              .finally(() => subResolve())
          }))
        }
        Promise.all(subPromises).then(() => resolve())
      }, () => resolve())
    }))
  }

  await Promise.all(promises)
  return files
}

const select = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.multiple = true
  // input.webkitdirectory = true

  input.onchange = async () => {
    const files = [] as QueuedFile[]
    if (input.webkitEntries && input.webkitEntries.length > 0) {
      for (const entry of input.webkitEntries) {
        files.push(...(await scanEntry(entry, '')))
      }
    } else {
      for (const file of input.files ?? []) {
        if (file.name) {
          const path = file.webkitRelativePath || file.name
          if (!allPaths.value.has(path)) {
            files.push({ file, path })
          }
        }
      }
    }

    if (files.length > 0 && checkLimit(files)) {
      queue.value.push(...files)
      startQueue()
    }
  }

  input.click()
}

const { isOverDropZone } = useDropZone(document.body, {
  onDrop: async (_, e) => {
    if (!active.value) {
      return
    }

    const promises = [] as Promise<void>[]
    const files = [] as QueuedFile[]
    const items = e.dataTransfer?.items ?? []
    for (const item of items) {
      if (!item.webkitGetAsEntry) {
        const file = item.getAsFile()
        if (file && file.name) {
          const path = file.webkitRelativePath || file.name
          if (!allPaths.value.has(path)) {
            files.push({ file, path })
          }
        }
      }
      const entry = item.webkitGetAsEntry()
      if (entry) {
        promises.push(new Promise<void>((resolve) => {
          scanEntry(entry, '')
            .then((subFiles) => files.push(...subFiles))
            .finally(() => resolve())
        }))
      }
    }
    await Promise.all(promises)

    if (files.length > 0 && checkLimit(files)) {
      queue.value.push(...files)
      startQueue()
    }
  }
})

const startQueue = async () => {
  if (uploading.value || !active.value) {
    return
  }
  uploading.value = queue.value.shift() || null
  uploadingDone.value = 0
  if (uploading.value) {
    const file = uploading.value

    try {
      // step 1: initiate multipart upload
      const { id: fileId, partSize } = (await useAxiosInstance().post<{
        id: string
        partSize: number
      }>('/api/upload')).data

      // step 2: upload parts
      const partTotal = Math.ceil(file.file.size / partSize)
      for (let i = 1; i <= partTotal; i++) {
        const start = (i - 1) * partSize
        const end = Math.min(i * partSize, file.file.size)
        const part = file.file.slice(start, end)

        for (let tries = 1; tries <= 3; tries++) {
          abortController = new AbortController()
          try {
            await useAxiosInstance().post(`/api/upload/${fileId}/${i}`, part, {
              headers: {
                'Content-Type': 'application/octet-stream',
              },
              signal: abortController.signal,
              onUploadProgress: (e) => {
                uploadingDone.value = start + e.loaded
              },
            })
            break
          } catch (e: AxiosError | any) {
            if (tries === 3 || e.name === 'AbortError') {
              throw e
            }
            console.error(e)
          }
        }
      }

      // step 3: complete multipart upload
      await useAxiosInstance().post(`/api/upload/${fileId}/complete`)
      uploadingDone.value = file.file.size
      emit('update:modelValue', [
        ...props.modelValue,
        { id: fileId, path: file.path }
      ])
    } catch (e) {
      console.error(e)
      failed.value.push(uploading.value)
    }

    uploading.value = null
    startQueue()
  }
}

const retry = () => {
  const failedFiles = failed.value.splice(0)
  queue.value.unshift(...failedFiles)
  startQueue()
}

const abort = () => {
  if (uploading.value && abortController) {
    abortController.abort()
  }
}

onActivated(() => {
  active.value = true
  startQueue()
})

onDeactivated(() => {
  active.value = false
})

onUnmounted(() => {
  abort()
})
</script>

<template>
  <div class="flex flex-col items-stretch gap-2">
    <div
      class="w-full py-8 flex flex-col items-center justify-center gap-4 bg-white dark:bg-neutral-800 rounded-lg border-2 border-dashed hover:border-gray-400 dark:hover:border-neutral-600 cursor-pointer select-none"
      :class="drop && isOverDropZone ? 'border-gray-400 dark:hover:border-neutral-600' : 'border-gray-300 dark:border-neutral-700'"
      @click="select"
    >
      <bi-cloud-upload class="w-12 h-12 text-gray-300 dark:text-neutral-600" />
      <p class="text-xs font-semibold text-gray-500 dark:text-neutral-400">
      <span class="hidden sm:inline-flex flex-col items-center">
        <template v-if="drop">
          <span>Drop your files anywhere</span>
          <span class="text-gray-400 dark:text-neutral-500 my-0.5">OR</span>
        </template>
        <span>Click to browse files</span>
      </span>
        <span class="sm:hidden">Touch to select files</span>
      </p>
    </div>
    <p
      v-if="uploading"
      class="text-blue-500 text-sm"
    >
      Uploading <strong>{{ uploading.path }}</strong>... ({{ formatSize(uploadingDone) }} / {{ formatSize(uploadingTotal) }}, {{ Math.round(uploadingProgress * 100) }}%)
    </p>
    <p
      v-if="props.modelValue.length || queue.length"
      class="text-green-500 text-sm"
    >
      <span v-if="props.modelValue.length"><strong>{{ props.modelValue.length }}</strong> file{{ props.modelValue.length > 1 ? 's' : '' }} uploaded</span>
      <span v-if="props.modelValue.length && queue.length">, </span>
      <span v-if="queue.length"><strong>{{ queue.length }}</strong> more file{{ queue.length > 1 ? 's' : '' }} to upload</span>
    </p>
    <p
      v-if="failed.length"
      class="text-red-500 text-sm"
    >
      <strong>{{ failed.length }}</strong> file{{ failed.length > 1 ? 's' : '' }} failed to upload, <a class="underline underline-offset-2 hover:text-red-700 dark:hover:text-red-300" href="#" @click.prevent="retry">click to retry</a>
    </p>
    <Teleport to="body">
      <div
        v-if="active && drop && isOverDropZone"
        class="fixed top-0 left-0 w-screen h-screen z-50 bg-black/10 dark:bg-black/30 pointer-events-none"
      />
    </Teleport>
  </div>
</template>
