import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/App.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/components/auth/Login.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/components/auth/Register.vue'),
    },
  ],
})

export default router
