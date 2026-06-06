import request from '@/utils/request'
import type { Asset, PageResult } from '@/types'

export interface AssetQuery {
  page?: number
  page_size?: number
  asset_type?: string
  status?: string
  expiring_days?: number
  keyword?: string
}

export function listAssets(params?: AssetQuery) {
  return request.get<PageResult<Asset>>('/assets', { params })
}

export function getAsset(id: number) {
  return request.get<Asset>(`/assets/${id}`)
}

export function createAsset(data: Partial<Asset>) {
  return request.post<Asset>('/assets', data)
}

export function updateAsset(id: number, data: Partial<Asset>) {
  return request.put<Asset>(`/assets/${id}`, data)
}

export function deleteAsset(id: number) {
  return request.delete(`/assets/${id}`)
}
