<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h2 class="text-3xl font-bold mb-6 text-gray-800">個人資訊</h2>

        <!-- Profile Info Section -->
        <div class="card mb-6">
          <h3 class="text-xl font-semibold mb-4">基本資訊</h3>
          <form @submit.prevent="handleUpdate" class="space-y-6">
            <div class="flex items-center space-x-6">
              <!-- Avatar Display -->
              <div class="flex-shrink-0">
                <div class="relative">
                  <img
                    v-if="avatarPreview || formData.avatar_url"
                    :src="avatarPreview || formData.avatar_url"
                    alt="頭像"
                    class="w-24 h-24 rounded-full object-cover border-4 border-gray-200"
                  />
                  <div
                    v-else
                    class="w-24 h-24 rounded-full bg-gray-300 flex items-center justify-center border-4 border-gray-200"
                  >
                    <svg class="w-12 h-12 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                      />
                    </svg>
                  </div>
                  <label
                    for="avatarInput"
                    class="absolute bottom-0 right-0 bg-primary-600 text-white rounded-full p-2 cursor-pointer hover:bg-primary-700 transition-colors"
                    title="上傳頭像"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
                      />
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
                      />
                    </svg>
                  </label>
                  <input
                    ref="avatarInputRef"
                    type="file"
                    id="avatarInput"
                    accept="image/*"
                    class="hidden"
                    @change="handleAvatarUpload"
                  />
                </div>
              </div>

              <!-- User Info -->
              <div class="flex-1 space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">暱稱</label>
                  <input
                    v-model="formData.nickname"
                    type="text"
                    class="input-field"
                    placeholder="請輸入您的暱稱"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">電子郵件</label>
                  <input
                    v-model="formData.email"
                    type="email"
                    class="input-field bg-gray-100"
                    readonly
                    placeholder="電子郵件"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">角色</label>
                  <input
                    :value="getRoleText(formData.role)"
                    type="text"
                    class="input-field bg-gray-100"
                    readonly
                  />
                </div>
              </div>
            </div>

            <div v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded-md border border-red-100">
              {{ error }}
            </div>
            <div v-if="success" class="text-green-600 text-sm bg-green-50 p-3 rounded-md border border-green-100">
              {{ success }}
            </div>

            <div class="flex justify-end">
              <button type="submit" :disabled="loading" class="btn-primary">
                {{ loading ? '儲存中...' : '儲存變更' }}
              </button>
            </div>
          </form>
        </div>

        <!-- Change Password Section -->
        <div class="card">
          <h3 class="text-xl font-semibold mb-4">修改密碼</h3>
          <form @submit.prevent="handlePasswordChange" class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                目前密碼 <span class="text-red-500">*</span>
              </label>
              <input
                v-model="passwordData.oldPassword"
                type="password"
                required
                class="input-field"
                placeholder="請輸入目前密碼"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                新密碼 <span class="text-red-500">*</span>
              </label>
              <input
                v-model="passwordData.newPassword"
                type="password"
                required
                class="input-field"
                placeholder="請輸入新密碼（至少6個字元）"
                minlength="6"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                確認新密碼 <span class="text-red-500">*</span>
              </label>
              <input
                v-model="passwordData.confirmPassword"
                type="password"
                required
                class="input-field"
                placeholder="請再次輸入新密碼"
              />
            </div>

            <div v-if="passwordError" class="text-red-600 text-sm bg-red-50 p-3 rounded-md border border-red-100">
              {{ passwordError }}
            </div>
            <div v-if="passwordSuccess" class="text-green-600 text-sm bg-green-50 p-3 rounded-md border border-green-100">
              {{ passwordSuccess }}
            </div>

            <div class="flex justify-end">
              <button type="submit" :disabled="passwordLoading" class="btn-primary">
                {{ passwordLoading ? '處理中...' : '修改密碼' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { userAPI } from '@/services/api'
import Navbar from '@/components/Navbar.vue'

const authStore = useAuthStore()
const avatarInputRef = ref(null)

const formData = ref({
  email: '',
  nickname: '',
  avatar_url: '',
  role: '',
})

const passwordData = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const avatarPreview = ref('')
const error = ref('')
const success = ref('')
const passwordError = ref('')
const passwordSuccess = ref('')
const loading = ref(false)
const passwordLoading = ref(false)

const getRoleText = (role) => {
  const roleMap = {
    user: '家屬',
    caregiver: '照護者',
    admin: '管理員',
  }
  return roleMap[role] || role
}

const handleAvatarUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // Validate file type
  if (!file.type.startsWith('image/')) {
    error.value = '請選擇圖片檔案'
    setTimeout(() => {
      error.value = ''
    }, 3000)
    return
  }

  // Validate file size (max 5MB)
  if (file.size > 5 * 1024 * 1024) {
    error.value = '圖片大小不能超過5MB'
    setTimeout(() => {
      error.value = ''
    }, 3000)
    return
  }

  try {
    // Convert to base64
    const reader = new FileReader()
    reader.onload = async (event) => {
      const base64Image = event.target.result

      // Update user with new avatar URL
      try {
        loading.value = true
        const updatedUser = await userAPI.updateUser(authStore.user.user_id, {
          avatar_url: base64Image,
        })
        formData.value.avatar_url = updatedUser.avatar_url || base64Image
        avatarPreview.value = base64Image
        authStore.updateUser(updatedUser)
        success.value = '頭像更新成功'
        setTimeout(() => {
          success.value = ''
        }, 3000)
      } catch (err) {
        console.error('Failed to update avatar:', err)
        error.value = err.response?.data?.error || err.message || '更新頭像失敗'
        setTimeout(() => {
          error.value = ''
        }, 3000)
      } finally {
        loading.value = false
      }
    }
    reader.onerror = () => {
      error.value = '處理圖片失敗'
      setTimeout(() => {
        error.value = ''
      }, 3000)
    }
    reader.readAsDataURL(file)
  } catch (err) {
    console.error('Failed to process image:', err)
    error.value = '處理圖片失敗'
    setTimeout(() => {
      error.value = ''
    }, 3000)
  }
}

const handleUpdate = async () => {
  error.value = ''
  success.value = ''
  loading.value = true

  try {
    const updateData = {}
    const currentNickname = authStore.user?.nickname || ''
    if (formData.value.nickname !== currentNickname) {
      updateData.nickname = formData.value.nickname
    }

    if (Object.keys(updateData).length === 0) {
      success.value = '沒有變更需要儲存'
      setTimeout(() => {
        success.value = ''
      }, 3000)
      return
    }

    const updatedUser = await userAPI.updateUser(authStore.user.user_id, updateData)
    authStore.updateUser(updatedUser)
    formData.value.nickname = updatedUser.nickname || ''
    success.value = '個人資訊更新成功'
    setTimeout(() => {
      success.value = ''
    }, 3000)
  } catch (err) {
    error.value = err.response?.data?.error || err.message || '更新失敗'
    setTimeout(() => {
      error.value = ''
    }, 3000)
  } finally {
    loading.value = false
  }
}

const handlePasswordChange = async () => {
  passwordError.value = ''
  passwordSuccess.value = ''

  // Validate passwords match
  if (passwordData.value.newPassword !== passwordData.value.confirmPassword) {
    passwordError.value = '新密碼與確認密碼不一致'
    return
  }

  // Validate password length
  if (passwordData.value.newPassword.length < 6) {
    passwordError.value = '新密碼長度至少需要6個字元'
    return
  }

  passwordLoading.value = true

  try {
    await userAPI.updatePassword(
      authStore.user.user_id,
      passwordData.value.oldPassword,
      passwordData.value.newPassword
    )
    passwordSuccess.value = '密碼修改成功'
    passwordData.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    }
    setTimeout(() => {
      passwordSuccess.value = ''
    }, 3000)
  } catch (err) {
    passwordError.value = err.response?.data?.error || err.message || '修改密碼失敗'
  } finally {
    passwordLoading.value = false
  }
}

onMounted(async () => {
  try {
    const user = await userAPI.getCurrentUser()
    formData.value = {
      email: user.email || '',
      nickname: user.nickname || '',
      avatar_url: user.avatar_url || '',
      role: user.role || '',
    }
    
    // Set avatar preview if exists
    if (user.avatar_url) {
      avatarPreview.value = user.avatar_url
    }
    
    authStore.updateUser(user)
  } catch (err) {
    error.value = '載入資料失敗'
    setTimeout(() => {
      error.value = ''
    }, 3000)
  }
})
</script>