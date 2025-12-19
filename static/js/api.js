// 替換為 Render 網址
const API_BASE_URL = './api/v1'; 

// Get token from localStorage
function getToken() {
  return localStorage.getItem('token');
}

// Set token to localStorage
function setToken(token) {
  localStorage.setItem('token', token);
}

// Remove token from localStorage
function removeToken() {
  localStorage.removeItem('token');
}

// Get user info from localStorage
function getUser() {
  const userStr = localStorage.getItem('user');
  return userStr ? JSON.parse(userStr) : null;
}

// Set user info to localStorage
function setUser(user) {
  localStorage.setItem('user', JSON.stringify(user));
}

// Remove user from localStorage
function removeUser() {
  localStorage.removeItem('user');
}

// API request helper
async function apiRequest(endpoint, options = {}) {
  const token = getToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Request failed');
    }

    return data;
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
}

// Auth API
const authAPI = {
  async register(email, password, role) {
    return apiRequest('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, role }),
    });
  },

  async login(email, password) {
    const data = await apiRequest('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    if (data.token) {
      setToken(data.token);
    }
    // Return full response including user info
    return data;
  },

  logout() {
    removeToken();
    removeUser();
    window.location.href = '/login.html';
  },
};

// User API
const userAPI = {
  async getCurrentUser() {
    return apiRequest('/users/me');
  },

  async getUser(id) {
    return apiRequest(`/users/${id}`);
  },

  async updateUser(id, userData) {
    return apiRequest(`/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(userData),
    });
  },

  async updatePassword(id, oldPassword, newPassword) {
    return apiRequest(`/users/${id}/password`, {
      method: 'PUT',
      body: JSON.stringify({
        old_password: oldPassword,
        new_password: newPassword,
      }),
    });
  },
};

// Caregiver API
const caregiverAPI = {
  async createProfile(profileData) {
    return apiRequest('/caregivers/profile', {
      method: 'POST',
      body: JSON.stringify(profileData),
    });
  },

  async updateProfile(profileData) {
    return apiRequest('/caregivers/profile', {
      method: 'PUT',
      body: JSON.stringify(profileData),
    });
  },

  async getMyProfile() {
    return apiRequest('/caregivers/profile');
  },

  async getProfile(id) {
    return apiRequest(`/caregivers/${id}`);
  },

  async search(filters = {}) {
    const params = new URLSearchParams(filters);
    return apiRequest(`/caregivers/search?${params}`);
  },

  async getAvailability() {
    return apiRequest('/caregivers/availability');
  },

  async updateAvailability(availabilities) {
    return apiRequest('/caregivers/availability', {
      method: 'PUT',
      body: JSON.stringify({ availabilities }),
    });
  },

  async getSpecialAvailability() {
    return apiRequest('/caregivers/special-availability');
  },

  async updateSpecialAvailability(specialAvailabilities) {
    return apiRequest('/caregivers/special-availability', {
      method: 'PUT',
      body: JSON.stringify({ special_availabilities: specialAvailabilities }),
    });
  },
};

// License API
const licenseAPI = {
  async uploadLicense(formData) {
    const token = getToken();
    if (!token) {
      throw new Error('Not authenticated');
    }

    try {
      const response = await fetch(`${API_BASE_URL}/caregivers/licenses`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Upload failed');
      }

      return data;
    } catch (error) {
      console.error('License API Error:', error);
      throw error;
    }
  },

  async getMyLicenses() {
    return apiRequest('/caregivers/licenses');
  },

  async deleteLicense(id) {
    return apiRequest(`/caregivers/licenses/${id}`, {
      method: 'DELETE',
    });
  },
};

// Contract API
const contractAPI = {
  async createContract(contractData) {
    return apiRequest('/contracts', {
      method: 'POST',
      body: JSON.stringify(contractData),
    });
  },

  async getContracts() {
    return apiRequest('/contracts');
  },

  async getContract(id) {
    return apiRequest(`/contracts/${id}`);
  },

  async acceptContract(id) {
    return apiRequest(`/contracts/${id}/accept`, {
      method: 'PUT',
    });
  },

  async completeContract(id) {
    return apiRequest(`/contracts/${id}/complete`, {
      method: 'PUT',
    });
  },

  async renewContract(id, endDate) {
    return apiRequest(`/contracts/${id}/renew`, {
      method: 'POST',
      body: JSON.stringify({ end_date: endDate }),
    });
  },
};

// Review API
const reviewAPI = {
  async createReview(reviewData) {
    return apiRequest('/reviews', {
      method: 'POST',
      body: JSON.stringify(reviewData),
    });
  },

  async getCaregiverReviews(caregiverId) {
    return apiRequest(`/reviews/caregivers/${caregiverId}`);
  },
};

// Check if user is authenticated
function isAuthenticated() {
  return !!getToken();
}

// Redirect to login if not authenticated
function requireAuth() {
  if (!isAuthenticated()) {
    window.location.href = '/login.html';
    return false;
  }
  return true;
}

// Check if user is caregiver, redirect if not
function requireCaregiverRole() {
  if (!requireAuth()) return false;
  
  const user = getUser();
  if (!user || user.role !== 'caregiver') {
    showNotification('只有照護者可以訪問此頁面', 'error');
    setTimeout(() => {
      window.location.href = '/dashboard.html';
    }, 2000);
    return false;
  }
  return true;
}

// Check if user is regular user, redirect if not
function requireUserRole() {
  if (!requireAuth()) return false;
  
  const user = getUser();
  if (!user || user.role !== 'user') {
    showNotification('只有使用者可以訪問此頁面', 'error');
    setTimeout(() => {
      window.location.href = '/dashboard.html';
    }, 2000);
    return false;
  }
  return true;
}

// Check if user is admin, redirect if not
function requireAdminRole() {
  if (!requireAuth()) return false;
  
  const user = getUser();
  if (!user || user.role !== 'admin') {
    showNotification('只有管理員可以訪問此頁面', 'error');
    setTimeout(() => {
      window.location.href = '/dashboard.html';
    }, 2000);
    return false;
  }
  return true;
}

// Admin API
const adminAPI = {
  async getAllUsers() {
    return apiRequest('/admin/users');
  },

  async toggleUserStatus(userId, isActive) {
    return apiRequest(`/admin/users/${userId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ is_active: isActive }),
    });
  },

  async getAllCaregivers() {
    return apiRequest('/admin/caregivers');
  },

  async updateCaregiverStatus(profileId, isActive) {
    return apiRequest(`/admin/caregivers/${profileId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ is_active: isActive }),
    });
  },

  async getAllContracts() {
    return apiRequest('/admin/contracts');
  },

  async updateContractStatus(contractId, status) {
    return apiRequest(`/admin/contracts/${contractId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
    });
  },

  async getAllLicenses() {
    return apiRequest('/admin/licenses');
  },

  async getPendingLicenses() {
    return apiRequest('/admin/licenses/pending');
  },

  async reviewLicense(licenseId, status, note = '') {
    return apiRequest(`/admin/licenses/${licenseId}/review`, {
      method: 'PUT',
      body: JSON.stringify({ status, note }),
    });
  },

  async getStatistics() {
    return apiRequest('/admin/statistics');
  },
};

// Format date
function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

// Format time
function formatTime(timeString) {
  return timeString;
}

// Show notification
function showNotification(message, type = 'info') {
  const notification = document.createElement('div');
  notification.className = `fixed top-4 right-4 p-4 rounded-lg shadow-lg z-50 ${
    type === 'error' ? 'bg-red-500' : type === 'success' ? 'bg-green-500' : 'bg-blue-500'
  } text-white`;
  notification.textContent = message;
  document.body.appendChild(notification);

  setTimeout(() => {
    notification.remove();
  }, 3000);
}

