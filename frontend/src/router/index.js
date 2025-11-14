import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/models'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { showNav: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { showNav: false }
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/views/ResetPassword.vue'),
    meta: { showNav: false }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { showNav: true, requiresAuth: true }
  },
  {
    path: '/models',
    name: 'Models',
    component: () => import('@/views/Models.vue'),
    meta: { showNav: true, requiresAuth: true }
  },
  {
    path: '/models/:id',
    name: 'ModelDetail',
    component: () => import('@/views/ModelDetail.vue'),
    meta: { showNav: true, requiresAuth: true }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: () => import('@/views/Orders.vue'),
    meta: { showNav: true, requiresAuth: true }
  },
  {
    path: '/pay/:orderId',
    name: 'Pay',
    component: () => import('@/views/Pay.vue'),
    meta: { showNav: true, requiresAuth: true }
  },
  {
    path: '/library',
    name: 'Library',
    component: () => import('@/views/Library.vue'),
    meta: { showNav: true, requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/views/Settings.vue'),
    meta: { showNav: false }
  }
  ,
  {
    path: '/error',
    name: 'Error',
    component: () => import('@/views/Error.vue'),
    meta: { showNav: false }
  }
]

export default routes
