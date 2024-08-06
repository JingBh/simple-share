import { useSharePasswords } from './share-passwords.ts'
import type { Share } from '../types/share.ts'

const passwords = useSharePasswords()

export const getContentUrl = (share: Share, fileId?: string): string => {
  let url = `/api/shares/${share.name}`
  if (fileId) {
    url += `/files/${fileId}`
  } else {
    url += '/content'
  }

  if (passwords.value[share.name]) {
    const params = new URLSearchParams()
    params.append('password', passwords.value[share.name])
    url += `?${params.toString()}`
  }

  return url
}
