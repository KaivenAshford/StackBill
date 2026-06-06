<template>
  <div>
    <div class="form-card">
      <h3 class="form-title">{{ isEdit ? t('common.edit') : t('common.create') }} {{ t('nav.assets').toLowerCase() }}</h3>
      <n-form ref="formRef" :model="form" :rules="rules" :label-placement="isMobile ? 'top' : 'left'" label-width="100">
        <n-form-item :label="t('asset.name')" path="name">
          <n-input v-model:value="form.name" :placeholder="t('asset.name')" />
        </n-form-item>
        <n-form-item :label="t('asset.type')" path="asset_type">
          <n-select v-model:value="form.asset_type" :options="typeOptions" />
        </n-form-item>
        <n-form-item :label="t('asset.provider')">
          <n-input v-model:value="form.provider" :placeholder="t('asset.provider')" />
        </n-form-item>
        <n-form-item :label="t('asset.identifier')">
          <n-input
            v-model:value="form.identifier"
            :placeholder="t('asset.identifier')"
            :type="form.asset_type === 'api_key' && !showIdentifier ? 'password' : 'text'"
          >
            <template v-if="form.asset_type === 'api_key'" #suffix>
              <n-button quaternary size="tiny" @click="showIdentifier = !showIdentifier">
                {{ showIdentifier ? '🙈' : '👁' }}
              </n-button>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item :label="t('asset.costAmount')">
          <div class="currency-row">
            <n-select v-model:value="form.cost_currency" :options="currencyOptions" class="currency-select" />
            <n-input-number v-model:value="form.cost_amount" :min="0" :precision="2" class="currency-amount" :placeholder="t('asset.costAmount')" />
          </div>
        </n-form-item>
        <n-form-item :label="t('asset.status')">
          <n-select v-model:value="form.status" :options="statusOptions" />
        </n-form-item>
        <n-form-item :label="t('asset.expireDate')">
          <n-date-picker v-model:formatted-value="form.expire_date" type="date" value-format="yyyy-MM-dd" clearable style="width: 100%;" />
        </n-form-item>
        <n-form-item :label="t('asset.url')">
          <n-input v-model:value="form.url" :placeholder="t('asset.url')" />
        </n-form-item>
        <n-form-item :label="t('asset.linkedSubscription')">
          <n-select v-model:value="form.subscription_id" :options="subscriptionOptions" clearable :placeholder="t('asset.noLinkedSubscription')" />
        </n-form-item>
        <n-form-item :label="t('asset.description')">
          <n-input v-model:value="form.description" type="textarea" :rows="3" :placeholder="t('asset.description')" />
        </n-form-item>
        <n-form-item :label="t('asset.remark')">
          <n-input v-model:value="form.remark" type="textarea" :rows="3" :placeholder="t('asset.remark')" />
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
import { NForm, NFormItem, NInput, NInputNumber, NSelect, NDatePicker, NButton, useMessage } from 'naive-ui'
import { getAsset, createAsset, updateAsset } from '@/api/asset'
import type { Asset } from '@/types'
import { listSubscriptions } from '@/api/subscription'
import { currencyOptions } from '@/utils/currency'
import { useAssetStore } from '@/stores/asset'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const message = useMessage()

const formRef = ref<InstanceType<typeof NForm> | null>(null)
const id = Number(route.params.id)
const isEdit = computed(() => !isNaN(id) && id > 0 && route.name === 'AssetEdit')
const saving = ref(false)
const showIdentifier = ref(false)
const isMobile = ref(window.innerWidth <= 768)

function onResize() { isMobile.value = window.innerWidth <= 768 }
onMounted(() => { window.addEventListener('resize', onResize); onResize() })
onUnmounted(() => window.removeEventListener('resize', onResize))

const form = reactive({
  name: '',
  asset_type: 'domain',
  provider: '',
  identifier: '',
  cost_amount: 0,
  cost_currency: 'CNY',
  status: 'active',
  expire_date: null as string | null,
  url: '',
  subscription_id: null as number | null,
  description: '',
  remark: '',
})

const rules = {
  name: { required: true, message: () => t('asset.name'), trigger: 'blur' },
}

const typeOptions = [
  { label: () => t('asset.domain'), value: 'domain' },
  { label: () => t('asset.server'), value: 'server' },
  { label: () => t('asset.dockerService'), value: 'docker_service' },
  { label: () => t('asset.sslCertificate'), value: 'ssl_certificate' },
  { label: () => t('asset.apiKey'), value: 'api_key' },
  { label: () => t('asset.repository'), value: 'repository' },
  { label: () => t('asset.other'), value: 'other' },
]

const statusOptions = [
  { label: () => t('asset.active'), value: 'active' },
  { label: () => t('asset.inactive'), value: 'inactive' },
  { label: () => t('asset.expired'), value: 'expired' },
  { label: () => t('asset.warning'), value: 'warning' },
]

const subscriptionOptions = ref<{ label: string; value: number }[]>([])

onMounted(async () => {
  const subRes = await listSubscriptions({ page: 1, page_size: 200 })
  const items = (subRes.data as any)?.items || []
  subscriptionOptions.value = items.map((s: any) => ({ label: s.name, value: s.id }))

  if (isEdit.value) {
    const res = await getAsset(id)
    const a = res.data
    form.name = a.name
    form.asset_type = a.asset_type
    form.provider = a.provider
    form.identifier = a.identifier
    form.cost_amount = a.cost_amount
    form.cost_currency = a.cost_currency
    form.status = a.status
    form.expire_date = a.expire_date
    form.url = a.url
    form.subscription_id = a.subscription_id || null
    form.description = a.description
    form.remark = a.remark
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
      await updateAsset(id, { ...form } as Partial<Asset>)
    } else {
      await createAsset({ ...form } as Partial<Asset>)
    }
    message.success(t('common.success'))
    useAssetStore().invalidate()
    router.push('/assets')
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
