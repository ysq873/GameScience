import request from '@/utils/request'

export function getProfile() {
  return request({ url: '/user/profile', method: 'get' })
}

export function updateProfile(data) {
  // 注意：Kratos 中更新用户资料需要不同的API
  // 这里只是一个占位符，实际实现需要根据Kratos API来调整
  return request({
    url: '/user/profile',
    method: 'put',
    data
  })
}

export function addFavorite(item) {
  return request({
    url: '/user/favorites',
    method: 'post',
    data: { item }
  })
}
