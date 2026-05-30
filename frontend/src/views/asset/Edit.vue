<template>
  <div>
    <n-card :title="isEdit ? t('common.edit') : t('common.create')">
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="100">
        <n-form-item :label="t('asset.name')" path="name">
          <n-input v-model:value="form.name" />
        </n-form-item>
        <n-form-item :label="t('asset.type')" path="asset_type">
          <n-select v-model:value="form.asset_type" :options="typeOptions" />
        </n-form-item>
        <n-form-item :label="t('asset.provider')">
          <n-input v-model:value="form.provider" />
        </n-form-item>
        <n-form-item :label="t('asset.costAmount')">
          <div style="display:flex;gap:8px;width:100%;">
            <n-select v-model:value="form.cost_currency" :options="currencyOptions" style="width:140px;" />
            <n-input-number v-model:value="form.cost_amount" :min="0" :precision="2" style="flex:1;" />
          </div>
        </n-form-item>
        <n-form-item :label="t('asset.status')">
          <n-select v-model:value="form.status" :options="statusOptions" />
        </n-form-item>
        <n-form-item :label="t('asset.expireDate')">
          <n-date-picker v-model:formatted-value="form.expire_date" type="date" value-format="yyyy-MM-dd" clearable />
        </n-form-item>
        <n-form-item :label="t('asset.url')">
          <n-input v-model:value="form.url" />
        </n-form-item>
        <n-form-item :label="t('asset.remark')">
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
import { getAsset, createAsset, updateAsset } from '@/api/asset'
import { currencyOptions } from '@/utils/currency'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const message = useMessage()

const formRef = ref<InstanceType<typeof NForm> | null>(null)
const id = Number(route.params.id)
const isEdit = computed(() => !isNaN(id) && id > 0 && route.name === 'AssetEdit')
const saving = ref(false)

const form = reactive({
  name: '',
  asset_type: 'domain',
  provider: '',
  cost_amount: 0,
  cost_currency: 'CNY',
  status: 'active',
  expire_date: null as string | null,
  url: '',
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

onMounted(async () => {
  if (isEdit.value) {
    const res = await getAsset(id)
    const a = res.data
    form.name = a.name
    form.asset_type = a.asset_type
    form.provider = a.provider
    form.cost_amount = a.cost_amount
    form.cost_currency = a.cost_currency
    form.status = a.status
    form.expire_date = a.expire_date
    form.url = a.url
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
      await updateAsset(id, { ...form })
    } else {
      await createAsset({ ...form })
    }
    message.success(t('common.success'))
    router.push('/assets')
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    saving.value = false
  }
}
</script>
