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
