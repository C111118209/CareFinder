<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen flex items-center justify-center py-12">
      <div class="max-w-md w-full">
        <div class="card">
          <div class="text-center mb-6">
            <h1 class="text-3xl font-bold text-primary-600 mb-2">CareFinder</h1>
            <h2 class="text-2xl font-semibold text-gray-800">註冊</h2>
          </div>

          <form @submit.prevent="handleRegister" class="space-y-4">
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 mb-1">電子郵件</label>
              <input
                v-model="email"
                type="email"
                id="email"
                required
                class="input-field"
                placeholder="your@email.com"
              />
            </div>

            <div>
              <label for="password" class="block text-sm font-medium text-gray-700 mb-1">密碼</label>
              <input
                v-model="password"
                type="password"
                id="password"
                required
                class="input-field"
                placeholder="••••••••"
                minlength="6"
              />
            </div>

            <div>
              <label for="role" class="block text-sm font-medium text-gray-700 mb-1">註冊身份</label>
              <select v-model="role" id="role" required class="input-field">
                <option value="user">使用者/家屬</option>
                <option value="caregiver">照護者</option>
              </select>
            </div>

            <div v-if="error" class="text-red-600 text-sm">{{ error }}</div>

            <button type="submit" :disabled="loading" class="w-full btn-primary">
              {{ loading ? '註冊中...' : '註冊' }}
            </button>
          </form>

          <div class="mt-6 text-center">
            <p class="text-gray-600">
              已經有帳號？
              <router-link to="/login" class="text-primary-600 hover:text-primary-700 font-medium">立即登入</router-link>
            </p>
          </div>

          <div class="mt-4 text-center">
            <router-link to="/" class="text-gray-500 hover:text-gray-700 text-sm">返回首頁</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const role = ref('user')
const error = ref('')
const loading = ref(false)

onMounted(() => {
  // Set role from URL parameter
  const roleParam = route.query.role
  if (roleParam === 'user' || roleParam === 'caregiver') {
    role.value = roleParam
  }
})

const handleRegister = async () => {
  error.value = ''
  loading.value = true
  
  try {
    await authStore.register(email.value, password.value, role.value)
    router.push({ name: 'login', query: { registered: 'true' } })
  } catch (err) {
    error.value = err.response?.data?.error || err.message || '註冊失敗，請稍後再試'
  } finally {
    loading.value = false
  }
}
</script>