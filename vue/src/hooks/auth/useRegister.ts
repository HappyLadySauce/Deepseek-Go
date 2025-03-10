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
    const loading = ref(false)
    const isEmailVerified = ref(false)
    
    // 注册步骤: 1.填写基础信息 2.验证邮箱
    const currentStep = ref(1)

    // 表单
    const registerForm = reactive({
        username: '',
        email: '',
        password: '',
        confirmPassword: '',
        agreeTerms: false,
        captcha: ''
    })
    
    // 图形验证码相关
    const captchaCode = ref('')
    const captchaRef = ref()

    // 表单错误信息
    const formErrors = reactive<FormErrors>({
        username: '',
        email: '',
        password: '',
        confirmPassword: '',
        captcha: '',
        general: ''
    })

    // 验证密码
    const validatePass = (rule: any, value: string, callback: any) => {
        if (value === '') {
            callback(new Error('请输入密码'))
        } else {
            if (registerForm.confirmPassword !== '') {
                if (formRef.value) {
                    formRef.value.validateField('confirmPassword')
                }
            }
            callback()
        }
    }

    // 验证确认密码
    const validatePass2 = (rule: any, value: string, callback: any) => {
        if (value === '') {
            callback(new Error('请再次输入密码'))
        } else if (value !== registerForm.password) {
            callback(new Error('两次输入密码不一致'))
        } else {
            callback()
        }
    }

    // 表单验证规则
    const registerRules = {
        username: [
            { required: true, message: '请输入用户名', trigger: 'blur' },
            { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' },
            { pattern: /^[a-zA-Z0-9_-]+$/, message: '用户名只能包含字母、数字、下划线和连字符', trigger: 'blur' }
        ],
        email: [
            { required: true, message: '请输入邮箱地址', trigger: 'blur' },
            { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        password: [
            { required: true, message: '请输入密码', trigger: 'blur' },
            { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
            { validator: validatePass, trigger: 'blur' }
        ],
        confirmPassword: [
            { required: true, message: '请再次输入密码', trigger: 'blur' },
            { validator: validatePass2, trigger: 'blur' }
        ],
        agreeTerms: [
            { required: true, message: '请阅读并同意服务条款', trigger: 'change' },
            { 
                validator: (rule: any, value: boolean, callback: any) => {
                    if (!value) {
                        callback(new Error('请阅读并同意服务条款'))
                    } else {
                        callback()
                    }
                }, 
                trigger: 'change' 
            }
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

    // 注册
    const register = async () => {
        // 第一步：表单验证
        if (currentStep.value === 1) {
            if (!formRef.value) {
                console.error('表单实例未找到')
                return
            }
            
            try {
                await formRef.value.validate()
                currentStep.value = 2 // 进入邮箱验证步骤
                return
            } catch (error) {
                ElMessage.error('请正确填写注册信息')
                return
            }
        }
        
        // 第二步：邮箱验证后注册
        if (currentStep.value === 2) {
            if (!isEmailVerified.value) {
                ElMessage.warning('请先验证邮箱')
                return
            }
            
            try {
                loading.value = true
                clearFormErrors()
                
                // 发送注册请求
                const response = await api.registerUser({
                    username: registerForm.username,
                    email: registerForm.email,
                    password: registerForm.password,
                    captcha: registerForm.captcha,
                    captcha_code: captchaCode.value
                })
                
                ElMessage.success('注册成功')
                router.push('/auth/login')
            } catch (error: any) {
                // 使用统一的错误处理工具
                handleApiError(error, formErrors, refreshCaptcha)
            } finally {
                loading.value = false
            }
        }
    }

    // 跳转到登录页
    const goToLogin = () => {
        router.push('/auth/login')
    }
    
    // 返回上一步
    const previousStep = () => {
        if (currentStep.value > 1) {
            currentStep.value -= 1
        }
    }
    
    // 邮箱验证完成
    const emailVerified = () => {
        isEmailVerified.value = true
        ElMessage.success('邮箱验证成功，请继续完成注册')
    }

    return {
        registerForm,
        registerRules,
        register,
        loading,
        formRef,
        goToLogin,
        formErrors,
        captchaCode,
        updateCaptchaCode,
        refreshCaptcha,
        captchaRef,
        currentStep,
        previousStep,
        emailVerified,
        isEmailVerified
    }
}
