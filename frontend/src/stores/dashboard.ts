import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getDashboard, type DashboardData } from '@/api/dashboard'

export const useDashboardStore = defineStore('dashboard', () => {
  const data = ref<DashboardData | null>(null)
  const loaded = ref(false)

  async function ensureLoaded() {
    if (loaded.value) return
    const res = await getDashboard()
    data.value = res.data
    loaded.value = true
  }

  function invalidate() {
    data.value = null
    loaded.value = false
  }

  async function refresh() {
    invalidate()
    await ensureLoaded()
  }

  return { data, loaded, ensureLoaded, invalidate, refresh }
})
