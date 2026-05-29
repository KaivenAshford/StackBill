import request from '@/utils/request'
import type { Subscription, PageResult } from '@/types'

export interface SubscriptionQuery {
  page?: number
  page_size?: number
  category_id?: number
  status?: string
  upcoming_renewal?: boolean
}

export function listSubscriptions(params?: SubscriptionQuery) {
  return request.get<PageResult<Subscription>>('/subscriptions', { params })
}

export function getSubscription(id: number) {
  return request.get<Subscription>(`/subscriptions/${id}`)
}

export function createSubscription(data: Partial<Subscription>) {
  return request.post<Subscription>('/subscriptions', data)
}

export function updateSubscription(id: number, data: Partial<Subscription>) {
  return request.put<Subscription>(`/subscriptions/${id}`, data)
}

export function deleteSubscription(id: number) {
  return request.delete(`/subscriptions/${id}`)
}
