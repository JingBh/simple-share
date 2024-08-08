import { useSharePasswords } from './share-passwords.ts'
import type { Share } from '../types/share.ts'

const passwords = useSharePasswords()

const getUrl = (share: Share, fileId?: string, type?: string): string => {
  let url = `/api/shares/${share.name}`
  if (fileId) {
    url += `/files/${fileId}`
  } else {
    url += '/content'
  }
  if (type) {
    url += `/${type}`
  }

  if (passwords.value[share.name]) {
    const params = new URLSearchParams()
    params.append('password', passwords.value[share.name])
    url += `?${params.toString()}`
  }

  return url
}

export const getContentUrl = (share: Share, fileId?: string): string => {
  return getUrl(share, fileId)
}

export const getContentTypeUrl = (share: Share, fileId: string): string => {
  return getUrl(share, fileId, 'type')
}

export const getPreviewUrl = (share: Share, fileId: string): string => {
  return getUrl(share, fileId, 'preview')
}
