<script setup lang="ts">
import BiFileEarmarkText from 'bootstrap-icons/icons/file-earmark-text.svg?component'
import BiFolder2 from 'bootstrap-icons/icons/folder2.svg?component'
import BiFonts from 'bootstrap-icons/icons/fonts.svg?component'
import BiLink45deg from 'bootstrap-icons/icons/link-45deg.svg?component'

defineProps<{
  shareType: string
}>()
</script>

<template>
  <bi-folder2 v-if="shareType === 'directory'" />
  <bi-file-earmark-text v-else-if="shareType === 'file'" />
  <bi-fonts v-else-if="shareType === 'text'" />
  <bi-link45deg v-else-if="shareType === 'url'" />
</template>
