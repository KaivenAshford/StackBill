import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/auth/Login.vue'),
      meta: { guest: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/auth/Register.vue'),
      meta: { guest: true },
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { auth: true },
      children: [
        { path: '', name: 'Dashboard', component: () => import('@/views/dashboard/Index.vue') },
        { path: 'subscriptions', name: 'Subscriptions', component: () => import('@/views/subscription/Index.vue') },
        { path: 'subscriptions/:id', name: 'SubscriptionDetail', component: () => import('@/views/subscription/Detail.vue') },
        { path: 'subscriptions/new', name: 'SubscriptionNew', component: () => import('@/views/subscription/Edit.vue') },
        { path: 'subscriptions/:id/edit', name: 'SubscriptionEdit', component: () => import('@/views/subscription/Edit.vue') },
        { path: 'assets', name: 'Assets', component: () => import('@/views/asset/Index.vue') },
        { path: 'assets/:id', name: 'AssetDetail', component: () => import('@/views/asset/Detail.vue') },
        { path: 'assets/new', name: 'AssetNew', component: () => import('@/views/asset/Edit.vue') },
        { path: 'assets/:id/edit', name: 'AssetEdit', component: () => import('@/views/asset/Edit.vue') },
        { path: 'categories', name: 'Categories', component: () => import('@/views/category/Index.vue') },
        { path: 'reminders', name: 'Reminders', component: () => import('@/views/reminder/Index.vue') },
        { path: 'settings', name: 'Settings', component: () => import('@/views/settings/Index.vue') },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const store = useUserStore()
  if (to.meta.auth && !store.isLoggedIn()) {
    return { name: 'Login' }
  }
  if (to.meta.guest && store.isLoggedIn()) {
    return { name: 'Dashboard' }
  }
})

export default router
