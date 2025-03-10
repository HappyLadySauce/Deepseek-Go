import { reactive, ref } from 'vue'
import axios from '@/axios'

export default function () {
    // 获取token
    const token = ref<string | null>(localStorage.getItem('token'))

    // 表单
    const form = reactive({
        username: '',
        password: ''
    })
    
    // 表单验证规则
    const rules = reactive({
        username: [
            { required: true, message: '请输入用户名或者邮箱', trigger: 'blur' }
        ],
        password: [
            { required: true, message: '请输入密码', trigger: 'blur' }
        ]
    })

    // 登录
    const login = () => {
        const { username, password } = form
        axios.post('/auth/login', { username, password }).then(res => {
            console.log(res)
        })
        .catch(err => {
            console.log(err)
        })
    }

    // 返回
    return { form, rules, login }
}
