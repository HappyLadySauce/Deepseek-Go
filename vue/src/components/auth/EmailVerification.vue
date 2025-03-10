<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import axios from '@/axios'

const props = defineProps<{
  email: string
}>()

const emit = defineEmits<{
  (e: 'verified'): void
  (e: 'cancel'): void
}>()

const verificationCode = ref('')
const inputCode = ref('')
const countdown = ref(0)
const loading = ref(false)
const verifying = ref(false)

// 监听邮箱变化，重置状态
watch(() => props.email, () => {
  verificationCode.value = ''
  inputCode.value = ''
  countdown.value = 0
  loading.value = false
}, { immediate: true })

// 发送验证码
const sendCode = async () => {
  if (countdown.value > 0) return
  if (!props.email) {
    ElMessage.warning('请先输入邮箱地址')
    return
  }
  
  try {
    loading.value = true
    console.log('发送验证码到邮箱:', props.email)
    const response = await axios.post('/auth/send-verification', {
      email: props.email,
      action: 'send-verification'
    })
    
    console.log('验证码发送成功:', response.data)
    ElMessage.success('验证码已发送，请查收邮件')
    // 启动倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error: any) {
    console.error('验证码发送失败:', error)
    ElMessage.error(error.response?.data?.error || '发送验证码失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 验证验证码
const verifyCode = async () => {
  if (!inputCode.value) {
    ElMessage.warning('请输入验证码')
    return
  }
  
  try {
    verifying.value = true
    console.log('验证邮箱验证码:', { email: props.email, code: inputCode.value })
    const response = await axios.post('/auth/verify-email', {
      email: props.email,
      code: inputCode.value
    })
    
    console.log('验证码验证成功:', response.data)
    ElMessage.success('邮箱验证成功')
    emit('verified')
  } catch (error: any) {
    console.error('验证码验证失败:', error)
    ElMessage.error(error.response?.data?.error || '验证码验证失败，请重试')
  } finally {
    verifying.value = false
  }
}

// 取消验证
const cancel = () => {
  emit('cancel')
}
</script>

<template>
  <div class="email-verification-container">
    <h3>邮箱验证</h3>
    <p class="description">
      我们已向 <strong>{{ email }}</strong> 发送了验证码，请查收并在下方输入
    </p>
    
    <div class="verification-input">
      <el-input
        v-model="inputCode"
        placeholder="请输入验证码"
        maxlength="6"
        class="verification-code-input"
      />
      <el-button 
        type="primary" 
        @click="sendCode" 
        :disabled="countdown > 0 || loading"
        :loading="loading"
        class="send-button"
      >
        {{ countdown > 0 ? `重新发送(${countdown}s)` : '发送验证码' }}
      </el-button>
    </div>
    
    <div class="action-buttons">
      <el-button @click="cancel">取消</el-button>
      <el-button 
        type="primary" 
        @click="verifyCode" 
        :loading="verifying"
      >
        验证
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.email-verification-container {
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 500px;
  margin: 0 auto;
}

h3 {
  margin-top: 0;
  margin-bottom: 16px;
  text-align: center;
  color: #303133;
  font-size: 18px;
}

.description {
  margin-bottom: 24px;
  color: #606266;
  font-size: 14px;
  text-align: center;
  line-height: 1.5;
}

.verification-input {
  display: flex;
  margin-bottom: 24px;
  align-items: center;
}

.verification-code-input {
  flex: 1;
  margin-right: 10px;
}

.send-button {
  width: 120px;
  white-space: nowrap;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 10px;
}
</style> 