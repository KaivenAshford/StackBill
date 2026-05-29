<template>
  <div>
    <div style="display:flex;justify-content:space-between;margin-bottom:16px;">
      <h2 style="margin:0">{{ t('nav.reminders') }}</h2>
      <n-button @click="handleMarkAllRead" :disabled="!hasUnread">Mark All Read</n-button>
    </div>
    <n-list bordered>
      <n-list-item v-for="item in items" :key="item.id">
        <n-thing :title="item.title">
          <template #description>
            <n-tag size="small" :type="item.remind_type === 'service_warning' ? 'warning' : 'info'" style="margin-right:8px;">{{ item.remind_type }}</n-tag>
            <span style="color:#666;">{{ item.content }}</span>
          </template>
          <template #footer>
            <span style="color:#999;">{{ item.remind_date }}</span>
          </template>
          <template #action>
            <n-button v-if="!item.is_read" size="small" @click="handleMarkRead(item.id)">Mark Read</n-button>
          </template>
        </n-thing>
      </n-list-item>
    </n-list>
    <n-empty v-if="items.length === 0" :description="t('common.noData')" style="margin-top:40px;" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NList, NListItem, NThing, NTag, NButton, NEmpty, useMessage } from 'naive-ui'
import { listReminders, markReminderRead, markAllRemindersRead } from '@/api/reminder'
import type { Reminder } from '@/types'

const { t } = useI18n()
const message = useMessage()
const items = ref<Reminder[]>([])
const hasUnread = computed(() => items.value.some(i => !i.is_read))

onMounted(() => fetchData())

async function fetchData() {
  const res = await listReminders({ page: 1, page_size: 50 })
  items.value = (res.data as any).items
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
</script>
