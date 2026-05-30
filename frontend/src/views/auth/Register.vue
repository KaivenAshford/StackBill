<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1>StackBill</h1>
      <p>{{ t('auth.register') }}</p>
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="top">
        <n-form-item :label="t('auth.username')" path="username">
          <n-input v-model:value="form.username" :placeholder="t('auth.username')" />
        </n-form-item>
        <n-form-item :label="t('auth.email')" path="email">
          <n-input v-model:value="form.email" :placeholder="t('auth.email')" />
        </n-form-item>
        <n-form-item :label="t('auth.password')" path="password">
          <n-input v-model:value="form.password" type="password" :placeholder="t('auth.password')" />
        </n-form-item>
        <n-button type="primary" block :loading="loading" @click="handleRegister">
          {{ t('auth.register') }}
        </n-button>
      </n-form>
      <router-link to="/login">{{ t('auth.hasAccount') }}</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import { register } from '@/api/auth'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const store = useUserStore()

const formRef = ref<InstanceType<typeof NForm> | null>(null)
const loading = ref(false)
const form = reactive({ username: '', email: '', password: '' })

const rules = {
  username: [
    { required: true, message: () => t('auth.usernameRequired'), trigger: 'blur' },
    { min: 3, max: 50, message: () => t('auth.usernameMin'), trigger: 'blur' },
  ],
  email: [
    { required: true, message: () => t('auth.emailRequired'), trigger: 'blur' },
    { type: 'email' as const, message: () => t('auth.emailFormat'), trigger: 'blur' },
  ],
  password: [
    { required: true, message: () => t('auth.passwordRequired'), trigger: 'blur' },
    { min: 6, max: 50, message: () => t('auth.passwordMin'), trigger: 'blur' },
  ],
}

async function handleRegister() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    const res = await register(form.username, form.email, form.password)
    store.setToken(res.data.token)
    store.setUser(res.data.user)
    router.push('/')
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f5f5f5;
}
.auth-card {
  width: 380px;
  padding: 40px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  text-align: center;
}
.auth-card h1 { margin-bottom: 8px; }
.auth-card a { display: inline-block; margin-top: 16px; font-size: 14px; }
</style>
