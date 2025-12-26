<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <!-- Navigation Tabs -->
        <div class="mb-6 flex space-x-2 border-b border-gray-200">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'px-4 py-2 font-medium transition-colors',
              activeTab === tab.id
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-700 hover:text-gray-900'
            ]"
          >
            {{ tab.label }}
          </button>
        </div>

        <!-- Statistics Section -->
        <div v-if="activeTab === 'statistics'" class="space-y-6">
          <h2 class="text-3xl font-bold text-gray-800">系統統計報告</h2>
          <div v-if="loadingStats" class="text-center py-8">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p class="mt-4 text-gray-600">載入中...</p>
          </div>
          <div v-else-if="statistics" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div class="card bg-blue-50">
              <h3 class="text-lg font-semibold text-blue-800 mb-2">用戶統計</h3>
              <div class="text-3xl font-bold text-blue-600 mb-1">{{ statistics.users?.total || 0 }}</div>
              <div class="text-sm text-gray-600">
                <p>啟用：{{ statistics.users?.active || 0 }} | 停用：{{ statistics.users?.inactive || 0 }}</p>
                <p class="mt-2">
                  使用者：{{ statistics.users?.by_role?.user || 0 }} | 
                  照護者：{{ statistics.users?.by_role?.caregiver || 0 }} | 
                  管理員：{{ statistics.users?.by_role?.admin || 0 }}
                </p>
              </div>
            </div>
            <div class="card bg-green-50">
              <h3 class="text-lg font-semibold text-green-800 mb-2">照護者統計</h3>
              <div class="text-3xl font-bold text-green-600 mb-1">{{ statistics.caregivers?.total || 0 }}</div>
              <div class="text-sm text-gray-600">啟用：{{ statistics.caregivers?.active || 0 }}</div>
            </div>
            <div class="card bg-purple-50">
              <h3 class="text-lg font-semibold text-purple-800 mb-2">合約統計</h3>
              <div class="text-3xl font-bold text-purple-600 mb-1">{{ statistics.contracts?.total || 0 }}</div>
              <div class="text-sm text-gray-600">
                <p>待處理：{{ statistics.contracts?.pending || 0 }} | 進行中：{{ statistics.contracts?.active || 0 }}</p>
                <p>已完成：{{ statistics.contracts?.completed || 0 }} | 已取消：{{ statistics.contracts?.canceled || 0 }}</p>
              </div>
            </div>
            <div class="card bg-yellow-50">
              <h3 class="text-lg font-semibold text-yellow-800 mb-2">證照統計</h3>
              <div class="text-3xl font-bold text-yellow-600 mb-1">{{ statistics.licenses?.total || 0 }}</div>
              <div class="text-sm text-gray-600">
                <p>待審核：{{ statistics.licenses?.pending || 0 }} | 已通過：{{ statistics.licenses?.approved || 0 }}</p>
                <p>已拒絕：{{ statistics.licenses?.rejected || 0 }}</p>
              </div>
            </div>
            <div class="card bg-indigo-50">
              <h3 class="text-lg font-semibold text-indigo-800 mb-2">評價統計</h3>
              <div class="text-3xl font-bold text-indigo-600 mb-1">{{ statistics.reviews?.total || 0 }}</div>
              <div class="text-sm text-gray-600">
                平均評分：{{ statistics.reviews?.avg_rating ? statistics.reviews.avg_rating.toFixed(1) : '0.0' }}
              </div>
            </div>
          </div>
        </div>

        <!-- Users Management -->
        <div v-if="activeTab === 'users'" class="space-y-6">
          <div class="flex justify-between items-center">
            <h2 class="text-3xl font-bold text-gray-800">用戶管理</h2>
            <button @click="loadUsers" class="btn-secondary">重新載入</button>
          </div>
          <div v-if="loadingUsers" class="text-center py-8">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p class="mt-4 text-gray-600">載入中...</p>
          </div>
          <div v-else class="card overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">電子郵件</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">暱稱</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">角色</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">狀態</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="user in users" :key="user.user_id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.user_id }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.email }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.nickname || '-' }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      :class="{
                        'bg-purple-100 text-purple-800': user.role === 'admin',
                        'bg-green-100 text-green-800': user.role === 'caregiver',
                        'bg-blue-100 text-blue-800': user.role === 'user',
                      }"
                      class="px-2 py-1 text-xs rounded-full"
                    >
                      {{ getRoleText(user.role) }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      :class="user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                      class="px-2 py-1 text-xs rounded-full"
                    >
                      {{ user.is_active ? '啟用' : '停用' }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <button
                      v-if="user.role !== 'admin'"
                      @click="toggleUserStatus(user.user_id, !user.is_active)"
                      :class="user.is_active ? 'bg-red-500 hover:bg-red-600' : 'bg-green-500 hover:bg-green-600'"
                      class="btn-secondary text-sm text-white"
                    >
                      {{ user.is_active ? '停用' : '啟用' }}
                    </button>
                    <span v-else class="text-gray-400">-</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Caregivers Management -->
        <div v-if="activeTab === 'caregivers'" class="space-y-6">
          <div class="flex justify-between items-center">
            <h2 class="text-3xl font-bold text-gray-800">照護者管理</h2>
            <button @click="loadCaregivers" class="btn-secondary">重新載入</button>
          </div>
          <div v-if="loadingCaregivers" class="text-center py-8">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p class="mt-4 text-gray-600">載入中...</p>
          </div>
          <div v-else class="card overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">姓名</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">電子郵件</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">電話</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">評分</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">狀態</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="caregiver in caregivers" :key="caregiver.profile_id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ caregiver.profile_id }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ caregiver.full_name }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ caregiver.user?.email || '-' }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ caregiver.phone || '-' }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ caregiver.avg_rating ? caregiver.avg_rating.toFixed(1) : '0.0' }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      :class="caregiver.user?.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                      class="px-2 py-1 text-xs rounded-full"
                    >
                      {{ caregiver.user?.is_active ? '啟用' : '停用' }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <button
                      @click="toggleCaregiverStatus(caregiver.profile_id, !caregiver.user?.is_active)"
                      :class="caregiver.user?.is_active ? 'bg-red-500 hover:bg-red-600' : 'bg-green-500 hover:bg-green-600'"
                      class="btn-secondary text-sm text-white"
                    >
                      {{ caregiver.user?.is_active ? '停用' : '啟用' }}
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Contracts Management -->
        <div v-if="activeTab === 'contracts'" class="space-y-6">
          <div class="flex justify-between items-center">
            <h2 class="text-3xl font-bold text-gray-800">合約管理</h2>
            <button @click="loadContracts" class="btn-secondary">重新載入</button>
          </div>
          <div v-if="loadingContracts" class="text-center py-8">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p class="mt-4 text-gray-600">載入中...</p>
          </div>
          <div v-else class="card overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">使用者</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">照護者</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">開始日期</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">結束日期</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">狀態</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="contract in contracts" :key="contract.contract_id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ contract.contract_id }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ contract.user?.email || '-' }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ contract.caregiver?.email || '-' }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ formatDate(contract.start_date) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ formatDate(contract.end_date) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      :class="getStatusColor(contract.status)"
                      class="px-2 py-1 text-xs rounded-full"
                    >
                      {{ getStatusText(contract.status) }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <select
                      @change="updateContractStatus(contract.contract_id, $event.target.value)"
                      class="text-sm border rounded px-2 py-1"
                    >
                      <option value="pending" :selected="contract.status === 'pending'">待處理</option>
                      <option value="active" :selected="contract.status === 'active'">進行中</option>
                      <option value="completed" :selected="contract.status === 'completed'">已完成</option>
                      <option value="canceled" :selected="contract.status === 'canceled'">已取消</option>
                    </select>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Licenses Review -->
        <div v-if="activeTab === 'licenses'" class="space-y-6">
          <div class="flex justify-between items-center">
            <h2 class="text-3xl font-bold text-gray-800">證照審核</h2>
            <div class="flex space-x-2">
              <button @click="loadLicenses('all')" class="btn-secondary">所有證照</button>
              <button @click="loadLicenses('pending')" class="btn-primary">待審核</button>
            </div>
          </div>
          <div v-if="loadingLicenses" class="text-center py-8">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p class="mt-4 text-gray-600">載入中...</p>
          </div>
          <div v-else-if="licenses.length === 0" class="text-center py-8 text-gray-600">
            目前沒有證照
          </div>
          <div v-else class="space-y-4">
            <div
              v-for="license in licenses"
              :key="license.license_id"
              class="card"
            >
              <div class="flex justify-between items-start">
                <div class="flex-1">
                  <div class="flex items-center space-x-3 mb-2">
                    <h3 class="text-lg font-semibold text-gray-800">{{ license.name }}</h3>
                    <span
                      :class="getLicenseStatusColor(license.status)"
                      class="px-3 py-1 rounded-full text-sm font-medium"
                    >
                      {{ getLicenseStatusText(license.status) }}
                    </span>
                  </div>
                  <div class="text-gray-600 space-y-1 text-sm">
                    <p>
                      <span class="font-medium">照護者：</span>
                      {{ license.caregiver_profile?.full_name || '-' }}
                      ({{ license.caregiver_profile?.user?.email || '-' }})
                    </p>
                    <p><span class="font-medium">發證日期：</span>{{ formatDate(license.issue_date) }}</p>
                    <p><span class="font-medium">有效期限：</span>{{ formatDate(license.expiry_date) }}</p>
                    <p v-if="license.proof_url">
                      <span class="font-medium">證明文件：</span>
                      <button
                        @click="viewLicenseImage(license.proof_url)"
                        class="text-primary-600 hover:underline"
                      >
                        查看
                      </button>
                    </p>
                  </div>
                </div>
                <div v-if="license.status === 'pending'" class="flex space-x-2 ml-4">
                  <button
                    @click="reviewLicense(license.license_id, 'approved')"
                    class="btn-primary bg-green-500 hover:bg-green-600 text-white"
                  >
                    通過
                  </button>
                  <button
                    @click="reviewLicense(license.license_id, 'rejected')"
                    class="btn-secondary bg-red-500 hover:bg-red-600 text-white"
                  >
                    拒絕
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminAPI, licenseAPI } from '@/services/api'
import Navbar from '@/components/Navbar.vue'

const tabs = [
  { id: 'statistics', label: '統計報告' },
  { id: 'users', label: '用戶管理' },
  { id: 'caregivers', label: '照護者管理' },
  { id: 'contracts', label: '合約管理' },
  { id: 'licenses', label: '證照審核' },
]

const activeTab = ref('statistics')
const statistics = ref(null)
const users = ref([])
const caregivers = ref([])
const contracts = ref([])
const licenses = ref([])
const loadingStats = ref(false)
const loadingUsers = ref(false)
const loadingCaregivers = ref(false)
const loadingContracts = ref(false)
const loadingLicenses = ref(false)

const formatDate = (dateString) => {
  if (!dateString) return '無資料'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

const getRoleText = (role) => {
  const roleMap = {
    admin: '管理員',
    caregiver: '照護者',
    user: '使用者',
  }
  return roleMap[role] || role
}

const getStatusText = (status) => {
  const statusMap = {
    pending: '待處理',
    active: '進行中',
    completed: '已完成',
    canceled: '已取消',
  }
  return statusMap[status] || status
}

const getStatusColor = (status) => {
  const colorMap = {
    pending: 'bg-yellow-100 text-yellow-800',
    active: 'bg-green-100 text-green-800',
    completed: 'bg-blue-100 text-blue-800',
    canceled: 'bg-red-100 text-red-800',
  }
  return colorMap[status] || 'bg-gray-100 text-gray-800'
}

const getLicenseStatusText = (status) => {
  const statusMap = {
    pending: '待審核',
    approved: '已通過',
    rejected: '已拒絕',
  }
  return statusMap[status] || status
}

const getLicenseStatusColor = (status) => {
  const colorMap = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
  }
  return colorMap[status] || 'bg-gray-100 text-gray-800'
}

const loadStatistics = async () => {
  loadingStats.value = true
  try {
    statistics.value = await adminAPI.getStatistics()
  } catch (error) {
    console.error('Load statistics error:', error)
    alert('載入統計資料失敗')
  } finally {
    loadingStats.value = false
  }
}

const loadUsers = async () => {
  loadingUsers.value = true
  try {
    users.value = await adminAPI.getAllUsers()
  } catch (error) {
    console.error('Load users error:', error)
    alert('載入用戶列表失敗')
  } finally {
    loadingUsers.value = false
  }
}

const loadCaregivers = async () => {
  loadingCaregivers.value = true
  try {
    caregivers.value = await adminAPI.getAllCaregivers()
  } catch (error) {
    console.error('Load caregivers error:', error)
    alert('載入照護者列表失敗')
  } finally {
    loadingCaregivers.value = false
  }
}

const loadContracts = async () => {
  loadingContracts.value = true
  try {
    contracts.value = await adminAPI.getAllContracts()
  } catch (error) {
    console.error('Load contracts error:', error)
    alert('載入合約列表失敗')
  } finally {
    loadingContracts.value = false
  }
}

const loadLicenses = async (type = 'pending') => {
  loadingLicenses.value = true
  try {
    licenses.value = type === 'pending'
      ? await adminAPI.getPendingLicenses()
      : await adminAPI.getAllLicenses()
  } catch (error) {
    console.error('Load licenses error:', error)
    alert('載入證照列表失敗')
  } finally {
    loadingLicenses.value = false
  }
}

const toggleUserStatus = async (userId, isActive) => {
  if (!confirm(`確定要${isActive ? '啟用' : '停用'}此用戶嗎？`)) return
  try {
    await adminAPI.toggleUserStatus(userId, isActive)
    alert('用戶狀態已更新')
    loadUsers()
  } catch (error) {
    alert('更新失敗：' + (error.response?.data?.error || error.message))
  }
}

const toggleCaregiverStatus = async (profileId, isActive) => {
  if (!confirm(`確定要${isActive ? '啟用' : '停用'}此照護者嗎？`)) return
  try {
    await adminAPI.updateCaregiverStatus(profileId, isActive)
    alert('照護者狀態已更新')
    loadCaregivers()
  } catch (error) {
    alert('更新失敗：' + (error.response?.data?.error || error.message))
  }
}

const updateContractStatus = async (contractId, status) => {
  try {
    await adminAPI.updateContractStatus(contractId, status)
    alert('合約狀態已更新')
    loadContracts()
  } catch (error) {
    alert('更新失敗：' + (error.response?.data?.error || error.message))
  }
}

const reviewLicense = async (licenseId, status) => {
  const action = status === 'approved' ? '通過' : '拒絕'
  if (!confirm(`確定要${action}此證照嗎？`)) return
  try {
    await adminAPI.reviewLicense(licenseId, status)
    alert(`證照已${action}`)
    loadLicenses('pending')
  } catch (error) {
    alert('審核失敗：' + (error.response?.data?.error || error.message))
  }
}

const viewLicenseImage = async (proofUrl) => {
  if (!proofUrl) return
  const token = localStorage.getItem('token')
  if (!token) {
    alert('請先登入')
    return
  }

  try {
    let imageUrl = proofUrl
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    
    if (proofUrl.startsWith('/uploads/licenses/')) {
      const filename = proofUrl.replace('/uploads/licenses/', '')
      imageUrl = `${API_BASE_URL}/caregivers/licenses/image/${filename}`
    } else if (!proofUrl.startsWith('http')) {
      if (proofUrl.startsWith('/api/v1/')) {
        imageUrl = `http://localhost:8080${proofUrl}`
      } else {
        imageUrl = `${API_BASE_URL}${proofUrl}`
      }
    }

    // 使用 fetch 並添加 Authorization header
    const response = await fetch(imageUrl, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      throw new Error('無法載入圖片')
    }

    const blob = await response.blob()
    const blobUrl = URL.createObjectURL(blob)
    const newWindow = window.open()
    
    if (newWindow) {
      newWindow.document.write(`
        <html>
          <head><title>證照預覽</title></head>
          <body style="margin:0;padding:0;background:#222;display:flex;justify-content:center;align-items:center;height:100vh;">
            <img src="${blobUrl}" style="max-width:95%;max-height:95vh;box-shadow:0 0 20px rgba(0,0,0,0.5);" />
          </body>
        </html>
      `)
    }
  } catch (error) {
    console.error('View license image error:', error)
    alert('無法載入圖片：' + (error.message || '未知錯誤'))
  }
}

onMounted(() => {
  loadStatistics()
})
</script>