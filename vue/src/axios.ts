import axios from 'axios';
import type { InternalAxiosRequestConfig } from 'axios';
import { ElMessage } from 'element-plus';
import router from './router'

// 创建axios实例
const instance = axios.create({
  baseURL: "http://localhost:14020/api/v1", // 基础URL，添加/api/v1
  timeout: 10000, // 增加请求超时时间
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    // 如果token存在，则添加到请求头
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
instance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (!error.response) {
      // 网络错误或请求被取消
      console.error('网络错误或请求被取消:', error.message);
      ElMessage.error('网络连接失败，请检查网络设置');
    } else {
      // 服务器返回错误
      console.error('服务器错误:', error.response.status, error.response.data);
      
      // 根据状态码处理错误
      switch (error.response.status) {
        case 401:
          // token过期或无效
          localStorage.removeItem('token');
          ElMessage.error('登录已过期，请重新登录');
          router.push('/auth/login');
          break;
        case 403:
          ElMessage.error('没有权限访问');
          break;
        case 404:
          ElMessage.error('请求的资源不存在');
          break;
        case 500:
          ElMessage.error('服务器错误，请稍后再试');
          break;
        default:
          // 使用服务器返回的错误信息
          const errorMsg = error.response.data?.error || '请求失败，请稍后再试';
          ElMessage.error(errorMsg);
      }
    }
    return Promise.reject(error);
  }
);

// 模拟刷新 Token 的接口
async function refreshTokenAPI(refreshToken: string): Promise<string> {
  // 替换为实际的刷新 Token 接口
  const response = await axios.post<{ token: string }>("/api/refresh-token", {
    refreshToken,
  });
  return response.data.token;
}

export default instance;