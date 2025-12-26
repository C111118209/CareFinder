<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h2 class="text-3xl font-bold mb-6 text-gray-800">發起新合約</h2>

        <div class="card">
          <div v-if="loadingCaregiver" class="mb-6 p-4 bg-primary-50 rounded-lg">
            <p class="text-gray-600">載入照護者資訊中...</p>
          </div>
          <div v-else-if="caregiverInfo" class="mb-6 p-4 bg-primary-50 rounded-lg">
            <h3 class="font-semibold text-gray-800 mb-2">{{ caregiverInfo.full_name || '照護者' }}</h3>
            <p class="text-sm text-gray-600">
              <span v-if="caregiverInfo.address">📍 {{ caregiverInfo.address }}</span>
              <span v-if="caregiverInfo.service_rate" class="ml-4">💰 NT$ {{ caregiverInfo.service_rate }}/小時</span>
            </p>
            <p v-if="caregiverInfo.avg_rating" class="text-sm text-gray-600 mt-1">
              ⭐ 評分：{{ caregiverInfo.avg_rating.toFixed(1) }}
            </p>
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-6">
            <input type="hidden" :value="caregiverId" />

            <div class="grid md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  合約開始日期 <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="formData.start_date"
                  type="date"
                  required
                  class="input-field"
                  :min="minDate"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  合約結束日期 <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="formData.end_date"
                  type="date"
                  required
                  class="input-field"
                  :min="minDate"
                />
              </div>
            </div>

            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <p class="text-sm text-blue-800">
                <strong>注意：</strong>合約期限最長為 180 天（6個月）。結束日期必須晚於開始日期。
              </p>
            </div>

            <div v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded-md border border-red-100">
              {{ error }}
            </div>

            <div class="flex justify-end space-x-4">
              <button
                type="button"
                @click="$router.back()"
                class="btn-secondary"
              >
                取消
              </button>
              <button
                type="submit"
                :disabled="loading"
                class="btn-primary"
              >
                {{ loading ? '處理中...' : '發起合約' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { contractAPI, caregiverAPI } from '@/services/api'
import Navbar from '@/components/Navbar.vue'

const route = useRoute()
const router = useRouter()

const caregiverId = ref(null)
const caregiverInfo = ref(null)
const loadingCaregiver = ref(false)
const formData = ref({
  start_date: '',
  end_date: '',
})
const error = ref('')
const loading = ref(false)

const minDate = computed(() => {
  return new Date().toISOString().split('T')[0]
})

onMounted(async () => {
  // Get caregiver ID from route params or query
  caregiverId.value = route.params.caregiverId || route.query.caregiver_id

  if (!caregiverId.value) {
    alert('缺少照護者資訊')
    router.push('/caregiver-search')
    return
  }

  // Load caregiver info
  loadingCaregiver.value = true
  try {
    caregiverInfo.value = await caregiverAPI.getProfile(caregiverId.value)
  } catch (error) {
    console.error('Load caregiver error:', error)
  } finally {
    loadingCaregiver.value = false
  }

  // Set minimum date to today
  const today = new Date().toISOString().split('T')[0]
  formData.value.start_date = today
  formData.value.end_date = today
})

const handleSubmit = async () => {
  error.value = ''

  // Validate dates
  if (new Date(formData.value.start_date) >= new Date(formData.value.end_date)) {
    error.value = '結束日期必須晚於開始日期'
    return
  }

  // Check max duration (180 days)
  const start = new Date(formData.value.start_date)
  const end = new Date(formData.value.end_date)
  const diffTime = Math.abs(end - start)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays > 180) {
    error.value = '合約期限最長為 180 天'
    return
  }

  loading.value = true

  try {
    const contract = await contractAPI.createContract({
      caregiver_id: parseInt(caregiverId.value),
      start_date: formData.value.start_date,
      end_date: formData.value.end_date,
    })
    alert('合約已發起！')
    router.push(`/contract/${contract.contract_id}`)
  } catch (err) {
    error.value = err.response?.data?.error || err.message || '建立合約失敗'
  } finally {
    loading.value = false
  }
}
</script>