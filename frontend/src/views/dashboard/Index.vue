<template>
  <div>
    <h2>{{ t('nav.dashboard') }}</h2>
    <n-grid :cols="4" :x-gap="16" :y-gap="16" responsive="screen" item-responsive>
      <n-gi span="4 m:2 l:1">
        <n-card>{{ t('dashboard.monthlyExpense') }}</n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>{{ t('dashboard.yearlyExpense') }}</n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>{{ t('dashboard.subscriptionCount') }}</n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>{{ t('dashboard.assetCount') }}</n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { NGrid, NGi, NCard } from 'naive-ui'
const { t } = useI18n()
</script>
