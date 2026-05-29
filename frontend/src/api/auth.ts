import request from '@/utils/request'
import type { User } from '@/types'

export function login(username: string, password: string) {
  return request.post<{ token: string; user: User }>('/auth/login', { username, password })
}

export function register(username: string, email: string, password: string) {
  return request.post<{ token: string; user: User }>('/auth/register', { username, email, password })
}

export function getMe() {
  return request.get<User>('/auth/me')
}

export function updateProfile(data: { nickname?: string; avatar?: string }) {
  return request.put<User>('/users/profile', data)
}

export function updatePassword(data: { old_password: string; new_password: string }) {
  return request.put<null>('/users/password', data)
}