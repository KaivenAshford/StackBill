<template>
  <div>
    <div style="display:flex;justify-content:space-between;margin-bottom:16px;">
      <h2 style="margin:0">{{ t('nav.assets') }}</h2>
      <n-button type="primary" @click="$router.push('/assets/new')">{{ t('common.create') }}</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :bordered="false" :pagination="pagination" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NDataTable, NTag } from 'naive-ui'
import { listAssets } from '@/api/asset'
import type { Asset } from '@/types'

const { t } = useI18n()
const router = useRouter()
const items = ref<Asset[]>([])
const page = ref(1)
const pagination = reactive({ page: 1, pageSize: 20, itemCount: 0, showSizePicker: false, onChange: (p: number) => { page.value = p; fetchData() } })

const columns = [
  { title: t('asset.name'), key: 'name', render: (row: Asset) => h('a', { style: 'cursor:pointer;color:#1890ff', onClick: () => router.push(`/assets/${row.id}`) }, row.name) },
  { title: t('asset.type'), key: 'asset_type' },
  { title: t('asset.provider'), key: 'provider' },
  { title: t('asset.expireDate'), key: 'expire_date' },
  { title: t('asset.status'), key: 'status', render: (row: Asset) => h(NTag, { size: 'small', type: row.status === 'active' ? 'success' : row.status === 'warning' ? 'warning' : 'error' }, { default: () => row.status }) },
]

onMounted(() => fetchData())

async function fetchData() {
  const res = await listAssets({ page: page.value, page_size: 20 })
  items.value = (res.data as any).items
  pagination.itemCount = (res.data as any).total
  pagination.page = page.value
}
</script>
