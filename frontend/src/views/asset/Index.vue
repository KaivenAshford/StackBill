<template>
  <div>
    <h2>{{ t('nav.assets') }}</h2>
    <n-button type="primary" @click="$router.push('/assets/new')">{{ t('common.create') }}</n-button>
    <p style="margin-top:16px;color:#999;">{{ t('common.noData') }}</p>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { NButton } from 'naive-ui'
const { t } = useI18n()
</script>
