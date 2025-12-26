<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h2 class="text-3xl font-bold mb-6 text-gray-800">我的合約</h2>

        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p class="mt-4 text-gray-600">載入中...</p>
        </div>
        <div v-else-if="contracts.length === 0" class="text-center py-8 text-gray-600">
          目前沒有合約
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="contract in contracts"
            :key="contract.contract_id"
            class="bg-white p-6 rounded-lg shadow-md border border-gray-200"
          >
            <div class="flex flex-col md:flex-row justify-between items-start gap-4">
              <div class="flex-1 w-full">
                <div class="flex items-center space-x-4 mb-4">
                  <h3 class="text-xl font-semibold text-gray-800">合約 #{{ contract.contract_id }}</h3>
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
                  <span v-if="contract.is_renewal" class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                    續約
                  </span>
                </div>

                <div class="grid md:grid-cols-2 gap-4 text-gray-700 mb-4">
                  <div>
                    <p><span class="font-medium">開始日期：</span>{{ formatDate(contract.start_date) }}</p>
                    <p><span class="font-medium">結束日期：</span>{{ formatDate(contract.end_date) }}</p>
                  </div>
                  <div>
                    <p><span class="font-medium">合約期間：</span>{{ calculateDuration(contract.start_date, contract.end_date) }} 天</p>
                  </div>
                </div>

                <div class="border-t pt-4 mt-4">
                  <div class="grid md:grid-cols-2 gap-6">
                    <div class="bg-blue-50 p-4 rounded-lg">
                      <h4 class="font-semibold text-blue-800 mb-2">照護者資訊</h4>
                      <p class="text-sm">
                        <span class="font-medium">姓名：</span>{{ getCaregiverName(contract) }}
                      </p>
                      <p class="text-sm">
                        <span class="font-medium">電子郵件：</span>{{ getCaregiverEmail(contract) }}
                      </p>
                      <button
                        v-if="contract.caregiver_id"
                        @click="viewDetails(contract.contract_id, 'caregiver')"
                        class="mt-2 text-sm text-blue-600 hover:text-blue-800 underline"
                      >
                        查看詳細資料
                      </button>
                    </div>

                    <div class="bg-green-50 p-4 rounded-lg">
                      <h4 class="font-semibold text-green-800 mb-2">家屬資訊</h4>
                      <p class="text-sm">
                        <span class="font-medium">姓名：</span>{{ getUserName(contract) }}
                      </p>
                      <p class="text-sm">
                        <span class="font-medium">電子郵件：</span>{{ getUserEmail(contract) }}
                      </p>
                      <button
                        @click="viewDetails(contract.contract_id, 'user')"
                        class="mt-2 text-sm text-green-600 hover:text-green-800 underline"
                      >
                        查看詳細資料
                      </button>
                    </div>
                  </div>
                </div>

                <div class="mt-4">
                  <button
                    @click="viewContractDetails(contract.contract_id)"
                    class="text-sm text-blue-600 hover:text-blue-800 underline font-medium"
                  >
                    查看完整合約詳情 &rarr;
                  </button>
                </div>
              </div>

              <div class="flex flex-row md:flex-col gap-2 mt-2 md:mt-0 w-full md:w-auto">
                <button
                  v-if="contract.status === 'pending' && authStore.isCaregiver"
                  @click="acceptContract(contract.contract_id)"
                  class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition whitespace-nowrap"
                >
                  接受合約
                </button>
                <button
                  v-if="contract.status === 'active'"
                  @click="completeContract(contract.contract_id)"
                  class="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 transition whitespace-nowrap"
                >
                  完成服務
                </button>
                <button
                  v-if="contract.status === 'completed' && authStore.isUser"
                  @click="showReviewModal(contract.contract_id, contract.caregiver_id)"
                  class="bg-yellow-500 text-white px-4 py-2 rounded hover:bg-yellow-600 transition whitespace-nowrap"
                >
                  評價
                </button>
                <button
                  v-if="contract.status === 'active' && authStore.isUser && canRenew(contract.end_date)"
                  @click="renewContract(contract.contract_id)"
                  class="bg-gray-200 text-gray-800 px-4 py-2 rounded hover:bg-gray-300 transition whitespace-nowrap"
                >
                  續約
                </button>
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
import { useRouter } from 'vue-router'
import { contractAPI, reviewAPI } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'

const router = useRouter()
const authStore = useAuthStore()

const contracts = ref([])
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
    pending: '待確認',
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

const canRenew = (endDate) => {
  const end = new Date(endDate)
  const now = new Date()
  const daysUntilEnd = Math.ceil((end - now) / (1000 * 60 * 60 * 24))
  return daysUntilEnd <= 30 && daysUntilEnd >= 0
}

const viewContractDetails = (contractId) => {
  router.push(`/contract/${contractId}`)
}

const viewDetails = (contractId, type) => {
  router.push(`/contract/${contractId}?view=${type}`)
}

const acceptContract = async (contractId) => {
  if (!confirm('確定要接受這個合約嗎？')) return
  try {
    await contractAPI.acceptContract(contractId)
    alert('合約已接受！')
    loadContracts()
  } catch (error) {
    alert('操作失敗：' + (error.response?.data?.error || error.message))
  }
}

const completeContract = async (contractId) => {
  if (!confirm('確定要標記服務為已完成嗎？')) return
  try {
    await contractAPI.completeContract(contractId)
    alert('服務已標記為完成！')
    loadContracts()
  } catch (error) {
    alert('操作失敗：' + (error.response?.data?.error || error.message))
  }
}

const renewContract = async (contractId) => {
  const endDate = prompt('請輸入續約結束日期 (YYYY-MM-DD):')
  if (!endDate) return
  try {
    await contractAPI.renewContract(contractId, endDate)
    alert('續約已發起！')
    loadContracts()
  } catch (error) {
    alert('續約失敗：' + (error.response?.data?.error || error.message))
  }
}

const showReviewModal = async (contractId, caregiverId) => {
  const rating = prompt('請輸入評分 (1-5):')
  if (!rating || rating < 1 || rating > 5) {
    alert('評分必須在 1-5 之間')
    return
  }
  const comment = prompt('請輸入評價留言（可選）:') || ''
  try {
    await reviewAPI.createReview({
      contract_id: contractId,
      caregiver_id: caregiverId,
      rating: parseInt(rating),
      comment: comment,
    })
    alert('評價已提交！')
    loadContracts()
  } catch (error) {
    alert('提交評價失敗：' + (error.response?.data?.error || error.message))
  }
}

const loadContracts = async () => {
  loading.value = true
  try {
    contracts.value = await contractAPI.getContracts()
  } catch (error) {
    console.error('Load contracts error:', error)
    alert('載入失敗：' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadContracts()
})
</script>