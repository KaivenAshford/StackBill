<template>
  <div class="detail-page">
    <n-page-header @back="$router.back()" :title="sub?.name || ''">
      <template #extra>
        <n-space v-if="sub">
          <n-button @click="$router.push(`/subscriptions/${id}/edit`)">
            <template #icon>
              <Pencil :size="14" :stroke-width="1.5" />
            </template>
            {{ t('common.edit') }}
          </n-button>
          <n-button type="error" @click="confirmDelete">
            <template #icon>
              <Trash2 :size="14" :stroke-width="1.5" />
            </template>
            {{ t('common.delete') }}
          </n-button>
        </n-space>
      </template>
    </n-page-header>

    <!-- Loading state -->
    <div v-if="loading" class="detail-card section-gap">
      <div v-for="i in 5" :key="i" class="detail-row">
        <div class="skeleton-shimmer" style="width:80px;height:12px;border-radius:4px;"></div>
        <div class="skeleton-shimmer" style="width:120px;height:14px;border-radius:4px;"></div>
      </div>
    </div>

    <!-- Error state -->
    <n-result v-else-if="error" status="error" :title="t('common.failed')" :description="error" class="section-gap">
      <template #footer>
        <n-button @click="fetchData">{{ t('common.retry') || 'Retry' }}</n-button>
      </template>
    </n-result>

    <!-- Data -->
    <div v-else-if="sub" class="detail-grid section-gap">
      <div class="detail-card">
        <div class="detail-row">
          <span class="detail-label">{{ t('subscription.amount') }}</span>
          <span class="detail-value detail-value--amount">{{ formatAmount(sub.amount, sub.currency) }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('subscription.cycle') }}</span>
          <span class="detail-value">{{ cycleLabel(sub.billing_cycle) }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('subscription.status') }}</span>
          <n-tag :type="statusType(sub.status)" size="small" round>{{ statusLabel(sub.status) }}</n-tag>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('subscription.nextPayment') }}</span>
          <span class="detail-value">{{ sub.next_payment_date || '-' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('subscription.url') }}</span>
          <span class="detail-value">
            <a v-if="sub.website_url" :href="sub.website_url" target="_blank" rel="noopener noreferrer" class="detail-link">{{ sub.website_url }}</a>
            <template v-else>-</template>
          </span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('subscription.remark') }}</span>
          <span class="detail-value">{{ sub.remark || '-' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NPageHeader, NButton, NTag, NSpace, NResult, useMessage, useDialog } from 'naive-ui'
import { Pencil, Trash2 } from '@lucide/vue'
import { getSubscription, deleteSubscription } from '@/api/subscription'
import { formatAmount } from '@/utils/currency'
import { useSubscriptionLabels } from '@/utils/mappings'
import type { Subscription } from '@/types'
import { useSubscriptionStore } from '@/stores/subscription'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const { cycleLabel, statusLabel, statusType } = useSubscriptionLabels()
const id = Number(route.params.id)
const sub = ref<Subscription | null>(null)
const loading = ref(true)
const error = ref('')

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  error.value = ''
  try {
    const res = await getSubscription(id)
    sub.value = res.data
  } catch (e: unknown) {
    error.value = (e as Error).message || t('common.failed')
  } finally {
    loading.value = false
  }
}

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
    useSubscriptionStore().invalidate()
    router.push('/subscriptions')
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>

<style scoped>
.detail-page {
  animation: fadeIn 0.3s ease-out;
}

.detail-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: border-color var(--transition-smooth);
}

.detail-card:hover {
  border-color: var(--color-border-strong);
}

.detail-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
  gap: var(--spacing-md);
  transition: background var(--transition-fast);
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-row:hover {
  background: var(--gradient-card-accent);
}

.detail-label {
  font-size: 14px;
  color: var(--color-text-secondary);
  flex-shrink: 0;
  font-weight: 500;
}

.detail-value {
  font-size: 14px;
  color: var(--color-text-primary);
  font-weight: 500;
  text-align: right;
  word-break: break-all;
}

.detail-value--amount {
  font-family: var(--font-heading);
  font-size: 18px;
  font-weight: 700;
  color: var(--color-accent);
}

.detail-link {
  color: var(--color-accent);
  word-break: break-all;
  transition: color var(--transition-fast);
}

.detail-link:hover {
  color: var(--color-accent-hover);
  text-decoration: underline;
}

@media (max-width: 768px) {
  .detail-row {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-xs);
  }
  .detail-value {
    text-align: left;
  }
}
</style>
