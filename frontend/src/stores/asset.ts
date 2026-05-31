import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Asset, PageResult } from '@/types'
import { listAssets } from '@/api/asset'

export const useAssetStore = defineStore('asset', () => {
  const assets = ref<Asset[]>([])
  const total = ref(0)
  const loaded = ref(false)

  async function ensureLoaded(page = 1, pageSize = 20) {
    if (loaded.value) return
    const res = await listAssets({ page, page_size: pageSize })
    assets.value = (res.data as unknown as PageResult<Asset>).items
    total.value = (res.data as unknown as PageResult<Asset>).total
    loaded.value = true
  }

  function invalidate() {
    assets.value = []
    total.value = 0
    loaded.value = false
  }

  async function refresh(page = 1, pageSize = 20) {
    invalidate()
    await ensureLoaded(page, pageSize)
  }

  return { assets, total, loaded, ensureLoaded, invalidate, refresh }
})
