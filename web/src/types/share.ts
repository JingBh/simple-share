export type ShareType = 'file' | 'directory' | 'text' | 'url'

export interface Share {
  type: ShareType
  name: string
  password?: string
  expiry?: number
  size: number
  createdAt?: string
  files?: ShareFile[]
  creator?: ShareCreator
}

export interface ShareSettings {
  type: ShareType
  name: string
  password: string
  expiry: number // in days, 0 means never
}

export interface GenerateShareSettings extends ShareSettings {
  nameRandom: boolean
  nameRandomLength: number
}

export interface ShareCreator {
  subject: string
  username?: string
}

export interface ShareFile {
  id: string
  path: string
  size?: number
}

export interface ShareFileUpload {
  id: string
  path: string
}
