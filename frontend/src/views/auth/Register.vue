<template>
  <div class="auth-page">
    <!-- Left brand panel -->
    <div class="auth-brand-panel">
      <div class="brand-grid-bg"></div>
      <div class="brand-content">
        <div class="brand-logo">
          <div class="logo-icon">
            <span class="logo-icon-inner">SB</span>
          </div>
          <span class="logo-text">StackBill</span>
        </div>
        <h1 class="brand-headline">{{ t('auth.brandSlogan') }}</h1>
        <div class="brand-features">
          <div class="brand-feature stagger-1">
            <div class="feature-icon">
              <CreditCard :size="20" :stroke-width="1.5" />
            </div>
            <span>{{ t('auth.featureSubscription') }}</span>
          </div>
          <div class="brand-feature stagger-2">
            <div class="feature-icon">
              <Server :size="20" :stroke-width="1.5" />
            </div>
            <span>{{ t('auth.featureAsset') }}</span>
          </div>
          <div class="brand-feature stagger-3">
            <div class="feature-icon">
              <Bell :size="20" :stroke-width="1.5" />
            </div>
            <span>{{ t('auth.featureReminder') }}</span>
          </div>
        </div>
      </div>
      <div class="brand-decoration">
        <div class="deco-circle deco-1"></div>
        <div class="deco-circle deco-2"></div>
        <div class="deco-circle deco-3"></div>
        <div class="deco-glow"></div>
      </div>
    </div>

    <!-- Right form panel -->
    <div class="auth-form-panel">
      <div class="auth-form-container">
        <div class="auth-mobile-brand">
          <div class="logo-icon-sm">SB</div>
          <span class="logo-text-sm">StackBill</span>
        </div>
        <h2 class="auth-title">{{ t('auth.createAccount') }}</h2>
        <p class="auth-subtitle">{{ t('auth.registerSubtitle') }}</p>

        <n-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleRegister">
          <n-form-item :label="t('auth.username')" path="username">
            <n-input v-model:value="form.username" :placeholder="t('auth.username')" size="large">
              <template #prefix>
                <User :size="18" :stroke-width="1.5" style="color: var(--color-text-muted);" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item :label="t('auth.email')" path="email">
            <n-input v-model:value="form.email" :placeholder="t('auth.email')" size="large">
              <template #prefix>
                <Mail :size="18" :stroke-width="1.5" style="color: var(--color-text-muted);" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item :label="t('auth.password')" path="password">
            <n-input v-model:value="form.password" type="password" :placeholder="t('auth.password')" show-password-on="click" size="large">
              <template #prefix>
                <Lock :size="18" :stroke-width="1.5" style="color: var(--color-text-muted);" />
              </template>
            </n-input>
          </n-form-item>
          <n-button type="primary" block :loading="loading" attr-type="submit" size="large" class="auth-submit">
            {{ t('auth.register') }}
          </n-button>
        </n-form>
        <p class="auth-switch">
          {{ t('auth.hasAccount') }}
          <router-link to="/login">{{ t('auth.login') }}</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { User, Mail, Lock, CreditCard, Server, Bell } from '@lucide/vue'
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
  min-height: 100vh;
  background: var(--color-bg);
}

/* ---- Left brand panel ---- */
.auth-brand-panel {
  position: relative;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-hero);
  overflow: hidden;
  padding: var(--spacing-2xl);
}

.brand-grid-bg {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(248, 250, 252, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(248, 250, 252, 0.03) 1px, transparent 1px);
  background-size: 60px 60px;
  mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
  -webkit-mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
}

.brand-content {
  position: relative;
  z-index: 2;
  max-width: 440px;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-2xl);
}

.logo-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  background: var(--gradient-accent);
  color: #0F172A;
  font-family: var(--font-heading);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.logo-icon::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.25) 0%, transparent 50%);
}

.logo-icon-inner {
  position: relative;
  z-index: 1;
  font-size: 18px;
  font-weight: 700;
}

.logo-text {
  font-family: var(--font-heading);
  font-size: 24px;
  font-weight: 700;
  color: #F8FAFC;
  letter-spacing: -0.02em;
}

.brand-headline {
  font-family: var(--font-heading);
  font-size: 32px;
  font-weight: 700;
  color: #F8FAFC;
  line-height: 1.2;
  letter-spacing: -0.02em;
  margin-bottom: var(--spacing-2xl);
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.brand-feature {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  color: #94A3B8;
  font-size: 15px;
}

.feature-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  background: rgba(248, 250, 252, 0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
  flex-shrink: 0;
  border: 1px solid rgba(248, 250, 252, 0.04);
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.brand-feature:hover .feature-icon {
  background: rgba(34, 197, 94, 0.1);
  border-color: rgba(34, 197, 94, 0.2);
}

.brand-decoration {
  position: absolute;
  inset: 0;
  z-index: 1;
  pointer-events: none;
}

.deco-circle {
  position: absolute;
  border-radius: 50%;
  border: 1px solid rgba(248, 250, 252, 0.03);
}

.deco-1 {
  width: 600px;
  height: 600px;
  top: -200px;
  right: -200px;
  background: radial-gradient(circle, rgba(34, 197, 94, 0.08) 0%, transparent 70%);
}

.deco-2 {
  width: 400px;
  height: 400px;
  bottom: -100px;
  left: -100px;
  background: radial-gradient(circle, rgba(34, 197, 94, 0.05) 0%, transparent 70%);
}

.deco-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  left: 60%;
  background: radial-gradient(circle, rgba(248, 250, 252, 0.03) 0%, transparent 70%);
}

.deco-glow {
  position: absolute;
  top: 20%;
  right: 10%;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(34, 197, 94, 0.04) 0%, transparent 60%);
  animation: float 8s ease-in-out infinite;
}

/* ---- Right form panel ---- */
.auth-form-panel {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-2xl);
  background: var(--color-bg);
}

.auth-form-container {
  width: 100%;
  max-width: 400px;
  animation: slideUp 0.5s ease-out;
}

.auth-mobile-brand {
  display: none;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-2xl);
}

.logo-icon-sm {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  background: var(--gradient-accent);
  color: #0F172A;
  font-family: var(--font-heading);
  font-size: 14px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text-sm {
  font-family: var(--font-heading);
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.auth-title {
  font-family: var(--font-heading);
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: var(--spacing-xs);
  letter-spacing: -0.02em;
}

.auth-subtitle {
  color: var(--color-text-muted);
  font-size: 15px;
  margin-bottom: var(--spacing-2xl);
}

.auth-submit {
  margin-top: var(--spacing-sm);
  height: 44px;
  font-weight: 600;
  font-size: 15px;
  border-radius: var(--radius-md);
}

.auth-switch {
  text-align: center;
  margin-top: var(--spacing-xl);
  font-size: 14px;
  color: var(--color-text-secondary);
}

.auth-switch a {
  font-weight: 500;
}

/* ---- Responsive ---- */
@media (max-width: 960px) {
  .auth-brand-panel {
    display: none;
  }

  .auth-mobile-brand {
    display: flex;
  }

  .auth-form-panel {
    padding: var(--spacing-lg);
  }
}
</style>
