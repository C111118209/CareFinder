import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

// Create axios instance
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authAPI = {
  async register(email, password, role) {
    return apiClient.post('/auth/register', { email, password, role })
  },
  
  async login(email, password) {
    const data = await apiClient.post('/auth/login', { email, password })
    if (data.token) {
      localStorage.setItem('token', data.token)
      if (data.user) {
        localStorage.setItem('user', JSON.stringify(data.user))
      }
    }
    return data
  },
  
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  },
}

// User API
export const userAPI = {
  async getCurrentUser() {
    return apiClient.get('/users/me')
  },
  
  async getUser(id) {
    return apiClient.get(`/users/${id}`)
  },
  
  async updateUser(id, userData) {
    return apiClient.put(`/users/${id}`, userData)
  },
  
  async updatePassword(id, oldPassword, newPassword) {
    return apiClient.put(`/users/${id}/password`, {
      old_password: oldPassword,
      new_password: newPassword,
    })
  },
}

// Caregiver API
export const caregiverAPI = {
  async createProfile(profileData) {
    return apiClient.post('/caregivers/profile', profileData)
  },
  
  async updateProfile(profileData) {
    return apiClient.put('/caregivers/profile', profileData)
  },
  
  async getMyProfile() {
    return apiClient.get('/caregivers/profile')
  },
  
  async getProfile(id) {
    return apiClient.get(`/caregivers/${id}`)
  },
  
  async search(filters = {}) {
    return apiClient.get('/caregivers/search', { params: filters })
  },
  
  async getAvailability() {
    return apiClient.get('/caregivers/availability')
  },
  
  async updateAvailability(availabilities) {
    return apiClient.put('/caregivers/availability', { availabilities })
  },
  
  async getSpecialAvailability() {
    return apiClient.get('/caregivers/special-availability')
  },
  
  async updateSpecialAvailability(specialAvailabilities) {
    return apiClient.put('/caregivers/special-availability', { special_availabilities: specialAvailabilities })
  },
}

// License API
export const licenseAPI = {
  async uploadLicense(formData) {
    const token = localStorage.getItem('token')
    const baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    return axios.post(`${baseURL}/caregivers/licenses`, formData, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'multipart/form-data',
      },
    }).then(response => response.data)
  },
  
  async getMyLicenses() {
    return apiClient.get('/caregivers/licenses')
  },
  
  async deleteLicense(id) {
    return apiClient.delete(`/caregivers/licenses/${id}`)
  },
  
  getLicenseImageUrl(filename) {
    const baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    return `${baseURL}/caregivers/licenses/image/${filename}`
  },
}

// Contract API
export const contractAPI = {
  async createContract(contractData) {
    return apiClient.post('/contracts', contractData)
  },
  
  async getContracts() {
    return apiClient.get('/contracts')
  },
  
  async getContract(id) {
    return apiClient.get(`/contracts/${id}`)
  },
  
  async acceptContract(id) {
    return apiClient.put(`/contracts/${id}/accept`)
  },
  
  async completeContract(id) {
    return apiClient.put(`/contracts/${id}/complete`)
  },
  
  async renewContract(id, endDate) {
    return apiClient.post(`/contracts/${id}/renew`, { end_date: endDate })
  },
}

// Review API
export const reviewAPI = {
  async createReview(reviewData) {
    return apiClient.post('/reviews', reviewData)
  },
  
  async getCaregiverReviews(caregiverId) {
    return apiClient.get(`/reviews/caregivers/${caregiverId}`)
  },
}

// Admin API
export const adminAPI = {
  async getAllUsers() {
    return apiClient.get('/admin/users')
  },
  
  async toggleUserStatus(userId, isActive) {
    return apiClient.put(`/admin/users/${userId}/status`, { is_active: isActive })
  },
  
  async getAllCaregivers() {
    return apiClient.get('/admin/caregivers')
  },
  
  async updateCaregiverStatus(profileId, isActive) {
    return apiClient.put(`/admin/caregivers/${profileId}/status`, { is_active: isActive })
  },
  
  async getAllContracts() {
    return apiClient.get('/admin/contracts')
  },
  
  async updateContractStatus(contractId, status) {
    return apiClient.put(`/admin/contracts/${contractId}/status`, { status })
  },
  
  async getAllLicenses() {
    return apiClient.get('/admin/licenses')
  },
  
  async getPendingLicenses() {
    return apiClient.get('/admin/licenses/pending')
  },
  
  async reviewLicense(licenseId, status, note = '') {
    return apiClient.put(`/admin/licenses/${licenseId}/review`, { status, note })
  },
  
  async getStatistics() {
    return apiClient.get('/admin/statistics')
  },
}

export default apiClient