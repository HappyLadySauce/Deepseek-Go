import { ElMessage } from 'element-plus'

/**
 * 表单错误结构
 */
export interface FormErrors {
  general?: string
  [key: string]: string | undefined
}

/**
 * 处理API错误并更新表单错误状态
 * @param error API错误对象
 * @param formErrors 表单错误状态对象
 * @param refreshCaptcha 刷新验证码的回调函数
 * @returns 更新后的表单错误
 */
export const handleApiError = (error: any, formErrors: FormErrors, refreshCaptcha?: () => void): FormErrors => {
  // 重置先前的错误
  Object.keys(formErrors).forEach(key => {
    formErrors[key] = undefined
  })

  // 网络错误处理
  if (!error.response) {
    formErrors.general = '网络错误，请稍后重试'
    return formErrors
  }

  const status = error.response.status
  const data = error.response.data

  // 显示错误消息
  if (data && data.error) {
    // 显示后端返回的错误消息
    ElMessage.error(`[${status}] ${data.error}`)
    
    // 根据错误消息内容决定将错误显示在哪个字段
    const errorMessage = data.error.toLowerCase()
    
    if (errorMessage.includes('用户名') || errorMessage.includes('username')) {
      formErrors.username = data.error
    } else if (errorMessage.includes('邮箱') || errorMessage.includes('email')) {
      formErrors.email = data.error
    } else if (errorMessage.includes('密码') || errorMessage.includes('password')) {
      formErrors.password = data.error
    } else if (errorMessage.includes('验证码') || errorMessage.includes('captcha')) {
      formErrors.captcha = data.error
      if (refreshCaptcha) refreshCaptcha()
    } else {
      // 如果无法确定具体字段，则显示为通用错误
      formErrors.general = data.error
    }
  } else {
    // 如果没有具体错误信息，根据状态码显示通用错误
    switch (status) {
      case 401:
        formErrors.general = '未授权，请重新登录'
        if (refreshCaptcha) refreshCaptcha()
        break
      case 403:
        formErrors.general = '没有权限执行此操作'
        break
      case 404:
        formErrors.general = '请求的资源不存在'
        break
      case 500:
        formErrors.general = '服务器内部错误'
        break
      default:
        formErrors.general = `未知错误，状态码: ${status}`
    }
  }

  return formErrors
} 