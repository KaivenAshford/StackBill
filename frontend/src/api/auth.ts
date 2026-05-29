import request from '@/utils/request'
import type { User } from '@/types'

export function login(username: string, password: string) {
  return request.post<{ token: string; user: User }>('/auth/login', { username, password })
}

export function register(username: string, email: string, password: string) {
  return request.post<{ token: string; user: User }>('/auth/register', { username, email, password })
}

export function getCurrentUser() {
  return request.get<User>('/auth/me')
}
