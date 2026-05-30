<template>
  <div>
    <div style="display:flex;justify-content:space-between;margin-bottom:16px;">
      <h2 style="margin:0">{{ t('nav.subscriptions') }}</h2>
      <n-button type="primary" @click="$router.push('/subscriptions/new')">{{ t('common.create') }}</n-button>
    </div>
    <n-spin :show="loading">
      <n-data-table :columns="columns" :data="items" :bordered="false" :pagination="pagination" />
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NDataTable, NTag, NButton, NSpace, NSpin, useMessage, useDialog } from 'naive-ui'
import { listSubscriptions, deleteSubscription } from '@/api/subscription'
import { formatAmount } from '@/utils/currency'
import type { Subscription } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const items = ref<Subscription[]>([])
const loading = ref(true)
const page = ref(1)
const pagination = reactive({ page: 1, pageSize: 20, itemCount: 0, showSizePicker: false, onChange: (p: number) => { page.value = p; fetchData() } })

const cycleMap: Record<string, string> = { weekly: 'subscription.weekly', monthly: 'subscription.monthly', quarterly: 'subscription.quarterly', yearly: 'subscription.yearly', one_time: 'subscription.oneTime', custom: 'subscription.cycle' }
const statusMap: Record<string, string> = { active: 'subscription.active', paused: 'subscription.paused', cancelled: 'subscription.cancelled', expired: 'subscription.expired' }
const statusTypeMap: Record<string, string> = { active: 'success', paused: 'warning', cancelled: 'default', expired: 'error' }

function cycleLabel(v: string) { return t(cycleMap[v] || v) }
function statusLabel(v: string) { return t(statusMap[v] || v) }
function statusType(v: string) { return statusTypeMap[v] || 'default' }

const columns = [
  { title: t('subscription.name'), key: 'name', render: (row: Subscription) => h('a', { style: 'cursor:pointer;color:#1890ff', onClick: () => router.push(`/subscriptions/${row.id}`) }, row.name) },
  { title: t('subscription.amount'), key: 'amount', render: (row: Subscription) => formatAmount(row.amount, row.currency) },
  { title: t('subscription.cycle'), key: 'billing_cycle', render: (row: Subscription) => cycleLabel(row.billing_cycle) },
  { title: t('subscription.status'), key: 'status', render: (row: Subscription) => h(NTag, { size: 'small', type: statusType(row.status) }, { default: () => statusLabel(row.status) }) },
  { title: t('subscription.nextPayment'), key: 'next_payment_date' },
  {
    title: t('common.edit'), key: 'actions', width: 120,
    render: (row: Subscription) => h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => router.push(`/subscriptions/${row.id}/edit`) }, { default: () => t('common.edit') }),
        h(NButton, { size: 'small', type: 'error', onClick: () => confirmDelete(row.id) }, { default: () => t('common.delete') }),
      ],
    }),
  },
]

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const res = await listSubscriptions({ page: page.value, page_size: 20 })
    items.value = (res.data as any).items
    pagination.itemCount = (res.data as any).total
    pagination.page = page.value
  } finally {
    loading.value = false
  }
}

function confirmDelete(id: number) {
  dialog.warning({
    title: t('common.confirm'),
    content: t('common.confirmDelete'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => handleDelete(id),
  })
}

async function handleDelete(id: number) {
  try {
    await deleteSubscription(id)
    message.success(t('common.success'))
    await fetchData()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>
