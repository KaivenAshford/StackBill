<template>
  <div>
    <h2>{{ t('nav.reminders') }}</h2>
    <p style="margin-top:16px;color:#999;">{{ t('common.noData') }}</p>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
