export interface CurrencyOption {
  label: string
  value: string
  symbol: string
}

const currencies: CurrencyOption[] = [
  { label: 'CNY', value: 'CNY', symbol: '¥' },
  { label: 'USD', value: 'USD', symbol: '$' },
  { label: 'EUR', value: 'EUR', symbol: '€' },
  { label: 'GBP', value: 'GBP', symbol: '£' },
  { label: 'JPY', value: 'JPY', symbol: '¥' },
  { label: 'HKD', value: 'HKD', symbol: 'HK$' },
  { label: 'TWD', value: 'TWD', symbol: 'NT$' },
  { label: 'KRW', value: 'KRW', symbol: '₩' },
  { label: 'SGD', value: 'SGD', symbol: 'S$' },
  { label: 'AUD', value: 'AUD', symbol: 'A$' },
  { label: 'CAD', value: 'CAD', symbol: 'C$' },
]

const symbolMap = Object.fromEntries(currencies.map(c => [c.value, c.symbol]))

export function getCurrencySymbol(code: string): string {
  return symbolMap[code] || code
}

export function formatAmount(amount: number, currency: string): string {
  return `${getCurrencySymbol(currency)}${amount.toFixed(2)}`
}

export const currencyOptions = currencies.map(c => ({ label: `${c.symbol} ${c.label}`, value: c.value }))
