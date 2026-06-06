<template>
  <div>
    <div class="form-card">
      <h3 class="form-title">{{ isEdit ? t('common.edit') : t('common.create') }} {{ t('nav.subscriptions').toLowerCase() }}</h3>
      <n-form ref="formRef" :model="form" :rules="rules" :label-placement="isMobile ? 'top' : 'left'" label-width="100">
        <n-form-item :label="t('subscription.name')" path="name">
          <n-input v-model:value="form.name" :placeholder="t('subscription.name')" />
        </n-form-item>
        <n-form-item :label="t('subscription.amount')" path="amount">
          <div class="currency-row">
            <n-select v-model:value="form.currency" :options="currencyOptions" class="currency-select" />
            <n-input-number v-model:value="form.amount" :min="0" :precision="2" class="currency-amount" :placeholder="t('subscription.amount')" />
          </div>
        </n-form-item>
        <n-form-item :label="t('subscription.billingCycle')" path="billing_cycle">
          <n-select v-model:value="form.billing_cycle" :options="cycleOptions" />
        </n-form-item>
        <n-form-item :label="t('subscription.status')" path="status">
          <n-select v-model:value="form.status" :options="statusOptions" />
        </n-form-item>
        <n-form-item :label="t('subscription.startDate')">
          <n-date-picker v-model:formatted-value="form.start_date" type="date" value-format="yyyy-MM-dd" clearable style="width: 100%;" />
        </n-form-item>
        <n-form-item :label="t('subscription.url')">
          <n-input v-model:value="form.website_url" :placeholder="t('subscription.url')" />
        </n-form-item>
        <n-form-item :label="t('subscription.remark')">
          <n-input v-model:value="form.remark" type="textarea" :rows="3" :placeholder="t('subscription.remark')" />
        </n-form-item>
      </n-form>
      <div class="form-actions">
        <n-button @click="$router.back()">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave">{{ t('common.save') }}</n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NForm, NFormItem, NInput, NInputNumber, NSelect, NDatePicker, NSpace, NButton, useMessage } from 'naive-ui'
import { getSubscription, createSubscription, updateSubscription } from '@/api/subscription'
import { currencyOptions } from '@/utils/currency'
import { useSubscriptionStore } from '@/stores/subscription'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const message = useMessage()

const formRef = ref<InstanceType<typeof NForm> | null>(null)
const id = Number(route.params.id)
const isEdit = computed(() => !isNaN(id) && id > 0 && route.name === 'SubscriptionEdit')
const saving = ref(false)
const isMobile = ref(window.innerWidth <= 768)

function onResize() { isMobile.value = window.innerWidth <= 768 }
onMounted(() => { window.addEventListener('resize', onResize); onResize() })
onUnmounted(() => window.removeEventListener('resize', onResize))

const form = reactive({
  name: '',
  amount: 0,
  currency: 'CNY',
  billing_cycle: 'monthly',
  billing_interval: 1,
  status: 'active',
  start_date: null as string | null,
  website_url: '',
  remark: '',
})

const rules = {
  name: { required: true, message: () => t('subscription.name'), trigger: 'blur' },
  amount: { type: 'number' as const, required: true, message: () => t('subscription.amount'), trigger: 'blur' },
}

const cycleOptions = [
  { label: () => t('subscription.weekly'), value: 'weekly' },
  { label: () => t('subscription.monthly'), value: 'monthly' },
  { label: () => t('subscription.quarterly'), value: 'quarterly' },
  { label: () => t('subscription.yearly'), value: 'yearly' },
  { label: () => t('subscription.oneTime'), value: 'one_time' },
]

const statusOptions = [
  { label: () => t('subscription.active'), value: 'active' },
  { label: () => t('subscription.paused'), value: 'paused' },
  { label: () => t('subscription.cancelled'), value: 'cancelled' },
  { label: () => t('subscription.expired'), value: 'expired' },
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
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateSubscription(id, { ...form })
    } else {
      await createSubscription({ ...form })
    }
    message.success(t('common.success'))
    useSubscriptionStore().invalidate()
    router.push('/subscriptions')
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.form-card {
  max-width: 640px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  animation: slideUp 0.3s ease-out;
}

.form-title {
  font-family: var(--font-heading);
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--spacing-lg);
  letter-spacing: -0.01em;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-lg);
  padding-top: var(--spacing-lg);
  border-top: 1px solid var(--color-border);
}

.currency-row { display: flex; gap: var(--spacing-sm); width: 100%; }
.currency-select { width: 140px; }
.currency-amount { flex: 1; }

@media (max-width: 768px) {
  .form-card {
    max-width: 100%;
    padding: var(--spacing-md);
  }
  :deep(.n-form-item-label) {
    padding-bottom: var(--spacing-xs) !important;
  }
}
</style>
