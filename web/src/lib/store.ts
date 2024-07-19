import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { AxiosError } from 'axios'

import { useAxiosInstance } from './axios.ts'

export const useStore = defineStore('app', () => {
  const authDisabled = ref(false)
  const userinfo = ref<{
    subject: string
    username: string
  } | null>(null)
  const userinfoLoaded = ref(false)
  const userinfoLoading = ref(false)

  const loggedIn = computed(() => {
    return !!userinfo.value?.subject
  })

  const fetchUserinfo = async () => {
    userinfoLoading.value = true
    try {
      userinfo.value = (await useAxiosInstance().get<{
        subject: string
        username: string
      }>('auth/userinfo')).data
    } catch (e: AxiosError | any) {
      if (e.response?.status === 404) {
        authDisabled.value = true
      } else {
        console.error(e)
        authDisabled.value = false
      }
    }
    userinfoLoading.value = false
    userinfoLoaded.value = true
  }

  return {
    authDisabled,
    fetchUserinfo,
    loggedIn,
    userinfo,
    userinfoLoaded,
    userinfoLoading
  }
})
