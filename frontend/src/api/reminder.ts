import request from '@/utils/request'
import type { Reminder, PageResult } from '@/types'

export interface ReminderQuery {
  page?: number
  page_size?: number
  type?: string
  is_read?: boolean
}

export function listReminders(params?: ReminderQuery) {
  return request.get<PageResult<Reminder>>('/reminders', { params })
}

export function markReminderRead(id: number) {
  return request.put(`/reminders/${id}/read`)
}

export function markAllRemindersRead() {
  return request.put('/reminders/read-all')
}

export function deleteReminder(id: number) {
  return request.delete(`/reminders/${id}`)
}
