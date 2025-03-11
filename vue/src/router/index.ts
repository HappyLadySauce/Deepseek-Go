import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/components/auth/Login.vue'
import Register from '@/components/auth/Register.vue'
import ForgotPassword from '@/components/auth/ForgotPassword.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // 重定向到登录页
    {
      path: '/',
      redirect: '/auth/login'
    },
    // 认证相关路由
    {
      path: '/auth',
      component: () => import('@/components/auth/AuthLayout.vue'),
      children: [
        { path: 'login', component: Login },
        { path: 'register', component: Register },
        { path: 'forgot-password', component: ForgotPassword }
      ]
    },
    // 单独访问 /auth 时，会重定向到 /auth/login
    {
      path: '/auth',
      redirect: '/auth/login'
    },
    // 主应用路由
    {
      path: '/overview',
      component: () => import('@/views/overview/index.vue'),
    },
    {
      path: '/monitor',
      component: () => import('@/views/monitor/system.vue'),
      children: [
        { path: 'system', component: () => import('@/views/monitor/system.vue') },
        { path: 'network', component: () => import('@/views/monitor/network.vue') },
        { path: 'alarm', component: () => import('@/views/monitor/alarm.vue') },
      ]
    },
    {
      path: '/settings',
      component: () => import('@/views/settings/index.vue'),
    },
    // AI聊天路由
    {
      path: '/chat',
      component: () => import('@/views/chat/index.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  
  if (requiresAuth && !token) {
    // 需要认证但未登录，重定向到登录页
    next({
      path: '/auth/login',
      query: { redirect: to.fullPath } // 保存重定向地址
    })
  } else if (token && (to.path.startsWith('/auth') || to.path.startsWith('/auth/login') || to.path.startsWith('/auth/register') || to.path.startsWith('/auth/forgot-password'))) {
    // 已登录但访问登录页，重定向到主页
    next({ path: '/overview' })
  } else {
    next()
  }
})

export default router
