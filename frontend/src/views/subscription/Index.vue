<template>
  <div>
    <div class="page-toolbar">
      <h2 class="page-heading">{{ t('nav.subscriptions') }}</h2>
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
        <n-button type="primary" @click="$router.push('/subscriptions/new')">
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
        v-model:value="filters.status"
        :options="statusFilterOptions"
        :placeholder="t('common.allStatus')"
        clearable
        class="filter-select"
        @update:value="fetchFiltered"
      />
      <n-select
        v-model:value="filters.category_id"
        :options="categoryFilterOptions"
        :placeholder="t('common.allCategories')"
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
          <div class="skel-cell" style="width:80px"></div>
          <div class="skel-cell" style="width:60px"></div>
          <div class="skel-cell" style="width:50px"></div>
          <div class="skel-cell" style="width:90px"></div>
          <div class="skel-cell" style="width:60px"></div>
        </div>
        <div v-for="i in 5" :key="i" class="skeleton-row">
          <div class="skel-cell" style="width:100px"></div>
          <div class="skel-cell" style="width:70px"></div>
          <div class="skel-cell" style="width:50px"></div>
          <div class="skel-cell skel-tag" style="width:40px"></div>
          <div class="skel-cell" style="width:80px"></div>
          <div class="skel-cell" style="width:50px"></div>
        </div>
      </div>
      <div v-show="!loading">
        <n-data-table :columns="columns" :data="subscriptionStore.subscriptions" :bordered="false" :pagination="pagination" :scroll-x="680" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NDataTable, NTag, NButton, NSpace, NTooltip, NInput, NSelect, NUpload, useMessage, useDialog } from 'naive-ui'
import type { UploadCustomRequestOptions } from 'naive-ui'
import { Plus, Pencil, Trash2, Search, Download, Upload } from '@lucide/vue'
import { deleteSubscription } from '@/api/subscription'
import request from '@/utils/request'
import { formatAmount } from '@/utils/currency'
import { useSubscriptionLabels } from '@/utils/mappings'
import type { Subscription } from '@/types'
import type { SubscriptionQuery } from '@/api/subscription'
import { useSubscriptionStore } from '@/stores/subscription'
import { useCategoryStore } from '@/stores/category'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const subscriptionStore = useSubscriptionStore()
const categoryStore = useCategoryStore()
const { cycleLabel, statusLabel, statusType } = useSubscriptionLabels()
const loading = computed(() => !subscriptionStore.loaded)
const page = ref(1)
const pagination = reactive({ page: 1, pageSize: 20, get itemCount() { return subscriptionStore.total }, showSizePicker: false, onChange: (p: number) => { page.value = p; fetchData() } })

const filters = reactive({
  keyword: '',
  status: null as string | null,
  category_id: null as number | null,
})

const hasActiveFilters = computed(() => filters.keyword || filters.status || filters.category_id)

const statusFilterOptions = [
  { label: () => t('subscription.active'), value: 'active' },
  { label: () => t('subscription.paused'), value: 'paused' },
  { label: () => t('subscription.cancelled'), value: 'cancelled' },
  { label: () => t('subscription.expired'), value: 'expired' },
]

const categoryFilterOptions = computed(() =>
  categoryStore.categories
    .filter(c => c.type === 'subscription')
    .map(c => ({ label: c.name, value: c.id }))
)

let debounceTimer: ReturnType<typeof setTimeout> | null = null
function debouncedFetch() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => fetchFiltered(), 300)
}

function buildQuery(): SubscriptionQuery {
  const query: SubscriptionQuery = { page: page.value, page_size: 20 }
  if (filters.keyword) query.keyword = filters.keyword
  if (filters.status) query.status = filters.status
  if (filters.category_id) query.category_id = filters.category_id
  return query
}

function clearFilters() {
  filters.keyword = ''
  filters.status = null
  filters.category_id = null
  page.value = 1
  fetchData()
}

async function fetchFiltered() {
  page.value = 1
  await fetchData()
}

const columns = [
  {
    title: t('subscription.name'),
    key: 'name',
    render: (row: Subscription) => h('a', {
      class: 'table-link',
      onClick: () => router.push(`/subscriptions/${row.id}`),
    }, row.name),
  },
  {
    title: t('subscription.amount'),
    key: 'amount',
    render: (row: Subscription) => h('span', { style: 'font-family: var(--font-heading); font-weight: 600;' }, formatAmount(row.amount, row.currency)),
  },
  {
    title: t('subscription.billingCycle'),
    key: 'billing_cycle',
    render: (row: Subscription) => cycleLabel(row.billing_cycle),
  },
  {
    title: t('subscription.status'),
    key: 'status',
    render: (row: Subscription) => h(NTag, { size: 'small', type: statusType(row.status), round: true }, { default: () => statusLabel(row.status) }),
  },
  { title: t('subscription.nextPayment'), key: 'next_payment_date' },
  {
    title: t('common.edit'),
    key: 'actions',
    width: 100,
    align: 'center' as const,
    render: (row: Subscription) => h(NSpace, { size: 'small', justify: 'center' }, {
      default: () => [
        h(NTooltip, null, {
          trigger: () => h(NButton, { size: 'small', quaternary: true, 'aria-label': t('common.edit'), onClick: () => router.push(`/subscriptions/${row.id}/edit`) }, {
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

onMounted(async () => {
  await categoryStore.ensureLoaded('subscription')
  await subscriptionStore.ensureLoaded()
})

async function fetchData() {
  await subscriptionStore.refresh(buildQuery())
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
    await deleteSubscription(id)
    message.success(t('common.success'))
    await fetchData()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}

async function handleExport() {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/subscriptions/export', {
      headers: { Authorization: `Bearer ${token}` },
    })
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'subscriptions.csv'
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    message.error(t('common.failed'))
  }
}

async function handleImport({ file }: UploadCustomRequestOptions) {
  if (!file.file) return
  try {
    const formData = new FormData()
    formData.append('file', file.file)
    await request.post('/subscriptions/import', formData)
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
  grid-template-columns: 2fr 1fr 1fr 1fr 1.2fr 80px;
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
