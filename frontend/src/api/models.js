import apiRequest from '@/utils/request'

export function listModels(params) {
  return apiRequest.get('/models', { params })
}

export function getModel(id) {
  return apiRequest.get('/models/detail', { params: { id } })
}

export function updateStatus(id, status) {
  return apiRequest.put('/models/status', { status }, { params: { id } })
}

export function uploadModel(formData) {
  return apiRequest.post('/models/upload', formData)
}

export function getMyModels(params) {
  return apiRequest.get('/user/models', { params })
}
