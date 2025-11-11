import request from '@/utils/request'
import { getSession } from '@/api/auth'

// 使用 Kratos 会话获取用户档案
export async function getProfile() {
  try {
    const response = await getSession()
    if (response.data?.identity) {
      return {
        data: {
          id: response.data.identity.id,
          email: response.data.identity.traits.email,
          name: response.data.identity.traits.name || { first: '', last: '' },
          favorites: response.data.identity.traits.favorites || []
        }
      }
    } else {
      throw new Error('Not authenticated')
    }
  } catch (error) {
    console.error('Error getting profile:', error)
    throw error
  }
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