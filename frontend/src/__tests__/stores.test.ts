import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useSubscriptionStore } from '@/stores/subscription'
import { useAssetStore } from '@/stores/asset'
import { useCategoryStore } from '@/stores/category'

// Mock API modules
vi.mock('@/api/subscription', () => ({
  listSubscriptions: vi.fn().mockResolvedValue({
    data: {
      items: [
        { id: 1, name: 'Netflix', amount: 15.99, currency: 'USD', status: 'active', billing_cycle: 'monthly' },
        { id: 2, name: 'Spotify', amount: 9.99, currency: 'USD', status: 'active', billing_cycle: 'monthly' },
      ],
      total: 2,
      page: 1,
      page_size: 20,
    },
  }),
}))

vi.mock('@/api/asset', () => ({
  listAssets: vi.fn().mockResolvedValue({
    data: {
      items: [
        { id: 1, name: 'example.com', asset_type: 'domain', status: 'active' },
      ],
      total: 1,
      page: 1,
      page_size: 20,
    },
  }),
}))

vi.mock('@/api/category', () => ({
  listCategories: vi.fn().mockResolvedValue({
    data: [
      { id: 1, name: 'AI Tools', type: 'subscription', color: '#42b883' },
      { id: 2, name: 'Domains', type: 'asset', color: '#3273dc' },
    ],
  }),
}))

describe('useSubscriptionStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with empty subscriptions', () => {
    const store = useSubscriptionStore()
    expect(store.subscriptions).toEqual([])
    expect(store.total).toBe(0)
    expect(store.loaded).toBe(false)
  })

  it('loads subscriptions from API', async () => {
    const store = useSubscriptionStore()
    await store.ensureLoaded()
    expect(store.loaded).toBe(true)
    expect(store.subscriptions.length).toBe(2)
    expect(store.total).toBe(2)
    expect(store.subscriptions[0].name).toBe('Netflix')
  })

  it('invalidates cache', async () => {
    const store = useSubscriptionStore()
    await store.ensureLoaded()
    expect(store.loaded).toBe(true)
    store.invalidate()
    expect(store.loaded).toBe(false)
    expect(store.subscriptions).toEqual([])
  })

  it('refresh reloads data', async () => {
    const store = useSubscriptionStore()
    await store.ensureLoaded()
    await store.refresh()
    expect(store.loaded).toBe(true)
    expect(store.subscriptions.length).toBe(2)
  })
})

describe('useAssetStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with empty assets', () => {
    const store = useAssetStore()
    expect(store.assets).toEqual([])
    expect(store.total).toBe(0)
    expect(store.loaded).toBe(false)
  })

  it('loads assets from API', async () => {
    const store = useAssetStore()
    await store.ensureLoaded()
    expect(store.loaded).toBe(true)
    expect(store.assets.length).toBe(1)
    expect(store.total).toBe(1)
    expect(store.assets[0].name).toBe('example.com')
  })
})

describe('useCategoryStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with empty categories', () => {
    const store = useCategoryStore()
    expect(store.categories).toEqual([])
    expect(store.loaded).toBe(false)
  })

  it('loads categories from API', async () => {
    const store = useCategoryStore()
    await store.ensureLoaded()
    expect(store.loaded).toBe(true)
    expect(store.categories.length).toBe(2)
    expect(store.categories[0].name).toBe('AI Tools')
  })
})
