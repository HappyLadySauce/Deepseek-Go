import axios from '@/axios'

/**
 * 用户注册请求参数
 */
export interface RegisterParams {
  username: string
  email: string
  password: string
  captcha: string
  captcha_code: string
}

/**
 * 用户登录请求参数
 */
export interface LoginParams {
  username: string
  password: string
  captcha: string
  captcha_code: string
  remember?: boolean
}

/**
 * API服务类，处理与后端的通信
 */
export default {
  /**
   * 用户注册
   * @param params 注册参数
   * @returns 注册结果
   */
  registerUser(params: RegisterParams) {
    return axios.post('/auth/register', params)
  },

  /**
   * 用户登录
   * @param params 登录参数
   * @returns 登录结果
   */
  loginUser(params: LoginParams) {
    return axios.post('/auth/login', params)
  },

  /**
   * 用户注销
   * @returns 注销结果
   */
  logoutUser() {
    return axios.post('/auth/logout')
  },

  /**
   * 发送邮箱验证码
   * @param email 邮箱地址
   * @returns 发送结果
   */
  sendEmailVerification(email: string) {
    return axios.post('/auth/send-verification', { email })
  },

  /**
   * 验证邮箱验证码
   * @param email 邮箱地址
   * @param code 验证码
   * @returns 验证结果
   */
  verifyEmailCode(email: string, code: string) {
    return axios.post('/auth/verify-email', { email, code })
  },

  /**
   * 请求重置密码
   * @param email 邮箱地址
   * @returns 请求结果
   */
  requestPasswordReset(email: string) {
    return axios.post('/auth/forgot-password', { email })
  },

  /**
   * 重置密码
   * @param token 重置令牌
   * @param password 新密码
   * @returns 重置结果
   */
  resetPassword(token: string, password: string) {
    return axios.post('/auth/reset-password', { token, password })
  }
} 