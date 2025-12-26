<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          <p class="mt-4 text-gray-600">載入中...</p>
        </div>
        <div v-else-if="profile" class="space-y-6">
          <!-- 基本資訊卡片 -->
          <div class="card">
            <div class="flex justify-between items-start mb-6">
              <div>
                <h1 class="text-3xl font-bold text-gray-800 mb-2">{{ profile.full_name || '未提供姓名' }}</h1>
                <div class="flex items-center space-x-4 text-gray-600">
                  <span>⭐ {{ (profile.avg_rating || 0).toFixed(1) }}</span>
                  <span v-if="profile.address">📍 {{ profile.address }}</span>
                  <span v-if="profile.service_rate">💰 NT$ {{ profile.service_rate }}/小時</span>
                </div>
              </div>
              <router-link
                v-if="authStore.isUser"
                :to="`/create-contract/${profile.user_id}?profile_id=${profile.profile_id || ''}`"
                class="btn-primary inline-block"
              >
                發起合約
              </router-link>
            </div>

            <div v-if="profile.bio" class="mb-6">
              <h2 class="text-xl font-semibold mb-2">自我介紹</h2>
              <p class="text-gray-700 whitespace-pre-line">{{ profile.bio }}</p>
            </div>

            <div class="grid md:grid-cols-2 gap-6">
              <div>
                <h2 class="text-xl font-semibold mb-3">基本資訊</h2>
                <div class="space-y-2 text-gray-700">
                  <p><span class="font-medium">性別：</span>{{ getGenderText(profile.gender) }}</p>
                  <p v-if="profile.phone"><span class="font-medium">電話：</span>{{ profile.phone }}</p>
                  <p v-if="profile.address"><span class="font-medium">服務地點：</span>{{ profile.address }}</p>
                </div>
              </div>

              <div v-if="profile.availabilities && profile.availabilities.length > 0">
                <h2 class="text-xl font-semibold mb-3">可服務時間</h2>
                <div class="space-y-2">
                  <div
                    v-for="avail in profile.availabilities"
                    :key="avail.availability_id"
                    class="text-gray-700"
                  >
                    <p>{{ getDayName(avail.day_of_week) }} {{ avail.start_time }} - {{ avail.end_time }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="profile.licenses && profile.licenses.length > 0" class="mt-6">
              <h2 class="text-xl font-semibold mb-3">證照</h2>
              <div class="space-y-2">
                <div
                  v-for="license in profile.licenses"
                  :key="license.license_id"
                  class="bg-green-50 border border-green-200 rounded-lg p-3"
                >
                  <p class="font-medium text-green-800">{{ license.name }}</p>
                  <p class="text-sm text-green-600">有效期至：{{ formatDate(license.expiry_date) }}</p>
                  <p
                    v-if="license.status === 'approved'"
                    class="text-xs text-green-700 mt-1"
                  >
                    ✅ 已認證
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- 評價區塊 -->
          <div class="card">
            <h2 class="text-2xl font-semibold mb-4">評價 ({{ reviews.length }})</h2>
            <div v-if="reviews.length === 0" class="text-gray-600">尚無評價</div>
            <div v-else class="space-y-4">
              <div
                v-for="review in reviews"
                :key="review.review_id"
                class="border-b border-gray-200 pb-4 last:border-0"
              >
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center">
                    <span v-html="generateStars(review.rating)"></span>
                    <span class="ml-2 text-gray-600">{{ formatDate(review.created_at) }}</span>
                  </div>
                </div>
                <p v-if="review.comment" class="text-gray-700 mt-2">{{ review.comment }}</p>
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
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { caregiverAPI, reviewAPI } from '@/services/api'
import Navbar from '@/components/Navbar.vue'

const route = useRoute()
const authStore = useAuthStore()

const profile = ref(null)
const reviews = ref([])
const loading = ref(false)

const formatDate = (dateString) => {
  if (!dateString) return '無資料'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const getGenderText = (gender) => {
  if (gender === 'male') return '男性'
  if (gender === 'female') return '女性'
  return '未提供'
}

const getDayName = (dayOfWeek) => {
  const days = ['', '週一', '週二', '週三', '週四', '週五', '週六', '週日']
  return days[dayOfWeek] || `週${dayOfWeek}`
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

onMounted(async () => {
  loading.value = true
  try {
    const caregiverId = route.params.id
    profile.value = await caregiverAPI.getProfile(caregiverId)
    reviews.value = await reviewAPI.getCaregiverReviews(caregiverId)
  } catch (error) {
    console.error('Load error:', error)
    alert('載入失敗：' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
})
</script>