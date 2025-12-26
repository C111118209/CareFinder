import { defineStore } from 'pinia'
import { authAPI, userAPI } from '@/services/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    token: localStorage.getItem('token') || null,
    isAuthenticated: !!localStorage.getItem('token'),
  }),
  
  getters: {
    isUser: (state) => state.user?.role === 'user',
    isCaregiver: (state) => state.user?.role === 'caregiver',
    isAdmin: (state) => state.user?.role === 'admin',
  },
  
  actions: {
    async login(email, password) {
      try {
        const data = await authAPI.login(email, password)
        this.user = data.user
        this.token = data.token
        this.isAuthenticated = true
        return data
      } catch (error) {
        throw error
      }
    },
    
    async register(email, password, role) {
      try {
        const data = await authAPI.register(email, password, role)
        return data
      } catch (error) {
        throw error
      }
    },
    
    async fetchCurrentUser() {
      try {
        const user = await userAPI.getCurrentUser()
        this.user = user
        localStorage.setItem('user', JSON.stringify(user))
        return user
      } catch (error) {
        this.logout()
        throw error
      }
    },
    
    logout() {
      authAPI.logout()
      this.user = null
      this.token = null
      this.isAuthenticated = false
    },
    
    updateUser(userData) {
      this.user = { ...this.user, ...userData }
      localStorage.setItem('user', JSON.stringify(this.user))
    },
  },
})