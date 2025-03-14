<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

// 获取路由参数中的邮箱
const email = ref(route.query.email as string || '')

// 表单引用
const verifyFormRef = ref<FormInstance>()

// 表单数据
const verifyForm = reactive({
  email: email.value,
  code: ''
})

// 表单验证规则
const verifyRules = reactive<FormRules>({
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
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
    await verifyFormRef.value?.validateField('email')
    
    codeSending.value = true
    
    // 记录API请求详情，便于调试
    console.log(`准备发送验证码到邮箱: ${verifyForm.email}`)
    
    const success = await userStore.sendVerificationCode(verifyForm.email)
    
    if (success) {
      ElMessage.success('验证码已发送，请检查您的邮箱')
      startCountDown()
    }
  } catch (error) {
    // 邮箱验证失败
    console.error('发送验证码失败:', error)
    ElMessage.error('发送验证码失败，请稍后重试')
  } finally {
    codeSending.value = false
  }
}

// 验证邮箱
const handleVerify = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  
  await formEl.validate((valid) => {
    if (valid) {
      userStore.verifyCode(verifyForm.email, verifyForm.code).then(success => {
        if (success) {
          ElMessage.success('邮箱验证成功')
          router.push('/register?email=' + verifyForm.email)
        }
      })
    } else {
      ElMessage.error('请正确填写验证表单')
      return false
    }
  })
}

// 返回注册页
const goToRegister = () => {
  router.push('/register')
}
</script>

<template>
  <div class="verify-container">
    <div class="verify-box">
      <div class="verify-header">
        <h2>验证您的邮箱</h2>
        <p>请输入您收到的验证码</p>
      </div>
      
      <el-form
        ref="verifyFormRef"
        :model="verifyForm"
        :rules="verifyRules"
        label-position="top"
        class="verify-form"
      >
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="verifyForm.email"
            placeholder="请输入邮箱"
            type="email"
          />
        </el-form-item>
        
        <el-form-item label="验证码" prop="code">
          <div class="verification-code">
            <el-input
              v-model="verifyForm.code"
              placeholder="请输入验证码"
              class="code-input"
            />
            <el-button
              type="primary"
              :disabled="codeSending || countDown > 0"
              @click="sendVerificationCode"
              class="send-code-button"
            >
              {{ countDown > 0 ? `${countDown}秒后重新发送` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
        
        <div class="form-actions">
          <el-button
            type="primary"
            :loading="userStore.loading"
            @click="handleVerify(verifyFormRef)"
            class="verify-button"
          >
            验证邮箱
          </el-button>
        </div>
      </el-form>
      
      <div class="verify-footer">
        <p>已有账号？<el-button type="text" @click="router.push('/login')">登录</el-button></p>
        <p>返回注册？<el-button type="text" @click="goToRegister">注册</el-button></p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.verify-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.verify-box {
  width: 100%;
  max-width: 400px;
  padding: 30px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.verify-header {
  text-align: center;
  margin-bottom: 30px;
}

.verify-header h2 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 10px;
}

.verify-header p {
  font-size: 14px;
  color: #909399;
}

.verify-form {
  margin-bottom: 30px;
}

.verification-code {
  display: flex;
  gap: 10px;
}

.code-input {
  flex: 1;
}

.form-actions {
  margin-top: 30px;
}

.verify-button {
  width: 100%;
  padding: 12px 0;
}

.verify-footer {
  text-align: center;
  font-size: 14px;
  color: #606266;
}
</style> 