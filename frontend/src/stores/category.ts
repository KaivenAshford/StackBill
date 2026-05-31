import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Category } from '@/types'
import { listCategories } from '@/api/category'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([])
  const loaded = ref(false)

  async function ensureLoaded(type?: string) {
    if (loaded.value) return
    const res = await listCategories(type ? { type } : undefined)
    categories.value = res.data
    loaded.value = true
  }

  function invalidate() {
    categories.value = []
    loaded.value = false
  }

  async function refresh(type?: string) {
    invalidate()
    await ensureLoaded(type)
  }

  return { categories, loaded, ensureLoaded, invalidate, refresh }
})
