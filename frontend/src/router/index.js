import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/Home.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { guest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/Register.vue'),
      meta: { guest: true },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/Dashboard.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/Profile.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/caregiver-search',
      name: 'caregiver-search',
      component: () => import('@/views/CaregiverSearch.vue'),
      meta: { requiresAuth: true, requiresUser: true },
    },
    {
      path: '/caregiver/:id',
      name: 'caregiver-profile',
      component: () => import('@/views/CaregiverProfile.vue'),
    },
    {
      path: '/caregiver-profile-setup',
      name: 'caregiver-profile-setup',
      component: () => import('@/views/CaregiverProfileSetup.vue'),
      meta: { requiresAuth: true, requiresCaregiver: true },
    },
    {
      path: '/caregiver-availability',
      name: 'caregiver-availability',
      component: () => import('@/views/CaregiverAvailability.vue'),
      meta: { requiresAuth: true, requiresCaregiver: true },
    },
    {
      path: '/caregiver-licenses',
      name: 'caregiver-licenses',
      component: () => import('@/views/CaregiverLicenses.vue'),
      meta: { requiresAuth: true, requiresCaregiver: true },
    },
    {
      path: '/contracts',
      name: 'contracts',
      component: () => import('@/views/Contracts.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/contract/:id',
      name: 'contract-details',
      component: () => import('@/views/ContractDetails.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/create-contract/:caregiverId',
      name: 'create-contract',
      component: () => import('@/views/CreateContract.vue'),
      meta: { requiresAuth: true, requiresUser: true },
    },
    {
      path: '/admin',
      name: 'admin-dashboard',
      component: () => import('@/views/AdminDashboard.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
  ],
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }
  
  if (to.meta.guest && authStore.isAuthenticated) {
    next({ name: 'dashboard' })
    return
  }
  
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'dashboard' })
    return
  }
  
  if (to.meta.requiresUser && !authStore.isUser) {
    next({ name: 'dashboard' })
    return
  }
  
  if (to.meta.requiresCaregiver && !authStore.isCaregiver) {
    next({ name: 'dashboard' })
    return
  }
  
  next()
})

export default router