import axios from 'axios'

const instance = axios.create({
    baseURL: 'http://localhost:14020/api/v1',
})

// 请求拦截器
instance.interceptors.request.use(config => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// 响应拦截器
// instance.interceptors.response.use(response => {
//     return response.data
// }, error => {
//     return Promise.reject(error)
// })

// 导出实例
export default instance
