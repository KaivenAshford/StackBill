import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Subscription, PageResult } from '@/types'
import { listSubscriptions } from '@/api/subscription'

export const useSubscriptionStore = defineStore('subscription', () => {
  const subscriptions = ref<Subscription[]>([])
  const total = ref(0)

  async function fetchSubscriptions(page = 1, pageSize = 20) {
    const res = await listSubscriptions({ page, page_size: pageSize })
    subscriptions.value = (res.data as unknown as PageResult<Subscription>).items
    total.value = (res.data as unknown as PageResult<Subscription>).total
  }

  return { subscriptions, total, fetchSubscriptions }
})
