export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
}

export interface Subscription {
  id: number
  user_id: number
  name: string
  description: string
  category_id: number
  amount: number
  currency: string
  billing_cycle: string
  billing_interval: number
  next_payment_date: string | null
  start_date: string | null
  payment_method: string
  auto_renew: boolean
  status: string
  website_url: string
  remark: string
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  user_id: number
  name: string
  type: string
  color: string
  icon: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Asset {
  id: number
  user_id: number
  name: string
  asset_type: string
  provider: string
  identifier: string
  url: string
  expire_date: string | null
  cost_amount: number
  cost_currency: string
  billing_cycle: string
  status: string
  subscription_id: number
  description: string
  remark: string
  created_at: string
  updated_at: string
}

export interface Reminder {
  id: number
  user_id: number
  target_type: string
  target_id: number
  remind_type: string
  remind_date: string | null
  title: string
  content: string
  is_read: boolean
  amount?: number | null
  currency?: string
  expire_date?: string
  asset_status?: string
  created_at: string
  updated_at: string
}

export interface DashboardData {
  monthly_expense: number
  yearly_expense: number
  subscription_count: number
  asset_count: number
  upcoming_renewals: number
  expiring_assets: number
  warning_assets: number
  recent_subscriptions: Subscription[]
  recent_assets: Asset[]
  upcoming_renewal_list: Subscription[]
  expiring_asset_list: Asset[]
  category_expense: CategoryExpense[]
}

export interface CategoryExpense {
  category_id: number
  category_name: string
  amount: number
  color: string
}
