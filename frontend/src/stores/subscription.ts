import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Subscription, PageResult } from '@/types'
import { listSubscriptions, type SubscriptionQuery } from '@/api/subscription'

export const useSubscriptionStore = defineStore('subscription', () => {
  const subscriptions = ref<Subscription[]>([])
  const total = ref(0)
  const loaded = ref(false)
  const currentQuery = ref<SubscriptionQuery>({})

  async function ensureLoaded(params?: SubscriptionQuery) {
    if (loaded.value) return
    const query = params || currentQuery.value
    currentQuery.value = query
    const res = await listSubscriptions(query)
    subscriptions.value = (res.data as unknown as PageResult<Subscription>).items
    total.value = (res.data as unknown as PageResult<Subscription>).total
    loaded.value = true
  }

  function invalidate() {
    subscriptions.value = []
    total.value = 0
    loaded.value = false
  }

  async function refresh(params?: SubscriptionQuery) {
    invalidate()
    await ensureLoaded(params)
  }

  return { subscriptions, total, loaded, currentQuery, ensureLoaded, invalidate, refresh }
})
