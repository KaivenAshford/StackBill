<template>
  <n-layout has-sider class="main-layout">
    <!-- PC sidebar -->
    <n-layout-sider
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      :collapsed="collapsed"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
      class="sidebar"
      :native-scrollbar="false"
    >
      <div class="sidebar-inner">
        <div class="sidebar-gradient-overlay"></div>
        <div class="logo" @click="collapsed = !collapsed" role="button" :aria-label="t('nav.dashboard')" tabindex="0">
          <div class="logo-icon">
            <span class="logo-icon-text">SB</span>
          </div>
          <transition name="fade">
            <span v-if="!collapsed" class="logo-text">StackBill</span>
          </transition>
        </div>
        <div class="sidebar-nav">
          <n-menu
            :options="menuOptions"
            :value="currentRoute"
            @update:value="handleMenuClick"
          />
        </div>
        <div class="sidebar-footer">
          <div class="sidebar-footer-divider"></div>
          <div class="sidebar-user" v-if="!collapsed">
            <div class="user-avatar">
              <img v-if="store.user?.avatar" :src="store.user.avatar" :alt="store.user.nickname" />
              <UserRound v-else :size="16" :stroke-width="1.5" />
            </div>
            <div class="user-info">
              <span class="user-name">{{ store.user?.nickname || store.user?.username || '' }}</span>
              <span class="user-plan">Free</span>
            </div>
          </div>
        </div>
      </div>
    </n-layout-sider>

    <n-layout>
      <n-layout-header class="header" :class="{ 'header-glass': !isDark }">
        <div class="header-left">
          <n-button quaternary circle size="small" @click="collapsed = !collapsed" class="menu-toggle">
            <template #icon>
              <Menu :size="18" :stroke-width="1.5" />
            </template>
          </n-button>
          <div class="breadcrumb">
            <span class="breadcrumb-text">{{ currentPageTitle }}</span>
          </div>
        </div>
        <div class="header-right">
          <n-button quaternary circle @click="toggleTheme" class="header-btn" :aria-label="isDark ? t('settings.light') : t('settings.dark')">
            <template #icon>
              <Moon v-if="!isDark" :size="18" :stroke-width="1.5" />
              <Sun v-else :size="18" :stroke-width="1.5" />
            </template>
          </n-button>
          <n-badge :value="unreadCount" :max="99">
            <n-button quaternary circle :aria-label="t('reminder.unread')" @click="router.push('/reminders')" class="header-btn">
              <template #icon>
                <Bell :size="18" :stroke-width="1.5" />
              </template>
            </n-button>
          </n-badge>
          <n-button quaternary size="small" @click="handleLogout" class="logout-btn" :aria-label="t('auth.logout')">
            <template #icon>
              <LogOut :size="16" :stroke-width="1.5" />
            </template>
            {{ t('auth.logout') }}
          </n-button>
        </div>
      </n-layout-header>
      <n-layout-content class="content">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </n-layout-content>

      <!-- Mobile bottom nav -->
      <div class="mobile-nav">
        <div
          v-for="item in mobileNavItems"
          :key="item.key"
          class="mobile-nav-item"
          :class="{ active: currentRoute === item.key }"
          role="button"
          :aria-label="item.label"
          tabindex="0"
          @click="handleMenuClick(item.key)"
          @keydown.enter="handleMenuClick(item.key)"
        >
          <div class="mobile-nav-icon">
            <component :is="item.icon" :size="20" :stroke-width="1.5" />
          </div>
          <span class="mobile-nav-label">{{ item.label }}</span>
        </div>
      </div>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NButton, NBadge } from 'naive-ui'
import {
  LayoutDashboard,
  CreditCard,
  Server,
  FolderTree,
  Bell,
  Settings,
  LogOut,
  UserRound,
  Moon,
  Sun,
  Menu,
} from '@lucide/vue'
import { useUserStore } from '@/stores/user'
import { listReminders } from '@/api/reminder'
import { useTheme } from '@/composables/useTheme'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const store = useUserStore()
const collapsed = ref(false)
const unreadCount = ref(0)
const { isDark, toggleTheme } = useTheme()

onMounted(async () => {
  if (store.isLoggedIn() && !store.user) {
    await store.fetchUser()
  }
  if (store.isLoggedIn()) {
    try {
      const res = await listReminders({ is_read: false, page: 1, page_size: 1 })
      unreadCount.value = res.data.total
    } catch {
      // ignore — non-critical
    }
  }
})

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

const currentPageTitle = computed(() => {
  const map: Record<string, string> = {
    dashboard: t('nav.dashboard'),
    subscriptions: t('nav.subscriptions'),
    assets: t('nav.assets'),
    categories: t('nav.categories'),
    reminders: t('nav.reminders'),
    settings: t('nav.settings'),
  }
  return map[currentRoute.value] || ''
})

const menuOptions = computed(() => [
  { label: t('nav.dashboard'), key: 'dashboard', icon: () => h(LayoutDashboard, { size: 20, strokeWidth: 1.5 }) },
  { label: t('nav.subscriptions'), key: 'subscriptions', icon: () => h(CreditCard, { size: 20, strokeWidth: 1.5 }) },
  { label: t('nav.assets'), key: 'assets', icon: () => h(Server, { size: 20, strokeWidth: 1.5 }) },
  { label: t('nav.categories'), key: 'categories', icon: () => h(FolderTree, { size: 20, strokeWidth: 1.5 }) },
  { label: t('nav.reminders'), key: 'reminders', icon: () => h(Bell, { size: 20, strokeWidth: 1.5 }) },
  { label: t('nav.settings'), key: 'settings', icon: () => h(Settings, { size: 20, strokeWidth: 1.5 }) },
])

const mobileNavItems = computed(() => [
  { key: 'dashboard', label: t('nav.dashboard'), icon: LayoutDashboard },
  { key: 'subscriptions', label: t('nav.subscriptions'), icon: CreditCard },
  { key: 'assets', label: t('nav.assets'), icon: Server },
  { key: 'reminders', label: t('nav.reminders'), icon: Bell },
  { key: 'settings', label: t('nav.settings'), icon: Settings },
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

/* ---- Sidebar ---- */
.sidebar {
  background: var(--gradient-sidebar) !important;
  border-right: none !important;
  display: block;
  position: relative;
  overflow: hidden;
}

.sidebar::after {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 1px;
  background: linear-gradient(
    180deg,
    transparent 0%,
    rgba(34, 197, 94, 0.1) 50%,
    transparent 100%
  );
}

.sidebar-inner {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  position: relative;
  z-index: 1;
}

.sidebar-gradient-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 200px;
  background: radial-gradient(
    ellipse at 30% 0%,
    rgba(34, 197, 94, 0.06) 0%,
    transparent 70%
  );
  pointer-events: none;
  z-index: 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-lg) var(--spacing-md);
  cursor: pointer;
  user-select: none;
  min-height: 64px;
  position: relative;
  z-index: 1;
}

.logo-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  background: var(--gradient-accent);
  color: #0F172A;
  font-family: var(--font-heading);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: relative;
  overflow: hidden;
}

.logo-icon::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.2) 0%,
    transparent 50%
  );
}

.logo-icon-text {
  position: relative;
  z-index: 1;
  font-size: 14px;
  font-weight: 700;
}

.logo-text {
  font-family: var(--font-heading);
  font-size: 18px;
  font-weight: 700;
  color: #F8FAFC;
  letter-spacing: -0.02em;
  white-space: nowrap;
}

.sidebar-nav {
  flex: 1;
  padding: 0 var(--spacing-sm);
  overflow-y: auto;
  position: relative;
  z-index: 1;
}

.sidebar-footer {
  padding: 0 var(--spacing-md) var(--spacing-md);
  position: relative;
  z-index: 1;
}

.sidebar-footer-divider {
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(248, 250, 252, 0.08) 50%,
    transparent 100%
  );
  margin-bottom: var(--spacing-md);
}

.sidebar-user {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm);
  border-radius: var(--radius-md);
  background: rgba(248, 250, 252, 0.04);
  transition: background var(--transition-fast);
}

.sidebar-user:hover {
  background: rgba(248, 250, 252, 0.06);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(248, 250, 252, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
  color: #94A3B8;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
  color: #E2E8F0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-plan {
  font-size: 11px;
  color: #64748B;
}

/* ---- Header ---- */
.header {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-lg);
  background: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  transition: background var(--transition-base), border-color var(--transition-base);
}

[data-theme="dark"] .header {
  background: rgba(20, 28, 46, 0.8);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.menu-toggle {
  display: none;
}

.breadcrumb-text {
  font-family: var(--font-heading);
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.header-btn {
  color: var(--color-text-secondary) !important;
  transition: color var(--transition-fast), background var(--transition-fast);
}

.header-btn:hover {
  color: var(--color-accent) !important;
}

.logout-btn {
  margin-left: var(--spacing-xs);
  color: var(--color-text-secondary) !important;
  transition: color var(--transition-fast);
}

/* ---- Content ---- */
.content {
  padding: var(--spacing-lg);
  background: var(--color-bg);
  min-height: calc(100vh - 56px);
}

/* ---- Mobile bottom nav ---- */
.mobile-nav {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: var(--color-bg-card);
  border-top: 1px solid var(--color-border);
  justify-content: space-around;
  align-items: center;
  padding-bottom: env(safe-area-inset-bottom, 0);
  z-index: var(--z-mobile-nav);
}

[data-theme="dark"] .mobile-nav {
  background: rgba(20, 28, 46, 0.9);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.mobile-nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  text-align: center;
  cursor: pointer;
  padding: var(--spacing-xs) 0;
  min-height: 44px;
  justify-content: center;
  position: relative;
  transition: color var(--transition-fast);
  -webkit-tap-highlight-color: transparent;
}

.mobile-nav-item .mobile-nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 28px;
  border-radius: var(--radius-full);
  color: var(--color-text-muted);
  transition: all var(--transition-smooth);
}

.mobile-nav-item .mobile-nav-label {
  font-size: 11px;
  font-weight: 500;
  color: var(--color-text-muted);
  transition: color var(--transition-fast);
}

.mobile-nav-item.active .mobile-nav-icon {
  color: var(--color-accent);
  background: var(--color-accent-light);
}

.mobile-nav-item.active .mobile-nav-label {
  color: var(--color-accent);
  font-weight: 600;
}

.mobile-nav-item.active::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 2px;
  border-radius: 0 0 2px 2px;
  background: var(--gradient-accent);
}

/* ---- Page transition ---- */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: opacity 80ms ease;
}
.fade-slide-enter-from {
  opacity: 0;
}
.fade-slide-leave-to {
  opacity: 0;
}

/* ---- Logo text transition ---- */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 150ms ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ---- Responsive ---- */
@media (max-width: 768px) {
  .sidebar { display: none !important; }
  .mobile-nav { display: flex; }
  .content { padding: var(--spacing-md); padding-bottom: 76px; }
  .menu-toggle { display: flex; }
  .breadcrumb { display: none; }
  .logout-btn span { display: none; }
}
</style>
