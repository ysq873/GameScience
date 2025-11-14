import apiRequest from '@/utils/request'

export function listPurchases() {
  return apiRequest.get('/purchases')
}

export function generateDownloadToken(modelId) {
  return apiRequest.post('/models/download-token', null, { params: { model_id: modelId } })
}
