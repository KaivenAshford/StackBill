<template>
  <div>
    <n-page-header @back="$router.back()" :title="asset?.name || ''">
      <template #extra>
        <n-space>
          <n-button @click="$router.push(`/assets/${id}/edit`)">{{ t('common.edit') }}</n-button>
          <n-button type="error" @click="confirmDelete">{{ t('common.delete') }}</n-button>
        </n-space>
      </template>
    </n-page-header>
    <n-descriptions bordered :column="2" style="margin-top:16px;" v-if="asset">
      <n-descriptions-item :label="t('asset.type')">{{ typeLabel(asset.asset_type) }}</n-descriptions-item>
      <n-descriptions-item :label="t('asset.provider')">{{ asset.provider || '-' }}</n-descriptions-item>
      <n-descriptions-item :label="t('asset.costAmount')">{{ formatAmount(asset.cost_amount, asset.cost_currency) }}</n-descriptions-item>
      <n-descriptions-item :label="t('asset.status')">
        <n-tag :type="statusType(asset.status)" size="small">{{ statusLabel(asset.status) }}</n-tag>
      </n-descriptions-item>
      <n-descriptions-item :label="t('asset.expireDate')">{{ asset.expire_date || '-' }}</n-descriptions-item>
      <n-descriptions-item :label="t('asset.url')" :span="2">{{ asset.url || '-' }}</n-descriptions-item>
      <n-descriptions-item :label="t('asset.remark')" :span="2">{{ asset.remark || '-' }}</n-descriptions-item>
    </n-descriptions>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NPageHeader, NDescriptions, NDescriptionsItem, NButton, NTag, NSpace, useMessage, useDialog } from 'naive-ui'
import { getAsset, deleteAsset } from '@/api/asset'
import { formatAmount } from '@/utils/currency'
import type { Asset } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const id = Number(route.params.id)
const asset = ref<Asset | null>(null)

const typeMap: Record<string, string> = { domain: 'asset.domain', server: 'asset.server', docker_service: 'asset.dockerService', ssl_certificate: 'asset.sslCertificate', api_key: 'asset.apiKey', repository: 'asset.repository', other: 'asset.other' }
const statusMap: Record<string, string> = { active: 'asset.active', inactive: 'asset.inactive', expired: 'asset.expired', warning: 'asset.warning' }
const statusTypeMap: Record<string, string> = { active: 'success', inactive: 'default', expired: 'error', warning: 'warning' }

function typeLabel(v: string) { return t(typeMap[v] || v) }
function statusLabel(v: string) { return t(statusMap[v] || v) }
function statusType(v: string) { return statusTypeMap[v] || 'default' }

onMounted(async () => {
  const res = await getAsset(id)
  asset.value = res.data
})

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
    await deleteAsset(id)
    message.success(t('common.success'))
    router.push('/assets')
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>
