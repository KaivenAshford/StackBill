import { ref } from 'vue'

/**
 * Wraps an async fetch so that `loading` only flips to true when the
 * request takes longer than `delay` ms.  For fast APIs (< 150 ms) the
 * user never sees a loading / skeleton state — the data is already
 * there by the time the browser paints.
 */
export function useDeferredLoading(delay = 150) {
  const loading = ref(false)

  async function withLoading<T>(fn: () => Promise<T>): Promise<T> {
    const timer = setTimeout(() => { loading.value = true }, delay)
    try {
      return await fn()
    } finally {
      clearTimeout(timer)
      loading.value = false
    }
  }

  return { loading, withLoading }
}
