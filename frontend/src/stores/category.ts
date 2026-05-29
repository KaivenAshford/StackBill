import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Category } from '@/types'
import { listCategories } from '@/api/category'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([])
  const loaded = ref(false)

  async function fetchCategories(type?: string) {
    const res = await listCategories(type ? { type } : undefined)
    categories.value = res.data
    loaded.value = true
  }

  return { categories, loaded, fetchCategories }
})
