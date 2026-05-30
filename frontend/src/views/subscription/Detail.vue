<template>
  <div>
    <n-page-header @back="$router.back()" :title="sub?.name || ''">
      <template #extra>
        <n-space>
          <n-button @click="$router.push(`/subscriptions/${id}/edit`)">{{ t('common.edit') }}</n-button>
          <n-button type="error" @click="confirmDelete">{{ t('common.delete') }}</n-button>
        </n-space>
      </template>
    </n-page-header>
    <n-descriptions bordered :column="2" style="margin-top:16px;" v-if="sub">
      <n-descriptions-item :label="t('subscription.amount')">{{ formatAmount(sub.amount, sub.currency) }}</n-descriptions-item>
      <n-descriptions-item :label="t('subscription.cycle')">{{ cycleLabel(sub.billing_cycle) }}</n-descriptions-item>
      <n-descriptions-item :label="t('subscription.status')">
        <n-tag :type="statusType(sub.status)" size="small">{{ statusLabel(sub.status) }}</n-tag>
      </n-descriptions-item>
      <n-descriptions-item :label="t('subscription.nextPayment')">{{ sub.next_payment_date || '-' }}</n-descriptions-item>
      <n-descriptions-item :label="t('subscription.url')">{{ sub.website_url || '-' }}</n-descriptions-item>
      <n-descriptions-item :label="t('subscription.remark')" :span="2">{{ sub.remark || '-' }}</n-descriptions-item>
    </n-descriptions>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NPageHeader, NDescriptions, NDescriptionsItem, NButton, NTag, NSpace, useMessage, useDialog } from 'naive-ui'
import { getSubscription, deleteSubscription } from '@/api/subscription'
import { formatAmount } from '@/utils/currency'
import type { Subscription } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const id = Number(route.params.id)
const sub = ref<Subscription | null>(null)

const cycleMap: Record<string, string> = { weekly: 'subscription.weekly', monthly: 'subscription.monthly', quarterly: 'subscription.quarterly', yearly: 'subscription.yearly', one_time: 'subscription.oneTime', custom: 'subscription.cycle' }
const statusMap: Record<string, string> = { active: 'subscription.active', paused: 'subscription.paused', cancelled: 'subscription.cancelled', expired: 'subscription.expired' }
const statusTypeMap: Record<string, string> = { active: 'success', paused: 'warning', cancelled: 'default', expired: 'error' }

function cycleLabel(v: string) { return t(cycleMap[v] || v) }
function statusLabel(v: string) { return t(statusMap[v] || v) }
function statusType(v: string) { return statusTypeMap[v] || 'default' }

onMounted(async () => {
  const res = await getSubscription(id)
  sub.value = res.data
})

function confirmDelete() {
  dialog.warning({
    title: t('common.confirm'),
    content: t('common.confirmDelete'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: handleDelete,
  })
}

async function handleDelete() {
  try {
    await deleteSubscription(id)
    message.success(t('common.success'))
    router.push('/subscriptions')
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>
