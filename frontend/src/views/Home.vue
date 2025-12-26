<template>
  <div>
    <Navbar />
    <div class="bg-gray-50">
      <!-- Hero Section -->
      <section class="bg-gradient-to-r from-primary-600 to-primary-800 text-white py-20">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 class="text-4xl md:text-5xl font-bold mb-4">找到最適合的照護者</h2>
          <p class="text-xl mb-8 text-primary-100">專業、可靠、貼心的照護服務平台</p>
          <div class="flex justify-center space-x-4">
            <router-link
              :to="{ name: 'register', query: { role: 'user' } }"
              class="btn-primary bg-white text-primary-600 hover:bg-gray-100"
            >
              尋找照護者
            </router-link>
            <router-link
              :to="{ name: 'register', query: { role: 'caregiver' } }"
              class="btn-secondary bg-transparent border-2 border-white text-white hover:bg-white hover:text-primary-600"
            >
              成為照護者
            </router-link>
          </div>
        </div>
      </section>

      <!-- Features Section -->
      <section class="py-16">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h3 class="text-3xl font-bold text-center mb-12 text-gray-800">為什麼選擇 CareFinder？</h3>
          <div class="grid md:grid-cols-3 gap-8">
            <div class="card text-center">
              <div class="text-4xl mb-4">🔍</div>
              <h4 class="text-xl font-semibold mb-2">智能搜尋</h4>
              <p class="text-gray-600">根據地點、時間、需求快速找到合適的照護者</p>
            </div>
            <div class="card text-center">
              <div class="text-4xl mb-4">✅</div>
              <h4 class="text-xl font-semibold mb-2">證照認證</h4>
              <p class="text-gray-600">所有照護者都經過證照審核，確保專業品質</p>
            </div>
            <div class="card text-center">
              <div class="text-4xl mb-4">⭐</div>
              <h4 class="text-xl font-semibold mb-2">真實評價</h4>
              <p class="text-gray-600">查看其他使用者的真實評價與評分</p>
            </div>
          </div>
        </div>
      </section>

      <!-- CTA Section -->
      <section class="bg-primary-50 py-16">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h3 class="text-3xl font-bold mb-4 text-gray-800">準備開始了嗎？</h3>
          <p class="text-gray-600 mb-8">立即註冊，開始尋找或提供照護服務</p>
          <router-link to="/register" class="btn-primary">免費註冊</router-link>
        </div>
      </section>

      <!-- Footer -->
      <footer class="bg-gray-800 text-white py-8">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p>&copy; 2024 CareFinder. All rights reserved.</p>
        </div>
      </footer>
    </div>
  </div>
</template>

<script setup>
import Navbar from '@/components/Navbar.vue'
</script>