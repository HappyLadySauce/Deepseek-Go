import { reactive, ref } from 'vue'
import axios from '@/axios'
import type { FormInstance } from 'element-plus';

export default function useRegister() {
    const form = reactive({
        username: '',
        password: '',
        email: ''
    })
    
    const rules = reactive({
        username: [
            { required: true, message: '请输入用户名', trigger: 'blur' }
        ],
        password: [
            { required: true, message: '请输入密码', trigger: 'blur' }
        ],
        email: [
            { required: true, message: '请输入邮箱', trigger: 'blur' }
        ]
    })
    
    const formRef = ref<FormInstance>()

    const register = () => {
        formRef.value?.validate((valid: boolean) => {
            if (valid) {
                axios.post('/auth/register', form).then(res => {
                    console.log(res)
                })
            } else {
                console.log('error submit!')
            }
        })
    }

    return { form, rules, formRef, register }
}
