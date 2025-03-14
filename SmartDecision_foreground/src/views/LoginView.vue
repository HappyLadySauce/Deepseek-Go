<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { User, Lock, Key, Message, Moon, Sunny } from '@element-plus/icons-vue'
import { currentTheme, toggleTheme } from '../utils/theme'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

// 表单引用
const loginFormRef = ref<FormInstance>()

// 表单数据
const loginForm = reactive({
  username: '',
  password: ''
})

// 记住密码
const rememberMe = ref(false)

// 主题背景
const backgroundStyle = ref({
  backgroundImage: 'linear-gradient(to right, #4facfe 0%, #00f2fe 100%)'
})

// 加载本地存储的用户名和密码
onMounted(() => {
  const savedUsername = localStorage.getItem('rememberedUsername')
  const savedPassword = localStorage.getItem('rememberedPassword')
  
  if (savedUsername && savedPassword) {
    loginForm.username = savedUsername
    loginForm.password = savedPassword
    rememberMe.value = true
  }
})

// 表单验证规则
const loginRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名或邮箱', trigger: 'blur' },
    { min: 3, message: '用户名长度至少为3个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
  ]
})

// 处理登录
const handleLogin = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  
  await formEl.validate(async (valid) => {
    if (valid) {
      // 如果勾选了记住密码，保存到本地存储
      if (rememberMe.value) {
        localStorage.setItem('rememberedUsername', loginForm.username)
        localStorage.setItem('rememberedPassword', loginForm.password)
      } else {
        localStorage.removeItem('rememberedUsername')
        localStorage.removeItem('rememberedPassword')
      }
      
      const success = await userStore.login(loginForm)
      if (success) {
        // 登录成功后重定向
        const redirectPath = route.query.redirect as string || '/chat'
        router.push(redirectPath)
      }
    } else {
      ElMessage.error('请正确填写登录表单')
      return false
    }
  })
}

// 跳转到注册页
const goToRegister = () => {
  router.push('/register')
}

// 跳转到忘记密码页面
const goToForgotPassword = () => {
  router.push('/verify-email')
}

// 切换主题
const isDarkMode = computed(() => currentTheme.value === 'dark')
const handleToggleTheme = () => {
  toggleTheme()
}
</script>

<template>
  <div class="login-container" :style="backgroundStyle">
    <div class="theme-toggle-button" @click="handleToggleTheme">
      <el-icon :size="24"><component :is="isDarkMode ? Sunny : Moon" /></el-icon>
    </div>
    
    <div class="login-box">
      <div class="login-brand">
        <div class="logo">
          <el-icon :size="60" color="#ffffff"><Key /></el-icon>
        </div>
        <h1 class="brand-name">DeepSeek AI</h1>
        <p class="brand-slogan">智能决策助手，让AI为您解答一切</p>
      </div>
      
      <div class="login-form-container">
        <h2 class="login-title">登录</h2>
        <p class="login-subtitle">欢迎回来，请登录您的账号</p>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="用户名或邮箱"
              size="large"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              show-password
              size="large"
              @keyup.enter="handleLogin(loginFormRef)"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <div class="form-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <el-button link type="primary" @click="goToForgotPassword">忘记密码？</el-button>
          </div>
          
          <el-button
            type="primary"
            :loading="userStore.loading"
            @click="handleLogin(loginFormRef)"
            class="login-button"
            size="large"
            round
          >
            登录
          </el-button>
          
          <div class="register-link">
            <span>还没有账号？</span>
            <el-button link type="primary" @click="goToRegister">立即注册</el-button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-size: cover;
  background-position: center;
  position: relative;
}

.theme-toggle-button {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: white;
  transition: all 0.3s;
  z-index: 10;
}

.theme-toggle-button:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.login-box {
  display: flex;
  width: 900px;
  height: 600px;
  background-color: var(--card-bg);
  border-radius: 20px;
  box-shadow: var(--card-shadow);
  overflow: hidden;
}

.login-brand {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #4facfe;
  background-image: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  padding: 40px;
  text-align: center;
}

.logo {
  margin-bottom: 20px;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.brand-name {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 10px;
}

.brand-slogan {
  font-size: 16px;
  opacity: 0.9;
  max-width: 300px;
}

.login-form-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px 40px;
}

.login-title {
  font-size: 32px;
  font-weight: bold;
  color: var(--heading-color);
  margin-bottom: 10px;
}

.login-subtitle {
  font-size: 16px;
  color: var(--text-muted);
  margin-bottom: 40px;
}

.login-form {
  width: 100%;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  color: var(--text-color);
}

.login-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  margin-bottom: 20px;
  background-image: linear-gradient(to right, #4facfe 0%, #00f2fe 100%);
  border: none;
}

.register-link {
  text-align: center;
  font-size: 14px;
  color: var(--text-light);
}

/* 响应式布局 */
@media (max-width: 768px) {
  .login-box {
    flex-direction: column;
    width: 90%;
    height: auto;
    max-width: 450px;
  }
  
  .login-brand {
    padding: 30px 20px;
  }
  
  .login-form-container {
    padding: 40px 30px;
  }
  
  .logo {
    width: 80px;
    height: 80px;
  }
}
</style> 