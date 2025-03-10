<script setup lang="ts">
import { User, Lock } from '@element-plus/icons-vue'
import useLogin from '@/hooks/auth/useLogin'
import CaptchaImage from '@/components/auth/CaptchaImage.vue'

const { 
  loginForm, 
  loginRules, 
  loading, 
  formRef,
  handleLogin,
  goToRegister,
  formErrors,
  handleForgotPassword,
  captchaCode,
  updateCaptchaCode,
  captchaRef
} = useLogin()
</script>

<template>
  <div class="login-box">
    <div class="login-header">
      <img src="@/assets/logo.png" alt="Logo" class="logo" />
      <h2>欢迎回来</h2>
      <p class="subtitle">登录您的DeepSeek账号</p>
    </div>

    <div v-if="formErrors.general" class="error-alert">
      <el-alert
        :title="formErrors.general"
        type="error"
        show-icon
        :closable="false"
      />
    </div>

    <el-form
      ref="formRef"
      :model="loginForm"
      :rules="loginRules"
      label-position="top"
      @submit.prevent
      status-icon
    >
      <el-form-item 
        label="用户名或邮箱" 
        prop="username"
        :error="formErrors.username"
      >
        <el-input
          v-model="loginForm.username"
          :prefix-icon="User"
          placeholder="请输入用户名或邮箱"
          @keyup.enter="handleLogin"
        />
      </el-form-item>

      <el-form-item 
        label="密码" 
        prop="password"
        :error="formErrors.password"
      >
        <el-input
          v-model="loginForm.password"
          type="password"
          :prefix-icon="Lock"
          placeholder="请输入密码"
          show-password
          @keyup.enter="handleLogin"
        />
      </el-form-item>
      
      <el-form-item 
        label="验证码" 
        prop="captcha"
        :error="formErrors.captcha"
      >
        <div class="captcha-wrapper">
          <el-input
            v-model="loginForm.captcha"
            placeholder="请输入验证码"
            @keyup.enter="handleLogin"
          />
          <CaptchaImage 
            ref="captchaRef"
            @update:code="updateCaptchaCode" 
            :width="120" 
            :height="40"
          />
        </div>
      </el-form-item>

      <div class="remember-forgot">
        <el-checkbox v-model="loginForm.remember">记住密码</el-checkbox>
        <el-link type="primary" @click="handleForgotPassword">忘记密码？</el-link>
      </div>

      <el-button
        type="primary"
        class="login-button"
        :loading="loading"
        @click="handleLogin"
      >
        登录
      </el-button>

      <div class="register-link">
        还没有账号？
        <el-link type="primary" @click="goToRegister">立即注册</el-link>
      </div>
    </el-form>
  </div>
</template>

<style scoped>
.login-box {
  width: 100%;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.login-header {
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

.captcha-wrapper {
  display: flex;
  align-items: center;
}

.captcha-wrapper .el-input {
  margin-right: 10px;
  flex: 1;
}

.remember-forgot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.login-button {
  width: 100%;
  height: 40px;
  font-size: 16px;
  margin-bottom: 16px;
}

.register-link {
  text-align: center;
  font-size: 14px;
  color: #666;
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