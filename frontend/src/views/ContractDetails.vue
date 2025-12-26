<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="mb-4">
          <router-link
            to="/contracts"
            class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
          >
            ← 返回合約列表
          </router-link>
        </div>

        <h2 class="text-3xl font-bold mb-6 text-gray-800">合約詳情</h2>

        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p class="mt-4 text-gray-600">載入中...</p>
        </div>
        <div v-else-if="contract" class="space-y-6">
          <!-- 合約基本資訊 -->
          <div class="bg-white p-6 rounded-lg shadow-md border border-gray-200">
            <div class="flex justify-between items-start mb-4">
              <div>
                <h3 class="text-2xl font-semibold mb-2">合約 #{{ contract.contract_id }}</h3>
                <span
                  :class="{
                    'bg-yellow-100 text-yellow-800': contract.status === 'pending',
                    'bg-green-100 text-green-800': contract.status === 'active',
                    'bg-blue-100 text-blue-800': contract.status === 'completed',
                    'bg-red-100 text-red-800': contract.status === 'canceled',
                  }"
                  class="px-3 py-1 rounded-full text-sm font-medium"
                >
                  {{ getStatusText(contract.status) }}
                </span>
                <span v-if="contract.is_renewal" class="ml-2 text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                  續約
                </span>
              </div>
            </div>

            <div class="grid md:grid-cols-2 gap-4 text-gray-700">
              <div>
                <p><span class="font-medium">開始日期：</span>{{ formatDate(contract.start_date) }}</p>
                <p><span class="font-medium">結束日期：</span>{{ formatDate(contract.end_date) }}</p>
              </div>
              <div>
                <p><span class="font-medium">合約期間：</span>{{ calculateDuration(contract.start_date, contract.end_date) }} 天</p>
                <p v-if="contract.created_at">
                  <span class="font-medium">建立時間：</span>{{ formatDate(contract.created_at) }}
                </p>
              </div>
            </div>
          </div>

          <!-- 照護者詳細資料 -->
          <div class="bg-white p-6 rounded-lg shadow-md border-l-4 border-blue-500">
            <h3 class="text-xl font-semibold mb-4 text-blue-800 flex items-center">
              <span class="mr-2">👤</span> 照護者詳細資料
            </h3>
            <div class="grid md:grid-cols-2 gap-6">
              <div>
                <p><span class="font-medium">暱稱：</span>{{ getCaregiverName(contract) }}</p>
                <p><span class="font-medium">電子郵件：</span>{{ getCaregiverEmail(contract) }}</p>
              </div>
              <div v-if="extraProfile">
                <p v-if="extraProfile.phone">
                  <span class="font-medium">電話：</span>{{ extraProfile.phone }}
                </p>
                <p v-if="extraProfile.address">
                  <span class="font-medium">地址：</span>{{ extraProfile.address }}
                </p>
                <p v-if="extraProfile.service_rate">
                  <span class="font-medium">時薪：</span>NT$ {{ extraProfile.service_rate }}/小時
                </p>
              </div>
            </div>
            <div v-if="extraProfile && extraProfile.bio" class="mt-4 pt-4 border-t">
              <p class="font-medium mb-2">簡介：</p>
              <p class="text-gray-700">{{ extraProfile.bio }}</p>
            </div>
          </div>

          <!-- 家屬詳細資料 -->
          <div class="bg-white p-6 rounded-lg shadow-md border-l-4 border-green-500">
            <h3 class="text-xl font-semibold mb-4 text-green-800 flex items-center">
              <span class="mr-2">👨‍👩‍👧‍👦</span> 家屬詳細資料
            </h3>
            <div class="grid md:grid-cols-2 gap-6">
              <div>
                <p><span class="font-medium">暱稱：</span>{{ getUserName(contract) }}</p>
                <p><span class="font-medium">電子郵件：</span>{{ getUserEmail(contract) }}</p>
              </div>
            </div>
          </div>

          <!-- 操作按鈕 -->
          <div v-if="contract.status === 'pending' && authStore.isCaregiver" class="bg-white p-6 rounded-lg shadow-md">
            <button
              @click="handleAccept"
              class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium"
            >
              接受合約
            </button>
          </div>

          <div v-if="contract.status === 'active' && authStore.isCaregiver" class="bg-white p-6 rounded-lg shadow-md">
            <button
              @click="handleComplete"
              class="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-lg font-medium"
            >
              完成合約
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { contractAPI, caregiverAPI } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const contract = ref(null)
const extraProfile = ref(null)
const loading = ref(false)

const formatDate = (dateString) => {
  if (!dateString) return '無資料'
  const date = new Date(dateString)
  return isNaN(date.getTime())
    ? '格式錯誤'
    : date.toLocaleDateString('zh-TW', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
      })
}

const getStatusText = (status) => {
  const statusMap = {
    pending: '待接受',
    active: '進行中',
    completed: '已完成',
    canceled: '已取消',
  }
  return statusMap[status] || status
}

const getCaregiverName = (contract) => {
  const caregiver = contract.Caregiver || contract.caregiver || {}
  return caregiver.nickname || caregiver.email || '未提供'
}

const getCaregiverEmail = (contract) => {
  const caregiver = contract.Caregiver || contract.caregiver || {}
  return caregiver.email || '未提供'
}

const getUserName = (contract) => {
  const user = contract.User || contract.user || {}
  return user.nickname || user.email || '未提供'
}

const getUserEmail = (contract) => {
  const user = contract.User || contract.user || {}
  return user.email || '未提供'
}

const calculateDuration = (startDate, endDate) => {
  const start = new Date(startDate)
  const end = new Date(endDate)
  const diffTime = Math.abs(end - start)
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24)) || 0
}

const handleAccept = async () => {
  if (!confirm('確定要接受這個合約嗎？')) return
  try {
    await contractAPI.acceptContract(contract.value.contract_id)
    contract.value.status = 'active'
    alert('合約已接受')
    router.push('/contracts')
  } catch (error) {
    alert('操作失敗：' + (error.response?.data?.error || error.message))
  }
}

const handleComplete = async () => {
  if (!confirm('確定要完成這個合約嗎？')) return
  try {
    await contractAPI.completeContract(contract.value.contract_id)
    contract.value.status = 'completed'
    alert('合約已完成')
    router.push('/contracts')
  } catch (error) {
    alert('操作失敗：' + (error.response?.data?.error || error.message))
  }
}

const loadContractDetails = async () => {
  loading.value = true
  try {
    const contractId = route.params.id
    contract.value = await contractAPI.getContract(contractId)

    // 嘗試取得照護者的詳細 Profile
    if (contract.value.caregiver_id) {
      try {
        extraProfile.value = await caregiverAPI.getProfile(contract.value.caregiver_id)
      } catch (e) {
        console.log('無法取得額外 Profile 資訊')
      }
    }
  } catch (error) {
    console.error('Load contract error:', error)
    alert('載入失敗：' + (error.response?.data?.error || error.message))
    router.push('/contracts')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (!route.params.id) {
    router.push('/contracts')
    return
  }
  loadContractDetails()
})
</script>