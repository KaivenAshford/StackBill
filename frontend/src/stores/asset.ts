import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Asset, PageResult } from '@/types'
import { listAssets } from '@/api/asset'

export const useAssetStore = defineStore('asset', () => {
  const assets = ref<Asset[]>([])
  const total = ref(0)

  async function fetchAssets(page = 1, pageSize = 20) {
    const res = await listAssets({ page, page_size: pageSize })
    assets.value = (res.data as unknown as PageResult<Asset>).items
    total.value = (res.data as unknown as PageResult<Asset>).total
  }

  return { assets, total, fetchAssets }
})
