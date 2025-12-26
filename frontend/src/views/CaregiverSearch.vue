<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h2 class="text-3xl font-bold mb-6 text-gray-800">搜尋照護者</h2>

        <!-- 搜尋條件 -->
        <div class="card mb-6 p-4 bg-white rounded shadow">
          <h3 class="text-lg font-semibold mb-4">搜尋條件</h3>
          <form @submit.prevent="performSearch" class="grid md:grid-cols-4 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">地點</label>
              <input
                v-model="filters.address"
                type="text"
                class="input-field border p-2 w-full"
                placeholder="例如：台北"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">性別</label>
              <select v-model="filters.gender" class="input-field border p-2 w-full">
                <option value="">不限</option>
                <option value="male">男性</option>
                <option value="female">女性</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">最低評分</label>
              <input
                v-model.number="filters.min_rating"
                type="number"
                min="0"
                max="5"
                step="0.1"
                class="input-field border p-2 w-full"
                placeholder="0"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">開始日期</label>
              <input
                v-model="filters.start_date"
                type="date"
                class="input-field border p-2 w-full"
              />
            </div>
            <div class="md:col-span-4 mt-4">
              <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                搜尋
              </button>
              <button
                type="button"
                @click="resetSearch"
                class="bg-gray-500 text-white px-4 py-2 rounded ml-2 hover:bg-gray-600"
              >
                重置
              </button>
            </div>
          </form>
        </div>

        <!-- 搜尋結果 -->
        <div v-if="loading" class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p class="mt-4 text-gray-600">載入中...</p>
        </div>
        <div v-else-if="caregivers.length === 0" class="text-center py-8 text-gray-600">
          {{ hasSearched ? '沒有找到符合條件的照護者' : '請輸入搜尋條件開始搜尋' }}
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="caregiver in caregivers"
            :key="caregiver.profile_id || caregiver.user_id"
            class="card bg-white p-4 rounded shadow hover:shadow-xl transition-shadow cursor-pointer mb-4"
            @click="viewProfile(caregiver.user_id)"
          >
            <div class="flex justify-between items-start">
              <div class="flex-1">
                <h3 class="text-xl font-semibold text-gray-800 mb-2">
                  {{ caregiver.full_name || '未提供姓名' }}
                </h3>
                <div class="space-y-2 text-gray-600">
                  <p>
                    <span class="font-medium">地點：</span>{{ caregiver.address || '未提供' }}
                  </p>
                  <p>
                    <span class="font-medium">性別：</span>{{ getGenderText(caregiver.gender) }}
                  </p>
                  <p>
                    <span class="font-medium">時薪：</span>NT$ {{ caregiver.service_rate || '0' }}
                  </p>
                  <div class="flex items-center">
                    <span class="font-medium mr-2">評分：</span>
                    <div class="flex items-center">
                      <span v-html="generateStars(caregiver.avg_rating || 0)"></span>
                      <span class="ml-2">{{ (caregiver.avg_rating || 0).toFixed(1) }}</span>
                    </div>
                  </div>
                </div>
              </div>
              <div class="ml-4">
                <router-link
                  :to="`/create-contract/${caregiver.user_id}?profile_id=${caregiver.profile_id || ''}`"
                  class="bg-blue-600 text-white px-4 py-2 rounded inline-block hover:bg-blue-700"
                  @click.stop
                >
                  發起合約
                </router-link>
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
import { caregiverAPI } from '@/services/api'
import Navbar from '@/components/Navbar.vue'

const router = useRouter()

const filters = ref({
  address: '',
  gender: '',
  min_rating: null,
  start_date: '',
})

const caregivers = ref([])
const loading = ref(false)
const hasSearched = ref(false)

const getGenderText = (gender) => {
  if (gender === 'male') return '男性'
  if (gender === 'female') return '女性'
  return '未提供'
}

const generateStars = (rating) => {
  const fullStars = Math.floor(rating)
  const hasHalfStar = rating % 1 >= 0.5
  let stars = ''

  for (let i = 0; i < fullStars; i++) {
    stars += '<span class="text-yellow-400">★</span>'
  }
  if (hasHalfStar) {
    stars += '<span class="text-yellow-400">☆</span>'
  }
  for (let i = fullStars + (hasHalfStar ? 1 : 0); i < 5; i++) {
    stars += '<span class="text-gray-300">★</span>'
  }
  return stars
}

const performSearch = async () => {
  loading.value = true
  hasSearched.value = true
  try {
    const params = {}
    if (filters.value.address) params.address = filters.value.address
    if (filters.value.gender) params.gender = filters.value.gender
    if (filters.value.min_rating) params.min_rating = filters.value.min_rating
    if (filters.value.start_date) params.start_date = filters.value.start_date

    const result = await caregiverAPI.search(params)
    caregivers.value = result.profiles || result || []
  } catch (error) {
    console.error('Search error:', error)
    alert('搜尋失敗：' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  filters.value = {
    address: '',
    gender: '',
    min_rating: null,
    start_date: '',
  }
  caregivers.value = []
  hasSearched.value = false
}

const viewProfile = (profileId) => {
  router.push(`/caregiver/${profileId}`)
}

onMounted(() => {
  // Optionally perform initial search
  // performSearch()
})
</script>