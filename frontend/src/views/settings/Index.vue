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
            <template #icon>
              <Sun :size="16" :stroke-width="1.5" />
            </template>
            {{ t('settings.light') }}
          </n-button>
          <n-button :type="isDark ? 'primary' : 'default'" @click="isDark = true">
            <template #icon>
              <Moon :size="16" :stroke-width="1.5" />
            </template>
            {{ t('settings.dark') }}
          </n-button>
        </n-button-group>
      </div>
      <div class="setting-row">
        <div class="setting-info">
          <span class="setting-label">{{ t('settings.language') }}</span>
        </div>
        <n-button-group>
          <n-button :type="locale === 'zh-CN' ? 'primary' : 'default'" @click="switchLocale('zh-CN')">
            中文
          </n-button>
          <n-button :type="locale === 'en-US' ? 'primary' : 'default'" @click="switchLocale('en-US')">
            English
          </n-button>
        </n-button-group>
      </div>
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
import { NCard, NForm, NFormItem, NInput, NButton, NButtonGroup, useMessage } from 'naive-ui'
import { Sun, Moon } from '@lucide/vue'
import { useUserStore } from '@/stores/user'
import { updateProfile, updatePassword } from '@/api/auth'
import { useTheme } from '@/composables/useTheme'

const { t, locale } = useI18n()
const message = useMessage()
const store = useUserStore()
const { isDark } = useTheme()

const profileLoading = ref(false)
const passwordLoading = ref(false)

const profileForm = reactive({ nickname: '', avatar: '' })
const passwordForm = reactive({ old_password: '', new_password: '', confirm_password: '' })

onMounted(() => {
  if (store.user) {
    profileForm.nickname = store.user.nickname || ''
    profileForm.avatar = store.user.avatar || ''
  }
})

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
    await updatePassword({ old_password: passwordForm.old_password, new_password: passwordForm.new_password })
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

@media (max-width: 768px) {
  .setting-row {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }
}
</style>
