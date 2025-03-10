import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/auth/Login.vue'
import Register from '@/views/auth/Register.vue'
import ForgotPassword from '@/views/auth/ForgotPassword.vue'

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
      component: () => import('@/views/auth/AuthLayout.vue'),
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
      component: () => import('@/views/home/index.vue'),
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
  } else if (token && to.path.startsWith('/auth')) {
    // 已登录但访问登录页，重定向到主页
    next({ path: '/chat' })
  } else {
    next()
  }
})

export default router
