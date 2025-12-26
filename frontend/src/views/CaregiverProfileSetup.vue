<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen flex flex-col">
      <div class="flex-grow max-w-3xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-10">
        <div class="text-center mb-10">
          <h2 class="text-3xl font-bold text-gray-900 mb-2">建立您的照護者檔案</h2>
          <p class="text-gray-500">完善的資料能讓案主更信任您，提高接案成功率。</p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-8">
          <!-- 基本資料 -->
          <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
            <div class="px-6 py-4 bg-gray-50 border-b border-gray-100 flex items-center gap-3">
              <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600">
                <i class="fa-solid fa-user"></i>
              </div>
              <h3 class="text-lg font-semibold text-gray-800">基本資料</h3>
            </div>

            <div class="p-6 md:p-8 space-y-6">
              <div class="grid md:grid-cols-2 gap-6">
                <div class="input-group">
                  <label class="block text-sm font-medium text-gray-700 mb-1.5">
                    真實姓名 <span class="text-red-500">*</span>
                  </label>
                  <div class="relative">
                    <span class="absolute left-3 top-3 text-gray-400">
                      <i class="fa-solid fa-signature"></i>
                    </span>
                    <input
                      v-model="formData.full_name"
                      type="text"
                      required
                      class="input-field pl-10"
                      placeholder="請輸入姓名"
                    />
                  </div>
                </div>

                <div class="input-group">
                  <label class="block text-sm font-medium text-gray-700 mb-1.5">性別</label>
                  <div class="relative">
                    <span class="absolute left-3 top-3 text-gray-400">
                      <i class="fa-solid fa-venus-mars"></i>
                    </span>
                    <select v-model="formData.gender" class="input-field pl-10 appearance-none bg-white">
                      <option value="">請選擇</option>
                      <option value="male">男性</option>
                      <option value="female">女性</option>
                      <option value="other">其他</option>
                    </select>
                  </div>
                </div>

                <div class="input-group">
                  <label class="block text-sm font-medium text-gray-700 mb-1.5">
                    聯絡電話 <span class="text-red-500">*</span>
                  </label>
                  <div class="relative">
                    <span class="absolute left-3 top-3 text-gray-400">
                      <i class="fa-solid fa-phone"></i>
                    </span>
                    <input
                      v-model="formData.phone"
                      type="tel"
                      required
                      class="input-field pl-10"
                      placeholder="0912345678"
                    />
                  </div>
                </div>

                <div class="input-group">
                  <label class="block text-sm font-medium text-gray-700 mb-1.5">
                    服務/居住地點 <span class="text-red-500">*</span>
                  </label>
                  <div class="relative">
                    <span class="absolute left-3 top-3 text-gray-400">
                      <i class="fa-solid fa-location-dot"></i>
                    </span>
                    <input
                      v-model="formData.address"
                      type="text"
                      required
                      class="input-field pl-10"
                      placeholder="例如：台北市信義區"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 服務詳情 -->
          <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
            <div class="px-6 py-4 bg-gray-50 border-b border-gray-100 flex items-center gap-3">
              <div class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center text-green-600">
                <i class="fa-solid fa-briefcase"></i>
              </div>
              <h3 class="text-lg font-semibold text-gray-800">服務詳情</h3>
            </div>

            <div class="p-6 md:p-8 space-y-6">
              <div class="input-group">
                <label class="block text-sm font-medium text-gray-700 mb-1.5">
                  期望時薪 (NT$) <span class="text-red-500">*</span>
                </label>
                <div class="relative max-w-xs">
                  <span class="absolute left-3 top-3 text-gray-400">
                    <i class="fa-solid fa-sack-dollar"></i>
                  </span>
                  <input
                    v-model.number="formData.service_rate"
                    type="number"
                    required
                    min="0"
                    step="10"
                    class="input-field pl-10"
                    placeholder="500"
                  />
                </div>
              </div>

              <div class="input-group">
                <label class="block text-sm font-medium text-gray-700 mb-1.5">自我介紹與經驗</label>
                <textarea
                  v-model="formData.bio"
                  rows="6"
                  class="input-field resize-y"
                  placeholder="請詳細描述您的照護經驗、擁有的證照、或是您的服務理念..."
                ></textarea>
                <p class="text-xs text-gray-400 mt-1 text-right">建議字數：100字以上</p>
              </div>
            </div>
          </div>

          <div v-if="error" class="bg-red-50 text-red-600 p-4 rounded-lg text-sm flex items-center gap-2">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error }}</span>
          </div>

          <div class="flex items-center justify-end space-x-4 pt-4 border-t border-gray-200">
            <button
              type="button"
              @click="$router.back()"
              class="px-6 py-2.5 rounded-lg text-gray-700 hover:bg-gray-100 font-medium transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="bg-blue-600 hover:bg-blue-700 text-white px-8 py-2.5 rounded-lg font-medium shadow-md hover:shadow-lg transition-all transform hover:-translate-y-0.5 flex items-center gap-2"
            >
              <i class="fa-solid fa-save"></i>
              <span>{{ loading ? '處理中...' : (isUpdate ? '更新檔案' : '儲存檔案') }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { caregiverAPI } from '@/services/api'
import Navbar from '@/components/Navbar.vue'

const formData = ref({
  full_name: '',
  gender: '',
  phone: '',
  address: '',
  bio: '',
  service_rate: 0,
})

const error = ref('')
const success = ref('')
const loading = ref(false)
const isUpdate = ref(false)

onMounted(() => {
  loadExistingProfile()
})

const handleSubmit = async () => {
  error.value = ''
  loading.value = true

  try {
    if (isUpdate.value) {
      await caregiverAPI.updateProfile(formData.value)
      alert('檔案更新成功！')
    } else {
      try {
        await caregiverAPI.createProfile(formData.value)
        alert('檔案建立成功！')
        isUpdate.value = true
      } catch (err) {
        if (err.response?.status === 409 || err.message.includes('already exists')) {
          await caregiverAPI.updateProfile(formData.value)
          alert('檔案更新成功！')
          isUpdate.value = true
        } else {
          throw err
        }
      }
    }
    // Reload to show updated data
    await loadExistingProfile()
  } catch (err) {
    error.value = err.response?.data?.error || err.message || '操作失敗，請稍後再試'
  } finally {
    loading.value = false
  }
}

const loadExistingProfile = async () => {
  try {
    const profile = await caregiverAPI.getMyProfile()
    if (profile) {
      formData.value = {
        full_name: profile.full_name || '',
        gender: profile.gender || '',
        phone: profile.phone || '',
        address: profile.address || '',
        bio: profile.bio || '',
        service_rate: profile.service_rate || 0,
      }
      isUpdate.value = true
    }
  } catch (error) {
    // Profile doesn't exist yet
    isUpdate.value = false
  }
}
</script>