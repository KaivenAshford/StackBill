import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'
import { getMe } from '@/api/auth'

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
    // Invalidate all cache stores on logout
    invalidateAllCaches()
  }

  function isLoggedIn() {
    return !!token.value
  }

  async function fetchUser() {
    if (!token.value) return
    try {
      const res = await getMe()
      user.value = res.data
    } catch {
      logout()
    }
  }

  return { user, token, setToken, setUser, logout, isLoggedIn, fetchUser }
})

function invalidateAllCaches() {
  // Dynamic imports to avoid circular dependency
  import('@/stores/dashboard').then(m => m.useDashboardStore().invalidate())
  import('@/stores/category').then(m => m.useCategoryStore().invalidate())
  import('@/stores/subscription').then(m => m.useSubscriptionStore().invalidate())
  import('@/stores/asset').then(m => m.useAssetStore().invalidate())
}