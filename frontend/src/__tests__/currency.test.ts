import { describe, it, expect } from 'vitest'
import { formatAmount, getCurrencySymbol, currencyOptions } from '@/utils/currency'

describe('getCurrencySymbol', () => {
  it('returns correct symbol for known currencies', () => {
    expect(getCurrencySymbol('CNY')).toBe('¥')
    expect(getCurrencySymbol('USD')).toBe('$')
    expect(getCurrencySymbol('EUR')).toBe('€')
    expect(getCurrencySymbol('GBP')).toBe('£')
    expect(getCurrencySymbol('JPY')).toBe('¥')
    expect(getCurrencySymbol('HKD')).toBe('HK$')
  })

  it('returns the code itself for unknown currencies', () => {
    expect(getCurrencySymbol('XYZ')).toBe('XYZ')
  })
})

describe('formatAmount', () => {
  it('formats amount with currency symbol', () => {
    expect(formatAmount(9.99, 'USD')).toBe('$9.99')
    expect(formatAmount(100, 'CNY')).toBe('¥100.00')
    expect(formatAmount(0, 'EUR')).toBe('€0.00')
  })

  it('formats large amounts correctly', () => {
    expect(formatAmount(1234567.89, 'USD')).toBe('$1234567.89')
  })

  it('formats unknown currency as code prefix', () => {
    expect(formatAmount(50, 'BTC')).toBe('BTC50.00')
  })
})

describe('currencyOptions', () => {
  it('contains expected currencies', () => {
    const values = currencyOptions.map(o => o.value)
    expect(values).toContain('CNY')
    expect(values).toContain('USD')
    expect(values).toContain('EUR')
    expect(values).toContain('JPY')
  })

  it('each option has label and value', () => {
    for (const opt of currencyOptions) {
      expect(opt.label).toBeTruthy()
      expect(opt.value).toBeTruthy()
    }
  })
})
