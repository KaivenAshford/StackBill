<template>
  <div>
    <div class="page-toolbar">
      <h2 class="page-heading">{{ t('nav.reminders') }}</h2>
      <n-button :type="hasUnread ? 'primary' : 'default'" @click="handleMarkAllRead" :disabled="!hasUnread">
        <template #icon>
          <CheckCheck :size="16" :stroke-width="1.5" />
        </template>
        {{ t('reminder.markAllRead') }}
      </n-button>
    </div>

    <!-- Single container — skeleton or real cards share the same list wrapper -->
    <div v-if="displayItems.length > 0" class="reminder-list">
      <div v-for="item in displayItems" :key="item._key" class="reminder-card" :class="{ 'reminder-card--skeleton': item._skeleton, 'reminder-card--unread': !item._skeleton && !(item as RealItem).is_read }">
        <div class="reminder-dot" v-if="!item._skeleton && !(item as RealItem).is_read"></div>
        <div class="reminder-icon">
          <component v-if="!item._skeleton" :is="remindIcon((item as RealItem).remind_type)" :size="18" :stroke-width="1.5" />
        </div>
        <div class="reminder-content">
          <div class="reminder-header">
            <span class="reminder-title" :class="{ 'reminder-title--unread': !item._skeleton && !(item as RealItem).is_read }">{{ item.title }}</span>
            <n-tag v-if="!item._skeleton" size="small" :type="(item as RealItem).remind_type === 'service_warning' ? 'warning' : 'info'" round class="tag-gap">
              {{ remindTypeLabel((item as RealItem).remind_type) }}
            </n-tag>
          </div>
          <p class="reminder-text">{{ item.content }}</p>
          <span class="reminder-date">{{ item.date }}</span>
        </div>
        <div v-if="!item._skeleton" class="reminder-actions">
          <n-button v-if="!(item as RealItem).is_read" size="small" quaternary @click="handleMarkRead((item as RealItem).id)">
            <template #icon>
              <Check :size="14" :stroke-width="1.5" />
            </template>
          </n-button>
          <n-button size="small" quaternary type="error" @click="confirmDelete((item as RealItem).id)">
            <template #icon>
              <Trash2 :size="14" :stroke-width="1.5" />
            </template>
          </n-button>
        </div>
      </div>
    </div>
    <n-empty v-if="initialized && items.length === 0" :description="t('common.noData')" class="empty-state-gap" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NTag, NEmpty, useMessage, useDialog } from 'naive-ui'
import { CheckCheck, Check, Trash2, CreditCard, Server, AlertTriangle } from '@lucide/vue'
import { listReminders, markReminderRead, markAllRemindersRead, deleteReminder } from '@/api/reminder'
import type { Reminder } from '@/types'
import { useDeferredLoading } from '@/composables/useDeferredLoading'

interface SkeletonItem { _key: string; _skeleton: true; title: string; content: string; date: string }
interface RealItem extends Reminder { _key: string; _skeleton: false }
type DisplayItem = SkeletonItem | RealItem

const SKELETON_COUNT = 4

function makeSkeletons(): SkeletonItem[] {
  return Array.from({ length: SKELETON_COUNT }, (_, i) => ({
    _key: `skel-${i}`,
    _skeleton: true as const,
    title: '',
    content: '',
    date: '',
  }))
}

const { t } = useI18n()
const message = useMessage()
const dialog = useDialog()
const items = ref<Reminder[]>([])
const { loading, withLoading } = useDeferredLoading()
const initialized = ref(false)
const hasUnread = computed(() => items.value.some(i => !i.is_read))

const displayItems = computed<DisplayItem[]>(() => {
  if (!initialized.value) return makeSkeletons()
  if (loading.value) return makeSkeletons()
  return items.value.map(r => ({ ...r, _key: `rem-${r.id}`, _skeleton: false as const }))
})

const remindTypeMap: Record<string, string> = {
  subscription_renewal: 'reminder.typeSubscriptionRenewal',
  asset_expiration: 'reminder.typeAssetExpiration',
  service_warning: 'reminder.typeServiceWarning',
}
function remindTypeLabel(v: string) { return t(remindTypeMap[v] || v) }

function remindIcon(type: string) {
  if (type === 'subscription_renewal') return CreditCard
  if (type === 'asset_expiration') return Server
  return AlertTriangle
}

onMounted(() => fetchData())

async function fetchData() {
  await withLoading(async () => {
    const res = await listReminders({ page: 1, page_size: 50 })
    items.value = (res.data as any).items
    initialized.value = true
  })
}

async function handleMarkRead(id: number) {
  try {
    await markReminderRead(id)
    await fetchData()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}

async function handleMarkAllRead() {
  try {
    await markAllRemindersRead()
    await fetchData()
    message.success(t('common.success'))
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
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
    await deleteReminder(id)
    message.success(t('common.success'))
    await fetchData()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>

<style scoped>
.reminder-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.reminder-card {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-md);
  padding: var(--spacing-md) var(--spacing-lg);
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  transition: box-shadow var(--transition-smooth), border-color var(--transition-smooth);
  position: relative;
  min-height: 82px;
}

.reminder-card:hover:not(.reminder-card--skeleton) {
  box-shadow: var(--shadow-sm);
  border-color: var(--color-border-strong);
}

.reminder-card--skeleton {
  pointer-events: none;
  opacity: 0.6;
}

.reminder-card--skeleton .reminder-icon {
  background: var(--color-bg-muted);
}

.reminder-card--skeleton .reminder-title {
  color: transparent;
  background: var(--color-bg-muted);
  border-radius: 4px;
  width: 140px;
  height: 14px;
  display: inline-block;
}

.reminder-card--skeleton .reminder-text {
  color: transparent;
  background: var(--color-bg-muted);
  border-radius: 4px;
  width: 200px;
  height: 12px;
}

.reminder-card--skeleton .reminder-date {
  color: transparent;
  background: var(--color-bg-muted);
  border-radius: 4px;
  width: 80px;
  height: 11px;
  display: inline-block;
}

.reminder-card--unread {
  border-left: 3px solid var(--color-accent);
  background: var(--gradient-card-accent);
}

.reminder-dot {
  position: absolute;
  top: var(--spacing-md);
  right: var(--spacing-lg);
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--gradient-accent);
  box-shadow: 0 0 8px rgba(34, 197, 94, 0.3);
  animation: dotBlink 2s ease-in-out infinite;
}

.reminder-icon {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md);
  background: var(--color-info-light);
  color: var(--color-info);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.reminder-content {
  flex: 1;
  min-width: 0;
}

.reminder-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: 4px;
  flex-wrap: wrap;
}

.reminder-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.reminder-title--unread {
  font-weight: 600;
}

.reminder-text {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 0 0 4px 0;
  line-height: 1.5;
}

.reminder-date {
  font-size: 12px;
  color: var(--color-text-muted);
  font-family: var(--font-mono);
  letter-spacing: -0.01em;
}

.reminder-actions {
  display: flex;
  gap: var(--spacing-xs);
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .reminder-card {
    flex-wrap: wrap;
    padding: var(--spacing-md);
  }
  .reminder-actions {
    width: 100%;
    justify-content: flex-end;
    padding-top: var(--spacing-sm);
    border-top: 1px solid var(--color-border);
  }
}
</style>
