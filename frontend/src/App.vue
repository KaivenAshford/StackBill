<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides" :locale="naiveLocale" :date-locale="naiveDateLocale">
    <n-message-provider>
      <n-dialog-provider>
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NMessageProvider, NDialogProvider, zhCN, dateZhCN, enUS, dateEnUS } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { useTheme } from '@/composables/useTheme'

const { locale } = useI18n()
const { theme, themeOverrides } = useTheme()

const naiveLocale = computed(() => locale.value === 'en-US' ? enUS : zhCN)
const naiveDateLocale = computed(() => locale.value === 'en-US' ? dateEnUS : dateZhCN)
</script>
