<script setup lang="ts">
import { User, Lock, Message } from '@element-plus/icons-vue'
import useRegister from '@/hooks/auth/useRegister'
import CaptchaImage from '@/components/auth/CaptchaImage.vue'
import EmailVerification from '@/components/auth/EmailVerification.vue'

const { 
  registerForm, 
  registerRules, 
  register, 
  loading, 
  formRef,
  goToLogin,
  formErrors,
  captchaCode,
  updateCaptchaCode,
  captchaRef,
  currentStep,
  previousStep,
  emailVerified,
  isEmailVerified
} = useRegister()

// 调试函数
const testRegister = () => {
  console.log('测试注册按钮被点击')
  console.log('当前步骤:', currentStep.value)
  console.log('邮箱是否验证:', isEmailVerified.value)
  register() // 调用原始注册函数
}
</script>

<template>
  <div class="register-box">
    <div class="register-header">
      <img src="@/assets/logo.png" alt="Logo" class="logo" />
      <h2>创建账号</h2>
      <p class="subtitle">加入DeepSeek，开启智能对话之旅</p>
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
    <el-steps :active="currentStep" finish-status="success" class="register-steps" align-center>
      <el-step title="填写信息"></el-step>
      <el-step title="验证邮箱"></el-step>
    </el-steps>

    <!-- 第一步：基本信息 -->
    <el-form
      v-if="currentStep === 1"
      ref="formRef"
      :model="registerForm"
      :rules="registerRules"
      label-position="top"
      @submit.prevent
      status-icon
    >
      <el-form-item 
        label="用户名" 
        prop="username"
        :error="formErrors.username"
      >
        <el-input
          v-model="registerForm.username"
          :prefix-icon="User"
          placeholder="请输入用户名"
        />
      </el-form-item>

      <el-form-item 
        label="邮箱" 
        prop="email"
        :error="formErrors.email"
      >
        <el-input
          v-model="registerForm.email"
          :prefix-icon="Message"
          placeholder="请输入邮箱"
        />
      </el-form-item>

      <el-form-item 
        label="密码" 
        prop="password"
        :error="formErrors.password"
      >
        <el-input
          v-model="registerForm.password"
          type="password"
          :prefix-icon="Lock"
          show-password
          placeholder="请输入密码"
          :default-value="registerForm.password"
        />
      </el-form-item>

      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input
          v-model="registerForm.confirmPassword"
          type="password"
          :prefix-icon="Lock"
          show-password
          placeholder="请再次输入密码"
          :default-value="registerForm.confirmPassword"
        />
      </el-form-item>
      
      <el-form-item 
        label="验证码" 
        prop="captcha"
        :error="formErrors.captcha"
      >
        <div class="captcha-wrapper">
          <el-input
            v-model="registerForm.captcha"
            placeholder="请输入验证码"
            :default-value="registerForm.captcha"
          />
          <CaptchaImage 
            ref="captchaRef"
            @update:code="updateCaptchaCode" 
            :width="120" 
            :height="40"
            :default-value="registerForm.captcha"
          />
        </div>
      </el-form-item>

      <el-form-item prop="agreeTerms">
        <el-checkbox v-model="registerForm.agreeTerms">
          我已阅读并同意
          <el-link type="primary">服务条款</el-link>
          和
          <el-link type="primary">隐私政策</el-link>
        </el-checkbox>
      </el-form-item>

      <el-button
        type="primary"
        class="register-button"
        :loading="loading"
        @click="register"
      >
        下一步
      </el-button>

      <div class="login-link">
        已有账号？
        <el-link type="primary" @click="goToLogin">立即登录</el-link>
      </div>
    </el-form>
    
    <!-- 第二步：邮箱验证 -->
    <div v-if="currentStep === 2" class="step-container">
      <!-- 修改邮箱验证为true -->
      <div v-if="!isEmailVerified" class="verification-container">
        <EmailVerification 
          :email="registerForm.email" 
          @verified="emailVerified"
          @cancel="previousStep"
        />
      </div>
      <div v-else class="verification-success">
        <el-result
          icon="success"
          title="邮箱验证成功"
          sub-title="您的邮箱已成功验证，点击下方按钮完成注册"
        >
          <template #extra>
            <el-button 
              type="primary" 
              @click="register" 
              :loading="loading"
            >
              完成注册
            </el-button>
            <el-button @click="previousStep">返回修改</el-button>
          </template>
        </el-result>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-box {
  width: 100%;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.register-header {
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

.register-steps {
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

.register-button {
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

.verification-success {
  text-align: center;
  margin: 20px 0;
}

.verification-container {
  margin-bottom: 20px;
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

:deep(.el-checkbox__label) {
  color: #666;
}
</style> 