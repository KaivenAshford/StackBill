<template>
  <div>
    <div class="page-toolbar">
      <h2 class="page-heading">{{ t('nav.assets') }}</h2>
      <n-space>
        <n-button @click="handleExport">
          <template #icon><Download :size="16" :stroke-width="1.5" /></template>
          {{ t('common.exportCSV') }}
        </n-button>
        <n-upload :show-file-list="false" accept=".csv" :custom-request="handleImport">
          <n-button>
            <template #icon><Upload :size="16" :stroke-width="1.5" /></template>
            {{ t('common.importCSV') }}
          </n-button>
        </n-upload>
        <n-button type="primary" @click="$router.push('/assets/new')">
          <template #icon><Plus :size="16" :stroke-width="1.5" /></template>
          {{ t('common.create') }}
        </n-button>
      </n-space>
    </div>

    <div class="filter-bar">
      <n-input
        v-model:value="filters.keyword"
        :placeholder="t('common.search')"
        clearable
        class="filter-search"
        @update:value="debouncedFetch"
      >
        <template #prefix>
          <Search :size="14" />
        </template>
      </n-input>
      <n-select
        v-model:value="filters.asset_type"
        :options="typeFilterOptions"
        :placeholder="t('common.allTypes')"
        clearable
        class="filter-select"
        @update:value="fetchFiltered"
      />
      <n-select
        v-model:value="filters.status"
        :options="statusFilterOptions"
        :placeholder="t('common.allStatus')"
        clearable
        class="filter-select"
        @update:value="fetchFiltered"
      />
      <n-button v-if="hasActiveFilters" quaternary size="small" @click="clearFilters">
        {{ t('common.clearFilter') }}
      </n-button>
    </div>

    <div class="table-card">
      <div v-if="loading" class="table-skeleton">
        <div class="skeleton-row skeleton-row--head">
          <div class="skel-cell" style="width:120px"></div>
          <div class="skel-cell skel-tag" style="width:50px"></div>
          <div class="skel-cell" style="width:80px"></div>
          <div class="skel-cell" style="width:70px"></div>
          <div class="skel-cell" style="width:90px"></div>
          <div class="skel-cell skel-tag" style="width:50px"></div>
          <div class="skel-cell" style="width:60px"></div>
        </div>
        <div v-for="i in 5" :key="i" class="skeleton-row">
          <div class="skel-cell" style="width:100px"></div>
          <div class="skel-cell skel-tag" style="width:40px"></div>
          <div class="skel-cell" style="width:70px"></div>
          <div class="skel-cell" style="width:60px"></div>
          <div class="skel-cell" style="width:80px"></div>
          <div class="skel-cell skel-tag" style="width:40px"></div>
          <div class="skel-cell" style="width:50px"></div>
        </div>
      </div>
      <div v-show="!loading">
        <n-data-table :columns="columns" :data="assetStore.assets" :bordered="false" :pagination="pagination" :scroll-x="760" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NDataTable, NTag, NButton, NSpace, NTooltip, NInput, NSelect, NUpload, useMessage, useDialog } from 'naive-ui'
import { Plus, Pencil, Trash2, Search, Download, Upload } from '@lucide/vue'
import { deleteAsset } from '@/api/asset'
import request from '@/utils/request'
import { formatAmount } from '@/utils/currency'
import { useAssetLabels } from '@/utils/mappings'
import type { Asset } from '@/types'
import type { AssetQuery } from '@/api/asset'
import { useAssetStore } from '@/stores/asset'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const assetStore = useAssetStore()
const { typeLabel, statusLabel, statusType } = useAssetLabels()
const loading = computed(() => !assetStore.loaded)
const page = ref(1)
const pagination = reactive({ page: 1, pageSize: 20, get itemCount() { return assetStore.total }, showSizePicker: false, onChange: (p: number) => { page.value = p; fetchData() } })

const filters = reactive({
  keyword: '',
  asset_type: null as string | null,
  status: null as string | null,
})

const hasActiveFilters = computed(() => filters.keyword || filters.asset_type || filters.status)

const typeFilterOptions = [
  { label: () => t('asset.domain'), value: 'domain' },
  { label: () => t('asset.server'), value: 'server' },
  { label: () => t('asset.dockerService'), value: 'docker_service' },
  { label: () => t('asset.sslCertificate'), value: 'ssl_certificate' },
  { label: () => t('asset.apiKey'), value: 'api_key' },
  { label: () => t('asset.repository'), value: 'repository' },
  { label: () => t('asset.other'), value: 'other' },
]

const statusFilterOptions = [
  { label: () => t('asset.active'), value: 'active' },
  { label: () => t('asset.inactive'), value: 'inactive' },
  { label: () => t('asset.expired'), value: 'expired' },
  { label: () => t('asset.warning'), value: 'warning' },
]

let debounceTimer: ReturnType<typeof setTimeout> | null = null
function debouncedFetch() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => fetchFiltered(), 300)
}

function buildQuery(): AssetQuery {
  const query: AssetQuery = { page: page.value, page_size: 20 }
  if (filters.keyword) query.keyword = filters.keyword
  if (filters.asset_type) query.asset_type = filters.asset_type
  if (filters.status) query.status = filters.status
  return query
}

function clearFilters() {
  filters.keyword = ''
  filters.asset_type = null
  filters.status = null
  page.value = 1
  fetchData()
}

async function fetchFiltered() {
  page.value = 1
  await fetchData()
}

const columns = [
  {
    title: t('asset.name'),
    key: 'name',
    render: (row: Asset) => h('a', {
      class: 'table-link',
      onClick: () => router.push(`/assets/${row.id}`),
    }, row.name),
  },
  {
    title: t('asset.type'),
    key: 'asset_type',
    render: (row: Asset) => h(NTag, { size: 'small', round: true }, { default: () => typeLabel(row.asset_type) }),
  },
  { title: t('asset.provider'), key: 'provider' },
  {
    title: t('asset.costAmount'),
    key: 'cost_amount',
    render: (row: Asset) => row.cost_amount
      ? h('span', { style: 'font-family: var(--font-heading); font-weight: 600;' }, formatAmount(row.cost_amount, row.cost_currency))
      : '-',
  },
  { title: t('asset.expireDate'), key: 'expire_date' },
  {
    title: t('asset.status'),
    key: 'status',
    render: (row: Asset) => h(NTag, { size: 'small', type: statusType(row.status), round: true }, { default: () => statusLabel(row.status) }),
  },
  {
    title: t('common.edit'),
    key: 'actions',
    width: 100,
    align: 'center' as const,
    render: (row: Asset) => h(NSpace, { size: 'small', justify: 'center' }, {
      default: () => [
        h(NTooltip, null, {
          trigger: () => h(NButton, { size: 'small', quaternary: true, 'aria-label': t('common.edit'), onClick: () => router.push(`/assets/${row.id}/edit`) }, {
            icon: () => h(Pencil, { size: 14, strokeWidth: 1.5 }),
          }),
          default: () => t('common.edit'),
        }),
        h(NTooltip, null, {
          trigger: () => h(NButton, { size: 'small', quaternary: true, type: 'error', 'aria-label': t('common.delete'), onClick: () => confirmDelete(row.id) }, {
            icon: () => h(Trash2, { size: 14, strokeWidth: 1.5 }),
          }),
          default: () => t('common.delete'),
        }),
      ],
    }),
  },
]

onMounted(() => assetStore.ensureLoaded())

async function fetchData() {
  await assetStore.refresh(buildQuery())
  pagination.page = page.value
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

async function handleExport() {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/assets/export', {
      headers: { Authorization: `Bearer ${token}` },
    })
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'assets.csv'
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    message.error(t('common.failed'))
  }
}

async function handleImport({ file }: { file: { file: File } }) {
  try {
    const formData = new FormData()
    formData.append('file', file.file)
    await request.post('/assets/import', formData)
    message.success(t('common.success'))
    await fetchData()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>

<style scoped>
.filter-bar {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-md);
  flex-wrap: wrap;
}

.filter-search {
  width: 200px;
}

.filter-select {
  width: 150px;
}

.table-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: border-color var(--transition-smooth);
  min-height: 320px;
}

.table-card:hover {
  border-color: var(--color-border-strong);
}

:deep(.table-link) {
  color: var(--color-accent);
  font-weight: 500;
  cursor: pointer;
  transition: color var(--transition-fast);
}
:deep(.table-link:hover) {
  color: var(--color-accent-hover);
}

.table-skeleton {
  padding: var(--spacing-md);
}

.skeleton-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1.2fr 1fr 80px;
  gap: var(--spacing-md);
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-md);
  min-width: 0;
}

@media (max-width: 768px) {
  .skeleton-row {
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-sm);
  }
  .skeleton-row--head {
    display: none;
  }
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }
  .filter-search,
  .filter-select {
    width: 100%;
  }
}

.skeleton-row--head {
  opacity: 0.4;
  border-bottom: 1px solid var(--color-border);
  padding-bottom: var(--spacing-md);
  margin-bottom: var(--spacing-xs);
}

.skel-cell {
  height: 14px;
  border-radius: 4px;
  background: var(--color-bg-muted);
  background-size: 200% 100%;
  animation: shimmer 1.5s ease-in-out infinite;
}

.skel-tag {
  height: 22px;
  border-radius: 10px;
  width: 40px !important;
}
</style>
