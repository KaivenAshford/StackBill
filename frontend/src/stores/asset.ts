import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Asset, PageResult } from '@/types'
import { listAssets, type AssetQuery } from '@/api/asset'

export const useAssetStore = defineStore('asset', () => {
  const assets = ref<Asset[]>([])
  const total = ref(0)
  const loaded = ref(false)
  const currentQuery = ref<AssetQuery>({})

  async function ensureLoaded(params?: AssetQuery) {
    if (loaded.value) return
    const query = params || currentQuery.value
    currentQuery.value = query
    const res = await listAssets(query)
    assets.value = (res.data as unknown as PageResult<Asset>).items
    total.value = (res.data as unknown as PageResult<Asset>).total
    loaded.value = true
  }

  function invalidate() {
    assets.value = []
    total.value = 0
    loaded.value = false
  }

  async function refresh(params?: AssetQuery) {
    invalidate()
    await ensureLoaded(params)
  }

  return { assets, total, loaded, currentQuery, ensureLoaded, invalidate, refresh }
})
