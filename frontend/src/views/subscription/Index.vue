<template>
  <div>
    <div style="display:flex;justify-content:space-between;margin-bottom:16px;">
      <h2 style="margin:0">{{ t('nav.subscriptions') }}</h2>
      <n-button type="primary" @click="$router.push('/subscriptions/new')">{{ t('common.create') }}</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :bordered="false" :pagination="pagination" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NDataTable, NTag } from 'naive-ui'
import { listSubscriptions } from '@/api/subscription'
import type { Subscription } from '@/types'

const { t } = useI18n()
const router = useRouter()
const items = ref<Subscription[]>([])
const page = ref(1)
const pagination = reactive({ page: 1, pageSize: 20, itemCount: 0, showSizePicker: false, onChange: (p: number) => { page.value = p; fetchData() } })

const columns = [
  { title: t('subscription.name'), key: 'name', render: (row: Subscription) => h('a', { style: 'cursor:pointer;color:#1890ff', onClick: () => router.push(`/subscriptions/${row.id}`) }, row.name) },
  { title: t('subscription.amount'), key: 'amount', render: (row: Subscription) => `${row.amount} ${row.currency}` },
  { title: t('subscription.cycle'), key: 'billing_cycle' },
  { title: t('subscription.status'), key: 'status', render: (row: Subscription) => h(NTag, { size: 'small', type: row.status === 'active' ? 'success' : 'default' }, { default: () => row.status }) },
  { title: t('subscription.nextPayment'), key: 'next_payment_date' },
]

onMounted(() => fetchData())

async function fetchData() {
  const res = await listSubscriptions({ page: page.value, page_size: 20 })
  items.value = (res.data as any).items
  pagination.itemCount = (res.data as any).total
  pagination.page = page.value
}
</script>
