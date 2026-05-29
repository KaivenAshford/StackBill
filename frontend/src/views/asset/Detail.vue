<template>
  <div>
    <n-page-header @back="$router.back()" :title="asset?.name || ''">
      <template #extra>
        <n-button @click="$router.push(`/assets/${id}/edit`)">{{ t('common.edit') }}</n-button>
      </template>
    </n-page-header>
    <n-descriptions bordered :column="2" style="margin-top:16px;" v-if="asset">
      <n-descriptions-item :label="t('asset.type')">{{ asset.asset_type }}</n-descriptions-item>
      <n-descriptions-item :label="t('asset.provider')">{{ asset.provider || '-' }}</n-descriptions-item>
      <n-descriptions-item :label="t('asset.status')">{{ asset.status }}</n-descriptions-item>
      <n-descriptions-item :label="t('asset.expireDate')">{{ asset.expire_date || '-' }}</n-descriptions-item>
      <n-descriptions-item label="URL" :span="2">{{ asset.url || '-' }}</n-descriptions-item>
      <n-descriptions-item label="Remark" :span="2">{{ asset.remark || '-' }}</n-descriptions-item>
    </n-descriptions>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NPageHeader, NDescriptions, NDescriptionsItem, NButton } from 'naive-ui'
import { getAsset } from '@/api/asset'
import type { Asset } from '@/types'

const { t } = useI18n()
const route = useRoute()
const id = Number(route.params.id)
const asset = ref<Asset | null>(null)

onMounted(async () => {
  const res = await getAsset(id)
  asset.value = res.data
})
</script>
