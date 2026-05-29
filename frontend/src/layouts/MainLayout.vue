<template>
  <n-layout has-sider class="main-layout">
    <!-- PC sidebar -->
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      :collapsed="collapsed"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
      class="sidebar"
    >
      <div class="logo">
        <h2 v-if="!collapsed">StackBill</h2>
        <h2 v-else>SB</h2>
      </div>
      <n-menu
        :options="menuOptions"
        :value="currentRoute"
        @update:value="handleMenuClick"
      />
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered class="header">
        <div class="header-right">
          <n-button quaternary @click="handleLogout">{{ t('auth.logout') }}</n-button>
        </div>
      </n-layout-header>
      <n-layout-content class="content">
        <router-view />
      </n-layout-content>

      <!-- Mobile bottom nav -->
      <div class="mobile-nav">
        <div
          v-for="item in mobileNavItems"
          :key="item.key"
          class="mobile-nav-item"
          :class="{ active: currentRoute === item.key }"
          @click="handleMenuClick(item.key)"
        >
          <span>{{ item.label }}</span>
        </div>
      </div>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NButton } from 'naive-ui'
import { useUserStore } from '@/stores/user'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const store = useUserStore()
const collapsed = ref(false)

const currentRoute = computed(() => {
  const path = route.path
  if (path === '/') return 'dashboard'
  const segment = '/' + path.split('/')[1]
  const map: Record<string, string> = {
    '/subscriptions': 'subscriptions',
    '/assets': 'assets',
    '/categories': 'categories',
    '/reminders': 'reminders',
    '/settings': 'settings',
  }
  return map[segment] || 'dashboard'
})

const menuOptions = computed(() => [
  { label: t('nav.dashboard'), key: 'dashboard' },
  { label: t('nav.subscriptions'), key: 'subscriptions' },
  { label: t('nav.assets'), key: 'assets' },
  { label: t('nav.categories'), key: 'categories' },
  { label: t('nav.reminders'), key: 'reminders' },
  { label: t('nav.settings'), key: 'settings' },
])

const mobileNavItems = computed(() => [
  { key: 'dashboard', label: t('nav.dashboard') },
  { key: 'subscriptions', label: t('nav.subscriptions') },
  { key: 'assets', label: t('nav.assets') },
  { key: 'reminders', label: t('nav.reminders') },
  { key: 'settings', label: t('nav.settings') },
])

const routeMap: Record<string, string> = {
  dashboard: '/',
  subscriptions: '/subscriptions',
  assets: '/assets',
  categories: '/categories',
  reminders: '/reminders',
  settings: '/settings',
}

function handleMenuClick(key: string) {
  router.push(routeMap[key] || '/')
}

function handleLogout() {
  store.logout()
  router.push('/login')
}
</script>

<style scoped>
.main-layout { height: 100vh; }
.logo { padding: 16px; text-align: center; border-bottom: 1px solid #efeff5; }
.logo h2 { margin: 0; font-size: 18px; }
.header { height: 56px; display: flex; align-items: center; justify-content: flex-end; padding: 0 24px; }
.content { padding: 24px; }
.sidebar { display: block; }

.mobile-nav {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: #fff;
  border-top: 1px solid #efeff5;
  justify-content: space-around;
  align-items: center;
}
.mobile-nav-item {
  flex: 1;
  text-align: center;
  font-size: 12px;
  color: #666;
  cursor: pointer;
}
.mobile-nav-item.active { color: #18a058; }

@media (max-width: 768px) {
  .sidebar { display: none !important; }
  .mobile-nav { display: flex; }
  .content { padding: 16px; padding-bottom: 72px; }
}
</style>
