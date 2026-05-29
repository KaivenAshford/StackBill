<template>
  <div>
    <n-page-header @back="$router.back()" :title="sub?.name || ''">
      <template #extra>
        <n-button @click="$router.push(`/subscriptions/${id}/edit`)">{{ t('common.edit') }}</n-button>
      </template>
    </n-page-header>
    <n-descriptions bordered :column="2" style="margin-top:16px;" v-if="sub">
      <n-descriptions-item :label="t('subscription.amount')">{{ sub.amount }} {{ sub.currency }}</n-descriptions-item>
      <n-descriptions-item :label="t('subscription.cycle')">{{ sub.billing_cycle }}</n-descriptions-item>
      <n-descriptions-item :label="t('subscription.status')">{{ sub.status }}</n-descriptions-item>
      <n-descriptions-item :label="t('subscription.nextPayment')">{{ sub.next_payment_date || '-' }}</n-descriptions-item>
      <n-descriptions-item :label="t('subscription.category')">{{ sub.category_id }}</n-descriptions-item>
      <n-descriptions-item label="URL">{{ sub.website_url || '-' }}</n-descriptions-item>
      <n-descriptions-item label="Remark" :span="2">{{ sub.remark || '-' }}</n-descriptions-item>
    </n-descriptions>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NPageHeader, NDescriptions, NDescriptionsItem, NButton } from 'naive-ui'
import { getSubscription } from '@/api/subscription'
import type { Subscription } from '@/types'

const { t } = useI18n()
const route = useRoute()
const id = Number(route.params.id)
const sub = ref<Subscription | null>(null)

onMounted(async () => {
  const res = await getSubscription(id)
  sub.value = res.data
})
</script>
