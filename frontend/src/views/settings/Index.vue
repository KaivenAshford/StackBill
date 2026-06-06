<template>
  <div class="settings-page">
    <n-card :title="t('settings.profile')" class="card-gap">
      <n-form :model="profileForm" @submit.prevent="handleUpdateProfile">
        <n-form-item :label="t('settings.nickname')" path="nickname">
          <n-input v-model:value="profileForm.nickname" :placeholder="t('settings.nickname')" />
        </n-form-item>
        <n-button type="primary" :loading="profileLoading" attr-type="submit">
          {{ t('common.save') }}
        </n-button>
      </n-form>
    </n-card>

    <n-card :title="t('settings.appearance')" class="card-gap">
      <div class="setting-row">
        <div class="setting-info">
          <span class="setting-label">{{ t('settings.theme') }}</span>
          <span class="setting-desc">{{ isDark ? t('settings.dark') : t('settings.light') }} mode</span>
        </div>
        <n-button-group>
          <n-button :type="!isDark ? 'primary' : 'default'" @click="isDark = false">
            <template #icon><Sun :size="16" :stroke-width="1.5" /></template>
            {{ t('settings.light') }}
          </n-button>
          <n-button :type="isDark ? 'primary' : 'default'" @click="isDark = true">
            <template #icon><Moon :size="16" :stroke-width="1.5" /></template>
            {{ t('settings.dark') }}
          </n-button>
        </n-button-group>
      </div>
      <div class="setting-row">
        <div class="setting-info">
          <span class="setting-label">{{ t('settings.language') }}</span>
        </div>
        <n-button-group>
          <n-button :type="locale === 'zh-CN' ? 'primary' : 'default'" @click="switchLocale('zh-CN')">中文</n-button>
          <n-button :type="locale === 'en-US' ? 'primary' : 'default'" @click="switchLocale('en-US')">English</n-button>
        </n-button-group>
      </div>
    </n-card>

    <n-card :title="t('settings.notifications')" class="card-gap">
      <div class="setting-row">
        <div class="setting-info">
          <span class="setting-label">{{ t('settings.emailEnabled') }}</span>
        </div>
        <n-switch v-model:value="notifForm.email_enabled" />
      </div>
      <div class="setting-row">
        <div class="setting-info">
          <span class="setting-label">{{ t('settings.remindDaysBefore') }}</span>
        </div>
        <n-input-number v-model:value="notifForm.remind_days_before" :min="1" :max="30" style="width:120px" />
      </div>
      <n-button type="primary" :loading="notifLoading" @click="handleSaveNotification" style="margin-top: var(--spacing-sm)">
        {{ t('common.save') }}
      </n-button>
    </n-card>

    <n-card :title="t('settings.webhooks')" class="card-gap">
      <div v-for="wh in webhooks" :key="wh.id" class="webhook-row">
        <div class="webhook-info">
          <span class="webhook-url">{{ wh.url }}</span>
          <n-tag size="small" :type="wh.active ? 'success' : 'default'" round>{{ wh.events }}</n-tag>
        </div>
        <n-space size="small">
          <n-button size="small" quaternary @click="toggleWebhook(wh)">
            {{ wh.active ? 'Disable' : 'Enable' }}
          </n-button>
          <n-button size="small" quaternary type="error" @click="handleDeleteWebhook(wh.id)">
            {{ t('common.delete') }}
          </n-button>
        </n-space>
      </div>
      <div v-if="!webhooks.length" style="color: var(--color-text-muted); font-size: 13px; margin-bottom: var(--spacing-sm)">
        {{ t('common.noData') }}
      </div>
      <n-button @click="showAddWebhook = true">{{ t('settings.addWebhook') }}</n-button>

      <n-modal v-model:show="showAddWebhook" :title="t('settings.addWebhook')" preset="card" style="max-width:480px">
        <n-form :model="webhookForm">
          <n-form-item :label="t('settings.webhookURL')">
            <n-input v-model:value="webhookForm.url" placeholder="https://example.com/webhook" />
          </n-form-item>
          <n-form-item :label="t('settings.webhookSecret')">
            <n-input v-model:value="webhookForm.secret" :placeholder="t('settings.webhookSecret')" />
          </n-form-item>
          <n-form-item :label="t('settings.webhookEvents')">
            <n-input v-model:value="webhookForm.events" placeholder="subscription_renewal,asset_expiration" />
          </n-form-item>
        </n-form>
        <template #action>
          <n-button @click="showAddWebhook = false">{{ t('common.cancel') }}</n-button>
          <n-button type="primary" :loading="webhookSaving" @click="handleAddWebhook">{{ t('common.save') }}</n-button>
        </template>
      </n-modal>
    </n-card>

    <n-card :title="t('settings.changePassword')">
      <n-form :model="passwordForm" @submit.prevent="handleChangePassword">
        <n-form-item :label="t('settings.oldPassword')" path="old_password">
          <n-input v-model:value="passwordForm.old_password" type="password" :placeholder="t('settings.oldPassword')" />
        </n-form-item>
        <n-form-item :label="t('settings.newPassword')" path="new_password">
          <n-input v-model:value="passwordForm.new_password" type="password" :placeholder="t('settings.newPassword')" />
        </n-form-item>
        <n-form-item :label="t('settings.confirmPassword')" path="confirm_password">
          <n-input v-model:value="passwordForm.confirm_password" type="password" :placeholder="t('settings.confirmPassword')" />
        </n-form-item>
        <n-button type="primary" :loading="passwordLoading" attr-type="submit">
          {{ t('common.save') }}
        </n-button>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NCard, NForm, NFormItem, NInput, NInputNumber, NButton, NButtonGroup, NSwitch, NSpace, NTag, NModal, useMessage } from 'naive-ui'
import { Sun, Moon } from '@lucide/vue'
import { useUserStore } from '@/stores/user'
import { updateProfile, updatePassword } from '@/api/auth'
import { useTheme } from '@/composables/useTheme'
import request from '@/utils/request'

interface WebhookItem { id: number; url: string; events: string; active: boolean }

const { t, locale } = useI18n()
const message = useMessage()
const store = useUserStore()
const { isDark } = useTheme()

const profileLoading = ref(false)
const passwordLoading = ref(false)
const notifLoading = ref(false)
const webhookSaving = ref(false)
const showAddWebhook = ref(false)

const profileForm = reactive({ nickname: '', avatar: '' })
const passwordForm = reactive({ old_password: '', new_password: '', confirm_password: '' })
const notifForm = reactive({ email_enabled: false, remind_days_before: 3 })
const webhookForm = reactive({ url: '', secret: '', events: '' })
const webhooks = ref<WebhookItem[]>([])

onMounted(async () => {
  if (store.user) {
    profileForm.nickname = store.user.nickname || ''
    profileForm.avatar = store.user.avatar || ''
  }
  await loadNotificationSettings()
  await loadWebhooks()
})

async function loadNotificationSettings() {
  try {
    const res = await request.get('/notification-settings')
    notifForm.email_enabled = (res.data as any).email_enabled ?? false
    notifForm.remind_days_before = (res.data as any).remind_days_before ?? 3
  } catch { /* use defaults */ }
}

async function loadWebhooks() {
  try {
    const res = await request.get('/webhooks')
    webhooks.value = (res.data as any) || []
  } catch { webhooks.value = [] }
}

async function handleSaveNotification() {
  notifLoading.value = true
  try {
    await request.put('/notification-settings', notifForm)
    message.success(t('common.success'))
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    notifLoading.value = false
  }
}

async function handleAddWebhook() {
  webhookSaving.value = true
  try {
    await request.post('/webhooks', webhookForm)
    message.success(t('common.success'))
    showAddWebhook.value = false
    webhookForm.url = ''
    webhookForm.secret = ''
    webhookForm.events = ''
    await loadWebhooks()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    webhookSaving.value = false
  }
}

async function toggleWebhook(wh: WebhookItem) {
  try {
    await request.put(`/webhooks/${wh.id}`, { active: !wh.active })
    await loadWebhooks()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}

async function handleDeleteWebhook(id: number) {
  try {
    await request.delete(`/webhooks/${id}`)
    message.success(t('common.success'))
    await loadWebhooks()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}

function switchLocale(lang: string) {
  locale.value = lang
  localStorage.setItem('locale', lang)
  document.documentElement.lang = lang === 'en-US' ? 'en' : 'zh-CN'
}

async function handleUpdateProfile() {
  profileLoading.value = true
  try {
    const res = await updateProfile({ nickname: profileForm.nickname, avatar: profileForm.avatar })
    store.setUser(res.data)
    message.success(t('common.success'))
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    profileLoading.value = false
  }
}

async function handleChangePassword() {
  if (!passwordForm.old_password) {
    message.error(t('settings.oldPassword'))
    return
  }
  if (!passwordForm.new_password || passwordForm.new_password.length < 6) {
    message.error(t('auth.passwordMin'))
    return
  }
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    message.error(t('settings.passwordMismatch'))
    return
  }

  passwordLoading.value = true
  try {
    await updatePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password,
    })
    message.success(t('common.success'))
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    passwordLoading.value = false
  }
}
</script>

<style scoped>
.settings-page {
  max-width: 600px;
}

.settings-page :deep(.n-card) {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  transition: border-color var(--transition-smooth);
}

.settings-page :deep(.n-card:hover) {
  border-color: var(--color-border-strong);
}

.settings-page :deep(.n-card .n-card-header) {
  font-family: var(--font-heading);
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.01em;
}

.card-gap {
  margin-bottom: var(--spacing-md);
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-sm) 0;
}

.setting-row + .setting-row {
  border-top: 1px solid var(--color-border);
  padding-top: var(--spacing-md);
  margin-top: var(--spacing-md);
}

.setting-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.setting-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.setting-desc {
  font-size: 12px;
  color: var(--color-text-muted);
  font-family: var(--font-mono);
  letter-spacing: -0.01em;
}

.webhook-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--color-border);
}

.webhook-row:last-of-type {
  border-bottom: none;
  margin-bottom: var(--spacing-sm);
}

.webhook-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  min-width: 0;
}

.webhook-url {
  font-size: 13px;
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 280px;
}

@media (max-width: 768px) {
  .setting-row {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }
}
</style>
