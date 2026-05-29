import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')

  function setToken(t: string) {
    token.value = t
    localStorage.setItem('token', t)
  }

  function setUser(u: User) {
    user.value = u
  }

  function logout() {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
  }

  function isLoggedIn() {
    return !!token.value
  }

  return { user, token, setToken, setUser, logout, isLoggedIn }
})
