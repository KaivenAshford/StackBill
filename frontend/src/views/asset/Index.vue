<template>
  <div>
    <div style="display:flex;justify-content:space-between;margin-bottom:16px;">
      <h2 style="margin:0">{{ t('nav.assets') }}</h2>
      <n-button type="primary" @click="$router.push('/assets/new')">{{ t('common.create') }}</n-button>
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
import { listAssets, deleteAsset } from '@/api/asset'
import { formatAmount } from '@/utils/currency'
import type { Asset } from '@/types'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const items = ref<Asset[]>([])
const loading = ref(true)
const page = ref(1)
const pagination = reactive({ page: 1, pageSize: 20, itemCount: 0, showSizePicker: false, onChange: (p: number) => { page.value = p; fetchData() } })

const typeMap: Record<string, string> = { domain: 'asset.domain', server: 'asset.server', docker_service: 'asset.dockerService', ssl_certificate: 'asset.sslCertificate', api_key: 'asset.apiKey', repository: 'asset.repository', other: 'asset.other' }
const statusMap: Record<string, string> = { active: 'asset.active', inactive: 'asset.inactive', expired: 'asset.expired', warning: 'asset.warning' }
const statusTypeMap: Record<string, string> = { active: 'success', inactive: 'default', expired: 'error', warning: 'warning' }

function typeLabel(v: string) { return t(typeMap[v] || v) }
function statusLabel(v: string) { return t(statusMap[v] || v) }
function statusType(v: string) { return statusTypeMap[v] || 'default' }

const columns = [
  { title: t('asset.name'), key: 'name', render: (row: Asset) => h('a', { style: 'cursor:pointer;color:#1890ff', onClick: () => router.push(`/assets/${row.id}`) }, row.name) },
  { title: t('asset.type'), key: 'asset_type', render: (row: Asset) => typeLabel(row.asset_type) },
  { title: t('asset.provider'), key: 'provider' },
  { title: t('asset.costAmount'), key: 'cost_amount', render: (row: Asset) => row.cost_amount ? formatAmount(row.cost_amount, row.cost_currency) : '-' },
  { title: t('asset.expireDate'), key: 'expire_date' },
  { title: t('asset.status'), key: 'status', render: (row: Asset) => h(NTag, { size: 'small', type: statusType(row.status) }, { default: () => statusLabel(row.status) }) },
  {
    title: t('common.edit'), key: 'actions', width: 120,
    render: (row: Asset) => h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => router.push(`/assets/${row.id}/edit`) }, { default: () => t('common.edit') }),
        h(NButton, { size: 'small', type: 'error', onClick: () => confirmDelete(row.id) }, { default: () => t('common.delete') }),
      ],
    }),
  },
]

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const res = await listAssets({ page: page.value, page_size: 20 })
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
    await deleteAsset(id)
    message.success(t('common.success'))
    await fetchData()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>
