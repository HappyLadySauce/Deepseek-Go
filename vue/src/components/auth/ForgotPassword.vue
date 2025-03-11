<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Message, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import type { FormInstance } from 'element-plus'
import axios from '@/axios'
import CaptchaImage from '@/components/auth/CaptchaImage.vue'
import EmailVerification from '@/components/auth/EmailVerification.vue'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

// 步骤控制: 1.输入邮箱 2.验证邮箱 3.重置密码
const currentStep = ref(1)
const isEmailVerified = ref(false)

// 图形验证码相关
const captchaCode = ref('')
const captchaRef = ref()

// 表单数据
const formData = reactive({
  email: '',
  captcha: '',
  password: '',
  confirmPassword: '',
})

// 表单错误
const formErrors = reactive({
  email: '',
  captcha: '',
  password: '',
  confirmPassword: '',
  general: ''
})

// 校验规则
const rules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { 
      validator: (rule: any, value: string, callback: any) => {
        if (value && value.toLowerCase() !== captchaCode.value.toLowerCase()) {
          callback(new Error('验证码不正确'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ],
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { 
      validator: (rule: any, value: string, callback: any) => {
        if (value === '') {
          callback(new Error('请再次输入新密码'))
        } else if (value !== formData.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
}

// 更新验证码
const updateCaptchaCode = (code: string) => {
  captchaCode.value = code
}

// 刷新验证码
const refreshCaptcha = () => {
  if (captchaRef.value) {
    captchaRef.value.refreshCaptcha()
  }
}

// 清除表单错误
const clearFormErrors = () => {
  formErrors.email = ''
  formErrors.captcha = ''
  formErrors.password = ''
  formErrors.confirmPassword = ''
  formErrors.general = ''
}

// 处理提交第一步
const handleStep1 = () => {
  if (!formRef.value) return
  
  formRef.value.validate().then(valid => {
    if (valid) {
      // 验证邮箱是否存在
      checkEmailExists()
    }
  })
}

// 检查邮箱是否存在
const checkEmailExists = async () => {
  try {
    loading.value = true
    // 这里应该调用API检查邮箱是否存在
    await axios.post('/auth/check-email', { email: formData.email })
    
    // 邮箱存在，进入验证步骤
    currentStep.value = 2
  } catch (error: any) {
    // 处理错误
    if (error.response && error.response.status === 404) {
      formErrors.email = '该邮箱未注册'
    } else {
      formErrors.general = error.response?.data?.error || '系统错误，请稍后再试'
    }
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

// 邮箱验证成功
const onEmailVerified = () => {
  isEmailVerified.value = true
  currentStep.value = 3
}

// 邮箱验证取消
const onEmailVerificationCancel = () => {
  currentStep.value = 1
}

// 处理重置密码
const handleResetPassword = () => {
  if (!formRef.value) return
  
  formRef.value.validate().then(valid => {
    if (valid) {
      resetPassword()
    }
  })
}

// 重置密码
const resetPassword = async () => {
  try {
    loading.value = true
    clearFormErrors()
    
    // 调用重置密码API
    await axios.post('/auth/reset-password', {
      email: formData.email,
      password: formData.password
    })
    
    ElMessage.success('密码重置成功，请用新密码登录')
    router.push('/auth/login')
  } catch (error: any) {
    formErrors.general = error.response?.data?.error || '密码重置失败，请稍后再试'
  } finally {
    loading.value = false
  }
}

// 返回上一步
const previousStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

// 返回登录页
const goToLogin = () => {
  router.push('/auth/login')
}
</script>

<template>
  <div class="forgot-password-box">
    <div class="forgot-password-header">
      <img src="@/assets/logo.png" alt="Logo" class="logo" />
      <h2>忘记密码</h2>
      <p class="subtitle">通过注册邮箱重置您的密码</p>
    </div>

    <div v-if="formErrors.general" class="error-alert">
      <el-alert
        :title="formErrors.general"
        type="error"
        show-icon
        :closable="false"
      />
    </div>
    
    <!-- 步骤条 -->
    <el-steps :active="currentStep" finish-status="success" class="steps">
      <el-step title="输入邮箱"></el-step>
      <el-step title="验证邮箱"></el-step>
      <el-step title="重置密码"></el-step>
    </el-steps>

    <!-- 第一步：输入邮箱 -->
    <el-form
      v-if="currentStep === 1"
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-position="top"
      @submit.prevent
      status-icon
    >
      <el-form-item 
        label="邮箱地址" 
        prop="email"
        :error="formErrors.email"
      >
        <el-input
          v-model="formData.email"
          :prefix-icon="Message"
          placeholder="请输入注册邮箱"
        />
      </el-form-item>
      
      <el-form-item 
        label="验证码" 
        prop="captcha"
        :error="formErrors.captcha"
      >
        <div class="captcha-wrapper">
          <el-input
            v-model="formData.captcha"
            placeholder="请输入验证码"
          />
          <CaptchaImage 
            ref="captchaRef"
            @update:code="updateCaptchaCode" 
            :width="120" 
            :height="40"
          />
        </div>
      </el-form-item>

      <el-button
        type="primary"
        class="submit-button"
        :loading="loading"
        @click="handleStep1"
      >
        下一步
      </el-button>

      <div class="login-link">
        <el-link type="primary" @click="goToLogin">返回登录</el-link>
      </div>
    </el-form>
    
    <!-- 第二步：验证邮箱 -->
    <div v-if="currentStep === 2" class="step-container">
      <EmailVerification 
        :email="formData.email" 
        @verified="onEmailVerified"
        @cancel="onEmailVerificationCancel"
      />
    </div>
    
    <!-- 第三步：重置密码 -->
    <el-form
      v-if="currentStep === 3"
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-position="top"
      @submit.prevent
      status-icon
    >
      <el-form-item 
        label="新密码" 
        prop="password"
        :error="formErrors.password"
      >
        <el-input
          v-model="formData.password"
          type="password"
          :prefix-icon="Lock"
          show-password
          placeholder="请输入新密码"
        />
      </el-form-item>

      <el-form-item 
        label="确认新密码" 
        prop="confirmPassword"
        :error="formErrors.confirmPassword"
      >
        <el-input
          v-model="formData.confirmPassword"
          type="password"
          :prefix-icon="Lock"
          show-password
          placeholder="请再次输入新密码"
        />
      </el-form-item>

      <div class="form-actions">
        <el-button @click="previousStep">上一步</el-button>
        <el-button
          type="primary"
          :loading="loading"
          @click="handleResetPassword"
        >
          重置密码
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<style scoped>
.forgot-password-box {
  width: 100%;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.forgot-password-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
}

h2 {
  margin: 0;
  font-size: 24px;
  color: #333;
  margin-bottom: 8px;
}

.subtitle {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.error-alert {
  margin-bottom: 20px;
}

.steps {
  margin-bottom: 30px;
}

.captcha-wrapper {
  display: flex;
  align-items: center;
}

.captcha-wrapper .el-input {
  margin-right: 10px;
  flex: 1;
}

.submit-button {
  width: 100%;
  height: 40px;
  font-size: 16px;
  margin-bottom: 16px;
}

.login-link {
  text-align: center;
  font-size: 14px;
  color: #666;
}

.step-container {
  margin-top: 20px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 16px;
}

:deep(.el-input__wrapper) {
  padding: 0 15px;
}

:deep(.el-input__inner) {
  height: 40px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #333;
}
</style> 