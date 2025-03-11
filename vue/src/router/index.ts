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
      path: '/index',
      component: () => import('@/views/overview/index.vue'),
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
    next({ path: '/index' })
  } else {
    next()
  }
})

export default router
