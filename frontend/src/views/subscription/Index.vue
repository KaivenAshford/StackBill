<template>
  <div>
    <div class="page-toolbar">
      <h2 class="page-heading">{{ t('nav.subscriptions') }}</h2>
      <n-button type="primary" @click="$router.push('/subscriptions/new')">
        <template #icon>
          <Plus :size="16" :stroke-width="1.5" />
        </template>
        {{ t('common.create') }}
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
import { NDataTable, NTag, NButton, NSpace, NTooltip, useMessage, useDialog } from 'naive-ui'
import { Plus, Pencil, Trash2 } from '@lucide/vue'
import { deleteSubscription } from '@/api/subscription'
import { formatAmount } from '@/utils/currency'
import { useSubscriptionLabels } from '@/utils/mappings'
import type { Subscription } from '@/types'
import { useSubscriptionStore } from '@/stores/subscription'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const subscriptionStore = useSubscriptionStore()
const { cycleLabel, statusLabel, statusType } = useSubscriptionLabels()
const loading = computed(() => !subscriptionStore.loaded)
const page = ref(1)
const pagination = reactive({ page: 1, pageSize: 20, get itemCount() { return subscriptionStore.total }, showSizePicker: false, onChange: (p: number) => { page.value = p; fetchData() } })

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
    title: t('subscription.cycle'),
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

onMounted(() => subscriptionStore.ensureLoaded())

async function fetchData() {
  await subscriptionStore.refresh(page.value)
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
    await subscriptionStore.refresh(page.value)
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>

<style scoped>
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
