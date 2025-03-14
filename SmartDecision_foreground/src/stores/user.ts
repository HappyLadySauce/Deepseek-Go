import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import router from '../router'

// 配置axios默认设置
axios.defaults.headers.post['Content-Type'] = 'application/json'
axios.defaults.headers.put['Content-Type'] = 'application/json'

// 定义类型
interface UserInfo {
  username: string
  email: string
  token: string
}

interface LoginData {
  username: string
  password: string
}

interface RegisterData {
  username: string
  email: string
  password: string
  verificationCode: string
}

interface ResetPasswordData {
  email: string
  code: string
  newPassword: string
}

interface UpdateProfileData {
  username: string
  email?: string
  currentPassword?: string
  newPassword?: string
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const userInfo = ref<UserInfo | null>(null)
  const isAuthenticated = ref(false)
  const loading = ref(false)

  // API基础URL
  const baseURL = 'http://localhost:14020/api/v1'

  // 登录
  const login = async (loginData: LoginData): Promise<boolean> => {
    loading.value = true
    try {
      console.log('登录请求:', {
        url: `${baseURL}/auth/login`,
        data: loginData
      })
      
      const response = await axios.post(`${baseURL}/auth/login`, loginData)
      
      console.log('登录成功:', response.data)
      
      // 存储用户信息和令牌
      userInfo.value = {
        username: response.data.username,
        email: response.data.email,
        token: response.data.token
      }
      
      // 设置认证状态
      isAuthenticated.value = true
      
      // 存储令牌到本地存储
      localStorage.setItem('token', response.data.token)
      localStorage.setItem('username', response.data.username)
      
      // 设置axios默认请求头
      axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.token}`
      
      // 仅在UI中显示成功消息，不在控制台重复记录
      ElMessage.success('登录成功')
      return true
    } catch (error: any) {
      console.error('登录失败:', error)
      
      if (error.response) {
        console.error('错误状态码:', error.response.status)
        console.error('错误数据:', error.response.data)
        // 使用后端返回的错误消息，如果没有则使用通用错误消息
        const errorMsg = error.response?.data?.error || '登录失败，请检查用户名和密码'
        ElMessage.error(errorMsg)
      } else {
        ElMessage.error('登录失败，请检查网络连接')
      }
      
      return false
    } finally {
      loading.value = false
    }
  }

  // 注册
  const register = async (registerData: RegisterData): Promise<boolean> => {
    loading.value = true
    try {
      console.log('注册请求:', {
        url: `${baseURL}/auth/register`,
        data: registerData
      })
      
      const response = await axios.post(`${baseURL}/auth/register`, registerData)
      
      console.log('注册成功:', response.data)
      // 只显示一条成功消息
      ElMessage.success('注册成功，请登录')
      return true
    } catch (error: any) {
      console.error('注册失败:', error)
      
      if (error.response) {
        console.error('错误状态码:', error.response.status)
        console.error('错误数据:', error.response.data)
        // 使用后端返回的错误消息，如果没有则使用通用错误消息
        const errorMsg = error.response?.data?.error || '注册失败，请检查输入信息'
        ElMessage.error(errorMsg)
      } else {
        ElMessage.error('注册失败，请检查网络连接')
      }
      
      return false
    } finally {
      loading.value = false
    }
  }

  // 发送验证码
  const sendVerificationCode = async (email: string): Promise<boolean> => {
    loading.value = true
    try {
      console.log('发送验证码请求:', {
        url: `${baseURL}/auth/send-verification-email`,
        data: { email }
      })
      
      const response = await axios.post(`${baseURL}/auth/send-verification-email`, { email })
      
      console.log('验证码发送成功:', response.data)
      // 使用后端返回的消息或默认消息
      ElMessage.success(response.data?.message || '验证码已发送到您的邮箱')
      return true
    } catch (error: any) {
      console.error('验证码发送失败:', error)
      
      if (error.response) {
        console.error('错误状态码:', error.response.status)
        console.error('错误数据:', error.response.data)
        ElMessage.error(error.response?.data?.error || `发送验证码失败`)
      } else if (error.request) {
        console.error('没有收到响应:', error.request)
        ElMessage.error('服务器无响应，请检查网络连接')
      } else {
        console.error('请求配置错误:', error.message)
        ElMessage.error(`发送请求出错`)
      }
      
      return false
    } finally {
      loading.value = false
    }
  }

  // 验证验证码
  const verifyCode = async (email: string, code: string): Promise<boolean> => {
    loading.value = true
    try {
      console.log('验证验证码请求:', {
        url: `${baseURL}/auth/verify-verification-code`,
        data: { email, code }
      })
      
      const response = await axios.post(`${baseURL}/auth/verify-verification-code`, { email, code })
      
      console.log('验证码验证成功:', response.data)
      // 使用后端返回的消息
      if (response.data?.valid) {
        ElMessage.success(response.data?.message || '验证码验证成功')
        return true
      } else {
        ElMessage.error('验证码无效')
        return false
      }
    } catch (error: any) {
      console.error('验证码验证失败:', error)
      
      if (error.response) {
        console.error('错误状态码:', error.response.status)
        console.error('错误数据:', error.response.data)
        ElMessage.error(error.response?.data?.error || '验证码验证失败')
      } else if (error.request) {
        console.error('没有收到响应:', error.request)
        ElMessage.error('服务器无响应，请检查网络连接')
      } else {
        console.error('请求配置错误:', error.message)
        ElMessage.error('验证码验证失败')
      }
      
      return false
    } finally {
      loading.value = false
    }
  }

  // 注销
  const logout = () => {
    // 清除用户信息
    userInfo.value = null
    isAuthenticated.value = false
    
    // 清除本地存储
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    
    // 清除请求头
    delete axios.defaults.headers.common['Authorization']
    
    // 重定向到登录页
    router.push('/login')
    ElMessage.success('已退出登录')
  }

  // 检查认证状态
  const checkAuth = () => {
    const token = localStorage.getItem('token')
    const username = localStorage.getItem('username')
    
    if (token && username) {
      isAuthenticated.value = true
      userInfo.value = {
        username,
        email: '', // 从本地存储中我们没有email，可以通过API获取
        token
      }
      
      // 设置axios默认请求头
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
    }
  }

  // 重置密码
  const resetPassword = async (resetData: ResetPasswordData): Promise<boolean> => {
    loading.value = true
    try {
      console.log('重置密码请求:', {
        url: `${baseURL}/auth/reset-password`,
        data: resetData
      })
      
      const response = await axios.post(`${baseURL}/auth/reset-password`, resetData)
      
      console.log('重置密码成功:', response.data)
      ElMessage.success('密码重置成功')
      return true
    } catch (error: any) {
      console.error('重置密码失败:', error)
      
      if (error.response) {
        console.error('错误状态码:', error.response.status)
        console.error('错误数据:', error.response.data)
        ElMessage.error(error.response?.data?.error || '密码重置失败')
      } else {
        ElMessage.error('密码重置失败，请检查网络连接')
      }
      
      return false
    } finally {
      loading.value = false
    }
  }

  // 更新个人信息
  const updateProfile = async (profileData: UpdateProfileData): Promise<boolean> => {
    loading.value = true
    try {
      console.log('更新个人信息请求:', {
        url: `${baseURL}/auth/update-profile`,
        data: profileData
      })
      
      const response = await axios.put(`${baseURL}/auth/update-profile`, profileData)
      
      console.log('更新个人信息成功:', response.data)
      
      // 更新本地存储的用户信息
      if (userInfo.value) {  // 先检查userInfo.value是否为null
        userInfo.value = {
          ...userInfo.value,
          username: response.data.username,
          email: response.data.email,
          token: userInfo.value.token  // 保留现有token
        }
        
        localStorage.setItem('username', response.data.username)
      }
      
      ElMessage.success('个人信息更新成功')
      return true
    } catch (error: any) {
      console.error('更新个人信息失败:', error)
      
      if (error.response) {
        console.error('错误状态码:', error.response.status)
        console.error('错误数据:', error.response.data)
        ElMessage.error(error.response?.data?.error || '个人信息更新失败')
      } else {
        ElMessage.error('个人信息更新失败，请检查网络连接')
      }
      
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    userInfo,
    isAuthenticated,
    loading,
    login,
    register,
    sendVerificationCode,
    verifyCode,
    logout,
    checkAuth,
    resetPassword,
    updateProfile
  }
}) 