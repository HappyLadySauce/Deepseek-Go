import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import type { FormInstance } from 'element-plus'
import api from '@/utils/api'
import type { FormErrors } from '@/utils/errorHandler'
import { handleApiError } from '@/utils/errorHandler'

export default function () {
    const router = useRouter()
    const formRef = ref<FormInstance>()
    // 获取token
    const token = ref<string | null>(localStorage.getItem('token'))
    const loading = ref(false)

    // 表单
    const loginForm = reactive({
        username: '',
        password: '',
        remember: false,
        captcha: '',
        captcha_code: ''
    })
    
    // 图形验证码相关
    const captchaCode = ref('')
    const captchaRef = ref()

    // 表单错误信息
    const formErrors = reactive<FormErrors>({
        username: '',
        password: '',
        captcha: '',
        general: ''
    })
    
    // 从localStorage获取保存的账号密码
    const savedCredentials = localStorage.getItem('credentials')
    if (savedCredentials) {
        try {
            const { username, password } = JSON.parse(savedCredentials)
            loginForm.username = username
            loginForm.password = password
            loginForm.remember = true
        } catch (e) {
            localStorage.removeItem('credentials')
        }
    }

    // 验证规则
    const loginRules = {
        username: [
            { required: true, message: '请输入用户名或邮箱', trigger: 'blur' }
        ],
        password: [
            { required: true, message: '请输入密码', trigger: 'blur' },
            { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
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
        Object.keys(formErrors).forEach(key => {
            formErrors[key] = ''
        })
    }

    // 处理登录
    const handleLogin = () => {
        if (!formRef.value) return

        formRef.value.validate().then(valid => {
            if (valid) {
                clearFormErrors()
                login()
            } else {
                ElMessage.error('请正确填写登录信息')
            }
        })
    }

    // 登录
    const login = async () => {
        try {
            loading.value = true
            
            // 登录请求
            const response = await api.loginUser({
                username: loginForm.username,
                password: loginForm.password,
                captcha: loginForm.captcha,
                captcha_code: captchaCode.value,
                remember: loginForm.remember
            })

            // 保存token
            const token = response.data.token
            localStorage.setItem('token', token)
            
            // 记住账号密码
            if (loginForm.remember) {
                localStorage.setItem('credentials', JSON.stringify({
                    username: loginForm.username,
                    password: loginForm.password
                }))
            } else {
                localStorage.removeItem('credentials')
            }

            // 登录成功提示
            ElMessage.success('登录成功')
            
            // 跳转到首页或上一页
            const redirect = router.currentRoute.value.query.redirect as string
            router.push(redirect || '/overview')
        } catch (error: any) {
            // 使用统一的错误处理工具
            handleApiError(error, formErrors, refreshCaptcha)
        } finally {
            loading.value = false
        }
    }
    
    // 登出
    const logout = () => {
        // 清除token
        localStorage.removeItem('token')
        
        // 清除用户信息
        // 如果有用户信息存储，这里需要清除
        
        // 提示
        ElMessage.success('已退出登录')
        
        // 跳转到登录页
        router.push('/auth/login')
    }

    // 跳转到注册页
    const goToRegister = () => {
        router.push('/auth/register')
    }

    // 处理忘记密码
    const handleForgotPassword = () => {
        router.push('/auth/forgot-password')
    }

    return {
        loginForm,
        loginRules,
        login,
        loading,
        formRef,
        handleLogin,
        goToRegister,
        handleForgotPassword,
        formErrors,
        logout,
        captchaCode,
        updateCaptchaCode,
        refreshCaptcha,
        captchaRef
    }
}
