import axios from 'axios'
import type { AxiosInstance } from 'axios'

let instance: AxiosInstance | null = null

export const useAxiosInstance = (): AxiosInstance => {
  if (instance === null) {
    instance = axios.create({
      withCredentials: true
    })
  }

  return instance
}
