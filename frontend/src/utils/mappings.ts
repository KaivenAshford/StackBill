import { useI18n } from 'vue-i18n'

type TagType = 'info' | 'success' | 'warning' | 'error' | 'default' | 'primary'

/**
 * Subscription label helpers — reactive with i18n locale changes.
 */
export function useSubscriptionLabels() {
  const { t } = useI18n()

  const cycleMap: Record<string, string> = {
    weekly: 'subscription.weekly',
    monthly: 'subscription.monthly',
    quarterly: 'subscription.quarterly',
    yearly: 'subscription.yearly',
    one_time: 'subscription.oneTime',
    custom: 'subscription.customCycle',
  }

  const statusMap: Record<string, string> = {
    active: 'subscription.active',
    paused: 'subscription.paused',
    cancelled: 'subscription.cancelled',
    expired: 'subscription.expired',
  }

  const statusTypeMap: Record<string, TagType> = {
    active: 'success',
    paused: 'warning',
    cancelled: 'default',
    expired: 'error',
  }

  function cycleLabel(v: string) { return t(cycleMap[v] || v) }
  function statusLabel(v: string) { return t(statusMap[v] || v) }
  function statusType(v: string): TagType { return statusTypeMap[v] || 'default' }

  return { cycleMap, statusMap, statusTypeMap, cycleLabel, statusLabel, statusType }
}

/**
 * Asset label helpers — reactive with i18n locale changes.
 */
export function useAssetLabels() {
  const { t } = useI18n()

  const typeMap: Record<string, string> = {
    domain: 'asset.domain',
    server: 'asset.server',
    docker_service: 'asset.dockerService',
    ssl_certificate: 'asset.sslCertificate',
    api_key: 'asset.apiKey',
    repository: 'asset.repository',
    other: 'asset.other',
  }

  const statusMap: Record<string, string> = {
    active: 'asset.active',
    inactive: 'asset.inactive',
    expired: 'asset.expired',
    warning: 'asset.warning',
  }

  const statusTypeMap: Record<string, TagType> = {
    active: 'success',
    inactive: 'default',
    expired: 'error',
    warning: 'warning',
  }

  function typeLabel(v: string) { return t(typeMap[v] || v) }
  function statusLabel(v: string) { return t(statusMap[v] || v) }
  function statusType(v: string): TagType { return statusTypeMap[v] || 'default' }

  return { typeMap, statusMap, statusTypeMap, typeLabel, statusLabel, statusType }
}

/**
 * Category label helpers — reactive with i18n locale changes.
 */
export function useCategoryLabels() {
  const { t } = useI18n()

  const typeLabels: Record<string, string> = {
    subscription: t('category.subscription'),
    asset: t('category.asset'),
  }

  return { typeLabels }
}
