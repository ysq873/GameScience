import apiRequest from '@/utils/request'

export function createOrder(items) {
  return apiRequest.post('/orders', { items })
}

export function payOrder(id) {
  return apiRequest.post('/orders/pay', null, { params: { id } })
}

export function listOrders() {
  return apiRequest.get('/orders')
}

export function getOrderDetail(id) {
  return apiRequest.get('/orders/detail', { params: { id } })
}

export function refundOrder(id) {
  return apiRequest.post('/orders/refund', null, { params: { id } })
}

export function mockCallback(orderId, status, idempotencyKey) {
  return apiRequest.post('/payments/callback', { orderId, status, idempotencyKey })
}
