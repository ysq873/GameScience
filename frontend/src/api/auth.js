// Kratos API service
import axios from 'axios'
import { ElMessage } from 'element-plus'

const kratosRequest = axios.create({
  timeout: 10000,
  withCredentials: true // Important for Kratos session handling
})

kratosRequest.interceptors.response.use(
  response => {
    return response
  },
  error => {
    // Only show error messages for non-401 errors, as 401s are expected in auth flows
    if (error.response?.status !== 401) {
      ElMessage.error(error.response?.data?.error?.message || error.response?.data?.message || '请求失败')
    }
    return Promise.reject(error)
  }
)

// 初始化注册流程
export function initRegistrationFlow() {
  return kratosRequest.get('/self-service/registration/browser', {
    headers: { Accept: 'application/json' }
  })
}

// 提交注册
export function register(flowId, data) {
  return kratosRequest.post(`/self-service/registration?flow=${flowId}`, data)
}

// 初始化登录流程
export function initLoginFlow() {
  return kratosRequest.get('/self-service/login/browser', {
    headers: { Accept: 'application/json' }
  })
}

export const getLoginFlow = (id) =>
  kratosRequest.get('/self-service/login/flows', { params: { id }, headers: { Accept: 'application/json' } })

// 提交登录
export function submitLogin(flowId, data) {
  return kratosRequest.post(`/self-service/login?flow=${flowId}`, data)
}

// 获取当前会话
export function getSession() {
  return kratosRequest.get('/sessions/whoami')
}

// 登出
export function logout() {
  return kratosRequest.get('/self-service/logout/browser')
}

// 初始化恢复流程
export function initRecoveryFlow() {
  return kratosRequest.get('/self-service/recovery/browser', {
    headers: { Accept: 'application/json' }
  })
}


export function recovery(flowId, data) {
  return kratosRequest.post(`/self-service/recovery?flow=${flowId}`, data)
}

export function initSettings(flowId) {
  return kratosRequest.get(`/self-service/settings/flows`, {
    params: { id: flowId },
    headers: { Accept: 'application/json' }
  })
}

export function settings(flowId, data) {
  return kratosRequest.post(`/self-service/settings?flow=${flowId}`, data)
}

// 获取用户档案
export function getProfile() {
  return getSession()
}
