import { describe, it, expect } from 'vitest'
import zhCN from '@/locales/zh-CN'
import enUS from '@/locales/en-US'

describe('i18n locale files', () => {
  it('zh-CN has all required keys', () => {
    expect(zhCN.common?.confirm).toBe('确认')
    expect(zhCN.common?.search).toBe('搜索')
    expect(zhCN.common?.clearFilter).toBe('清除筛选')
    expect(zhCN.common?.allStatus).toBe('全部状态')
    expect(zhCN.common?.allCategories).toBe('全部分类')
    expect(zhCN.common?.allTypes).toBe('全部类型')

    expect(zhCN.subscription?.paymentMethod).toBe('付款方式')
    expect(zhCN.subscription?.autoRenew).toBe('自动续费')
    expect(zhCN.subscription?.category).toBe('分类')

    expect(zhCN.asset?.identifier).toBe('标识符')
    expect(zhCN.asset?.description).toBe('描述')
    expect(zhCN.asset?.linkedSubscription).toBe('关联订阅')

    expect(zhCN.errors?.invalidParams).toBe('参数校验失败')
    expect(zhCN.errors?.duplicateUsername).toBe('用户名已存在')
    expect(zhCN.errors?.requestFailed).toBe('请求失败')
  })

  it('en-US has all required keys', () => {
    expect(enUS.common?.confirm).toBe('Confirm')
    expect(enUS.common?.clearFilter).toBe('Clear Filters')

    expect(enUS.subscription?.paymentMethod).toBe('Payment Method')
    expect(enUS.subscription?.autoRenew).toBe('Auto Renew')

    expect(enUS.asset?.linkedSubscription).toBe('Linked Subscription')

    expect(enUS.errors?.invalidParams).toBe('Invalid parameters')
    expect(enUS.errors?.requestFailed).toBe('Request failed')
  })

  it('both locales have the same top-level keys', () => {
    const zh = Object.keys(zhCN)
    const en = Object.keys(enUS)
    expect(zh.sort()).toEqual(en.sort())
  })

  it('errors section has matching keys across locales', () => {
    const zhErrKeys = Object.keys(zhCN.errors || {})
    const enErrKeys = Object.keys(enUS.errors || {})
    expect(zhErrKeys.sort()).toEqual(enErrKeys.sort())
  })
})
