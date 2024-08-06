import { useLocalStorage } from '@vueuse/core'

const passwords = useLocalStorage<Record<string, string>>('share_passwords', {}, {
  writeDefaults: false,
  initOnMounted: false
})

export const useSharePasswords = () => {
  return passwords
}
