import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Subscription, PageResult } from '@/types'
import { listSubscriptions } from '@/api/subscription'

export const useSubscriptionStore = defineStore('subscription', () => {
  const subscriptions = ref<Subscription[]>([])
  const total = ref(0)
  const loaded = ref(false)

  async function ensureLoaded(page = 1, pageSize = 20) {
    if (loaded.value) return
    const res = await listSubscriptions({ page, page_size: pageSize })
    subscriptions.value = (res.data as unknown as PageResult<Subscription>).items
    total.value = (res.data as unknown as PageResult<Subscription>).total
    loaded.value = true
  }

  function invalidate() {
    subscriptions.value = []
    total.value = 0
    loaded.value = false
  }

  async function refresh(page = 1, pageSize = 20) {
    invalidate()
    await ensureLoaded(page, pageSize)
  }

  return { subscriptions, total, loaded, ensureLoaded, invalidate, refresh }
})
