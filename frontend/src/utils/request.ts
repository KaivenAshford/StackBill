import axios from 'axios'
import router from '@/router'
import i18n from '@/i18n'

const ERROR_CODE_I18N_MAP: Record<number, string> = {
  40001: 'errors.invalidParams',
  40002: 'errors.invalidReminderId',
  40003: 'errors.incorrectPassword',
  40004: 'errors.usernameRequired',
  40005: 'errors.emailInvalid',
  40006: 'errors.passwordRequired',
  40007: 'errors.credentialsRequired',
  40100: 'errors.unauthorized',
  40101: 'errors.invalidCredentials',
  40301: 'errors.forbidden',
  40400: 'errors.notFound',
  40901: 'errors.duplicateUsername',
  40902: 'errors.duplicateEmail',
  40903: 'errors.duplicateCategory',
}

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

request.interceptors.response.use(
  (response) => {
    const data = response.data
    if (data.code !== undefined && data.code !== 0) {
      const i18nKey = ERROR_CODE_I18N_MAP[data.code]
      const msg = i18nKey
        ? (i18n.global as { t: (key: string) => string }).t(i18nKey)
        : (data.message || (i18n.global as { t: (key: string) => string }).t('errors.requestFailed'))
      return Promise.reject(new Error(msg))
    }
    return data
  },
  (error) => {
    const code = error.response?.data?.code
    const i18nKey = code ? ERROR_CODE_I18N_MAP[code] : undefined
    const msg = i18nKey
      ? (i18n.global as { t: (key: string) => string }).t(i18nKey)
      : (error.response?.data?.message || error.message || (i18n.global as { t: (key: string) => string }).t('errors.requestFailed'))
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      router.push('/login')
    }
    return Promise.reject(new Error(msg))
  },
)

export default request
