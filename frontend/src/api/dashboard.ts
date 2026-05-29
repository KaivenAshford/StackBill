import request from '@/utils/request'

export interface CategoryExpense {
  category_id: number
  category_name: string
  amount: number
  color: string
}

export interface DashboardData {
  monthly_expense: number
  yearly_expense: number
  subscription_count: number
  asset_count: number
  upcoming_renewals: number
  expiring_assets: number
  warning_assets: number
  recent_subscriptions: any[]
  recent_assets: any[]
  upcoming_renewal_list: any[]
  expiring_asset_list: any[]
  category_expense: CategoryExpense[]
}

export function getDashboard() {
  return request.get<DashboardData>('/dashboard')
}
