<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <!-- User Dashboard -->
        <div v-if="authStore.isUser">
          <h2 class="text-3xl font-bold mb-6 text-gray-800">歡迎回來！</h2>
          
          <div class="grid md:grid-cols-2 gap-6 mb-8">
            <div class="card">
              <h3 class="text-xl font-semibold mb-4">快速搜尋照護者</h3>
              <p class="text-gray-600 mb-4">根據您的需求尋找合適的照護者</p>
              <router-link to="/caregiver-search" class="btn-primary inline-block">開始搜尋</router-link>
            </div>
            
            <div class="card">
              <h3 class="text-xl font-semibold mb-4">我的合約</h3>
              <p class="text-gray-600 mb-4">查看和管理您的服務合約</p>
              <router-link to="/contracts" class="btn-primary inline-block">查看合約</router-link>
            </div>
          </div>
        </div>

        <!-- Caregiver Dashboard -->
        <div v-if="authStore.isCaregiver">
          <h2 class="text-3xl font-bold mb-6 text-gray-800">照護者儀表板</h2>
          
          <div class="grid md:grid-cols-2 gap-6 mb-8">
            <div class="card">
              <h3 class="text-xl font-semibold mb-4">我的檔案</h3>
              <p class="text-gray-600 mb-4">管理您的照護者檔案和可服務時間</p>
              <router-link to="/caregiver-profile-setup" class="btn-primary inline-block">管理檔案</router-link>
            </div>
            
            <div class="card">
              <h3 class="text-xl font-semibold mb-4">服務時間</h3>
              <p class="text-gray-600 mb-4">設定您的可服務時間段</p>
              <router-link to="/caregiver-availability" class="btn-primary inline-block">設定時間</router-link>
            </div>
            
            <div class="card">
              <h3 class="text-xl font-semibold mb-4">證照管理</h3>
              <p class="text-gray-600 mb-4">上傳和管理您的證照</p>
              <router-link to="/caregiver-licenses" class="btn-primary inline-block">管理證照</router-link>
            </div>
            
            <div class="card">
              <h3 class="text-xl font-semibold mb-4">我的合約</h3>
              <p class="text-gray-600 mb-4">查看和管理您的服務合約</p>
              <router-link to="/contracts" class="btn-primary inline-block">查看合約</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'

const router = useRouter()
const authStore = useAuthStore()

onMounted(() => {
  // Redirect admin to admin dashboard
  if (authStore.isAdmin) {
    router.push('/admin')
  }
})
</script>