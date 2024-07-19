import { partial } from 'filesize'

export const formatSize = partial({
  base: 2,
  round: 1,
  standard: 'jedec'
})
