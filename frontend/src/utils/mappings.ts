import { useI18n } from 'vue-i18n'

export function useSubscriptionLabels() {
  const { t } = useI18n()

  const cycleLabels: Record<string, string> = {
    weekly: t('subscription.weekly'),
    monthly: t('subscription.monthly'),
    quarterly: t('subscription.quarterly'),
    yearly: t('subscription.yearly'),
    one_time: t('subscription.oneTime'),
  }

  const statusLabels: Record<string, string> = {
    active: t('subscription.active'),
    paused: t('subscription.paused'),
    cancelled: t('subscription.cancelled'),
    expired: t('subscription.expired'),
  }

  return { cycleLabels, statusLabels }
}

export function useAssetLabels() {
  const { t } = useI18n()

  const typeLabels: Record<string, string> = {
    domain: t('asset.domain'),
    server: t('asset.server'),
    docker_service: t('asset.dockerService'),
    ssl_certificate: t('asset.sslCertificate'),
    api_key: t('asset.apiKey'),
    repository: t('asset.repository'),
    other: t('asset.other'),
  }

  const statusLabels: Record<string, string> = {
    active: t('asset.active'),
    inactive: t('asset.inactive'),
    expired: t('asset.expired'),
    warning: t('asset.warning'),
  }

  return { typeLabels, statusLabels }
}

export function useCategoryLabels() {
  const { t } = useI18n()

  const typeLabels: Record<string, string> = {
    subscription: t('category.subscription'),
    asset: t('category.asset'),
  }

  return { typeLabels }
}
