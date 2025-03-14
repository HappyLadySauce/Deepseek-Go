<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { User, Lock, Message, Key } from '@element-plus/icons-vue'

const userStore = useUserStore()
const router = useRouter()

// 表单引用
const registerFormRef = ref<FormInstance>()

// 表单数据
const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  code: ''
})

// 表单验证规则
const validatePass = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else {
    if (registerForm.confirmPassword !== '') {
      if (registerFormRef.value) {
        registerFormRef.value.validateField('confirmPassword')
      }
    }
    callback()
  }
}

const validateConfirmPass = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const registerRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, message: '用户名长度至少为3个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' },
    { validator: validatePass, trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPass, trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度应为6位', trigger: 'blur' }
  ]
})

// 发送验证码
const codeSending = ref(false)
const countDown = ref(0)
const timer = ref<number | null>(null)

const startCountDown = () => {
  countDown.value = 60
  timer.value = window.setInterval(() => {
    countDown.value--
    if (countDown.value <= 0) {
      if (timer.value) {
        clearInterval(timer.value)
        timer.value = null
      }
    }
  }, 1000)
}

const sendVerificationCode = async () => {
  // 验证邮箱
  try {
    await registerFormRef.value?.validateField('email')
    
    codeSending.value = true
    const success = await userStore.sendVerificationCode(registerForm.email)
    
    if (success) {
      ElMessage.success('验证码已发送，请检查您的邮箱')
      startCountDown()
    }
  } catch (error) {
    // 邮箱验证失败
  } finally {
    codeSending.value = false
  }
}

// 处理注册
const handleRegister = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  
  await formEl.validate(async (valid) => {
    if (valid) {
      // 先验证验证码
      const codeValid = await userStore.verifyCode(registerForm.email, registerForm.code)
      
      if (codeValid) {
        const success = await userStore.register({
          username: registerForm.username,
          email: registerForm.email,
          password: registerForm.password
        })
        
        if (success) {
          ElMessage.success('注册成功，请登录')
          router.push('/login')
        }
      }
    } else {
      ElMessage.error('请正确填写注册表单')
      return false
    }
  })
}

// 跳转到登录页
const goToLogin = () => {
  router.push('/login')
}

// 主题背景
const backgroundStyle = ref({
  backgroundImage: 'linear-gradient(to right, #4facfe 0%, #00f2fe 100%)'
})
</script>

<template>
  <div class="register-container" :style="backgroundStyle">
    <div class="register-box">
      <div class="register-brand">
        <div class="logo">
          <el-icon :size="60" color="#ffffff"><Key /></el-icon>
        </div>
        <h1 class="brand-name">DeepSeek AI</h1>
        <p class="brand-slogan">智能决策助手，专业问答平台</p>
      </div>
      
      <div class="register-form-container">
        <h2 class="register-title">注册</h2>
        <p class="register-subtitle">创建您的账号，探索AI的无限可能</p>
        
        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          class="register-form"
          label-position="top"
        >
          <el-form-item label="用户名" prop="username">
            <el-input 
              v-model="registerForm.username"
              placeholder="请输入用户名"
              size="large"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item label="邮箱" prop="email">
            <el-input
              v-model="registerForm.email"
              placeholder="请输入邮箱"
              type="email"
              size="large"
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="请输入密码"
              show-password
              size="large"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              show-password
              size="large"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item label="验证码" prop="code">
            <div class="verification-code">
              <el-input
                v-model="registerForm.code"
                placeholder="请输入验证码"
                class="code-input"
                size="large"
              />
              <el-button
                type="primary"
                :disabled="codeSending || countDown > 0"
                @click="sendVerificationCode"
                class="send-code-button"
                size="large"
              >
                {{ countDown > 0 ? `${countDown}秒后重新发送` : '发送验证码' }}
              </el-button>
            </div>
          </el-form-item>
          
          <el-button
            type="primary"
            :loading="userStore.loading"
            @click="handleRegister(registerFormRef)"
            class="register-button"
            size="large"
            round
          >
            注册
          </el-button>
          
          <div class="login-link">
            <span>已有账号？</span>
            <el-button link type="primary" @click="goToLogin">立即登录</el-button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-size: cover;
  background-position: center;
}

.register-box {
  display: flex;
  width: 900px;
  min-height: 700px;
  background-color: #fff;
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.register-brand {
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

.register-form-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 40px;
}

.register-title {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 10px;
}

.register-subtitle {
  font-size: 16px;
  color: #909399;
  margin-bottom: 30px;
}

.register-form {
  width: 100%;
}

.verification-code {
  display: flex;
  gap: 10px;
}

.code-input {
  flex: 1;
}

.register-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  margin: 20px 0;
  background-image: linear-gradient(to right, #4facfe 0%, #00f2fe 100%);
  border: none;
}

.login-link {
  text-align: center;
  font-size: 14px;
  color: #606266;
}

/* 响应式布局 */
@media (max-width: 768px) {
  .register-box {
    flex-direction: column;
    width: 90%;
    height: auto;
    max-width: 450px;
  }
  
  .register-brand {
    padding: 30px 20px;
  }
  
  .register-form-container {
    padding: 30px 20px;
  }
  
  .logo {
    width: 80px;
    height: 80px;
  }
  
  .verification-code {
    flex-direction: column;
  }
  
  .send-code-button {
    margin-top: 10px;
  }
}
</style> 