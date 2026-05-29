<template>
  <div>
    <n-card :title="isEdit ? t('common.edit') : t('common.create')">
      <n-form :model="form" label-placement="left" label-width="100">
        <n-form-item :label="t('subscription.name')" path="name">
          <n-input v-model:value="form.name" />
        </n-form-item>
        <n-form-item :label="t('subscription.amount')" path="amount">
          <n-input-number v-model:value="form.amount" :min="0" :precision="2" />
        </n-form-item>
        <n-form-item :label="t('subscription.cycle')" path="billing_cycle">
          <n-select v-model:value="form.billing_cycle" :options="cycleOptions" />
        </n-form-item>
        <n-form-item :label="t('subscription.status')" path="status">
          <n-select v-model:value="form.status" :options="statusOptions" />
        </n-form-item>
        <n-form-item label="Start Date">
          <n-date-picker v-model:formatted-value="form.start_date" type="date" value-format="yyyy-MM-dd" clearable />
        </n-form-item>
        <n-form-item label="URL">
          <n-input v-model:value="form.website_url" />
        </n-form-item>
        <n-form-item label="Remark">
          <n-input v-model:value="form.remark" type="textarea" :rows="3" />
        </n-form-item>
      </n-form>
      <n-space>
        <n-button @click="$router.back()">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave">{{ t('common.save') }}</n-button>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NCard, NForm, NFormItem, NInput, NInputNumber, NSelect, NDatePicker, NSpace, NButton, useMessage } from 'naive-ui'
import { getSubscription, createSubscription, updateSubscription } from '@/api/subscription'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const message = useMessage()

const id = Number(route.params.id)
const isEdit = computed(() => !isNaN(id) && route.name === 'SubscriptionEdit')
const saving = ref(false)

const form = reactive({
  name: '',
  amount: 0,
  currency: 'USD',
  billing_cycle: 'monthly',
  billing_interval: 1,
  status: 'active',
  start_date: null as string | null,
  website_url: '',
  remark: '',
})

const cycleOptions = [
  { label: 'Weekly', value: 'weekly' },
  { label: 'Monthly', value: 'monthly' },
  { label: 'Quarterly', value: 'quarterly' },
  { label: 'Yearly', value: 'yearly' },
  { label: 'One Time', value: 'one_time' },
]

const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Paused', value: 'paused' },
  { label: 'Cancelled', value: 'cancelled' },
  { label: 'Expired', value: 'expired' },
]

onMounted(async () => {
  if (isEdit.value) {
    const res = await getSubscription(id)
    const s = res.data
    form.name = s.name
    form.amount = s.amount
    form.currency = s.currency
    form.billing_cycle = s.billing_cycle
    form.billing_interval = s.billing_interval
    form.status = s.status
    form.start_date = s.start_date
    form.website_url = s.website_url
    form.remark = s.remark
  }
})

async function handleSave() {
  saving.value = true
  try {
    if (isEdit.value) {
      await updateSubscription(id, { ...form })
    } else {
      await createSubscription({ ...form })
    }
    message.success(t('common.success'))
    router.push('/subscriptions')
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    saving.value = false
  }
}
</script>
