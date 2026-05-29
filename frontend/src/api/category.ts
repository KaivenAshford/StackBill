import request from '@/utils/request'
import type { Category } from '@/types'

export function listCategories(params?: { type?: string }) {
  return request.get<Category[]>('/categories', { params })
}

export function getCategory(id: number) {
  return request.get<Category>(`/categories/${id}`)
}

export function createCategory(data: Partial<Category>) {
  return request.post<Category>('/categories', data)
}

export function updateCategory(id: number, data: Partial<Category>) {
  return request.put<Category>(`/categories/${id}`, data)
}

export function deleteCategory(id: number) {
  return request.delete(`/categories/${id}`)
}
