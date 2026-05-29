<template>
  <div>
    <n-card :title="isEdit ? t('common.edit') : t('common.create')">
      <n-form :model="form" label-placement="left" label-width="100">
        <n-form-item :label="t('asset.name')" path="name">
          <n-input v-model:value="form.name" />
        </n-form-item>
        <n-form-item :label="t('asset.type')" path="asset_type">
          <n-select v-model:value="form.asset_type" :options="typeOptions" />
        </n-form-item>
        <n-form-item :label="t('asset.provider')">
          <n-input v-model:value="form.provider" />
        </n-form-item>
        <n-form-item :label="t('asset.status')">
          <n-select v-model:value="form.status" :options="statusOptions" />
        </n-form-item>
        <n-form-item :label="t('asset.expireDate')">
          <n-date-picker v-model:formatted-value="form.expire_date" type="date" value-format="yyyy-MM-dd" clearable />
        </n-form-item>
        <n-form-item label="URL">
          <n-input v-model:value="form.url" />
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
import { NCard, NForm, NFormItem, NInput, NSelect, NDatePicker, NSpace, NButton, useMessage } from 'naive-ui'
import { getAsset, createAsset, updateAsset } from '@/api/asset'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const message = useMessage()

const id = Number(route.params.id)
const isEdit = computed(() => !isNaN(id) && route.name === 'AssetEdit')
const saving = ref(false)

const form = reactive({
  name: '',
  asset_type: 'domain',
  provider: '',
  status: 'active',
  expire_date: null as string | null,
  url: '',
  remark: '',
})

const typeOptions = [
  { label: 'Domain', value: 'domain' },
  { label: 'Server', value: 'server' },
  { label: 'Docker Service', value: 'docker_service' },
  { label: 'SSL Certificate', value: 'ssl_certificate' },
  { label: 'API Key', value: 'api_key' },
  { label: 'Repository', value: 'repository' },
  { label: 'Other', value: 'other' },
]

const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Expired', value: 'expired' },
  { label: 'Warning', value: 'warning' },
]

onMounted(async () => {
  if (isEdit.value) {
    const res = await getAsset(id)
    const a = res.data
    form.name = a.name
    form.asset_type = a.asset_type
    form.provider = a.provider
    form.status = a.status
    form.expire_date = a.expire_date
    form.url = a.url
    form.remark = a.remark
  }
})

async function handleSave() {
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
