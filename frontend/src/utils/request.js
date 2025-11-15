import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建两个实例：一个用于 Kratos，一个用于业务 API
export const kratosRequest = axios.create({
  timeout: 10000,
  withCredentials: true // 关键：允许携带 Cookie
})

export const apiRequest = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true
})

// Kratos 请求拦截器
kratosRequest.interceptors.response.use(
  response => response,
  error => {
    const status = error.response?.status
    if (status === 401) {
      ElMessage.error('请先登录')
    } else {
      const message = error.response?.data?.error?.message || 
                     error.response?.data?.ui?.messages?.[0]?.text || 
                     '请求失败'
      ElMessage.error(message)
    }
    return Promise.reject(error)
  }
)

// 业务 API 请求拦截器
apiRequest.interceptors.response.use(
  response => response,
  error => {
    const status = error.response?.status
    if (status === 401) {
      ElMessage.error('请先登录')
      window.location.href = '/login'
    } else {
      ElMessage.error(error.response?.data?.message || '请求失败')
    }
    return Promise.reject(error)
  }
)

export default apiRequest