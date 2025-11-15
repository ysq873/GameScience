import apiRequest from '@/utils/request'

export function addFavorite(modelId) {
  return apiRequest.post('/models/favorite', null, { params: { model_id: modelId } })
}

export function removeFavorite(modelId) {
  return apiRequest.delete('/models/favorite', { params: { model_id: modelId } })
}

export function listFavorites() {
  return apiRequest.get('/user/favorites/models')
}