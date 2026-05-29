<template>
  <div class="settings-page">
    <n-card :title="t('settings.profile')" style="margin-bottom: 24px;">
      <n-form :model="profileForm" @submit.prevent="handleUpdateProfile">
        <n-form-item :label="t('settings.nickname')" path="nickname">
          <n-input v-model:value="profileForm.nickname" :placeholder="t('settings.nickname')" />
        </n-form-item>
        <n-button type="primary" :loading="profileLoading" attr-type="submit">
          {{ t('common.save') }}
        </n-button>
      </n-form>
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
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import { updateProfile, updatePassword } from '@/api/auth'

const { t } = useI18n()
const message = useMessage()
const store = useUserStore()

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
</style>
