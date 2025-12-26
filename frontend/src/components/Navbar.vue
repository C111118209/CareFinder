<template>
  <nav class="bg-white shadow-md sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <div class="flex items-center">
          <router-link
            :to="authStore.isAuthenticated ? '/dashboard' : '/'"
            class="text-2xl font-bold text-primary-600"
          >
            CareFinder
          </router-link>
        </div>

        <!-- Desktop Navigation -->
        <div class="hidden md:flex items-center space-x-8">
          <!-- Navigation Menu -->
          <div v-if="authStore.isAuthenticated" class="flex space-x-4">
            <!-- Admin Navigation -->
            <template v-if="authStore.isAdmin">
              <router-link
                to="/dashboard"
                class="nav-link"
                :class="{ 'active': $route.name === 'dashboard' }"
              >
                儀表板
              </router-link>
              <router-link
                to="/admin"
                class="nav-link"
                :class="{ 'active': $route.name === 'admin-dashboard' }"
              >
                管理後台
              </router-link>
              <router-link
                to="/profile"
                class="nav-link"
                :class="{ 'active': $route.name === 'profile' }"
              >
                個人資訊
              </router-link>
            </template>

            <!-- Caregiver Navigation -->
            <template v-else-if="authStore.isCaregiver">
              <router-link
                to="/dashboard"
                class="nav-link"
                :class="{ 'active': $route.name === 'dashboard' }"
              >
                儀表板
              </router-link>
              <router-link
                to="/caregiver-profile-setup"
                class="nav-link"
                :class="{ 'active': $route.name === 'caregiver-profile-setup' }"
              >
                我的檔案
              </router-link>
              <router-link
                to="/caregiver-availability"
                class="nav-link"
                :class="{ 'active': $route.name === 'caregiver-availability' }"
              >
                服務時間
              </router-link>
              <router-link
                to="/caregiver-licenses"
                class="nav-link"
                :class="{ 'active': $route.name === 'caregiver-licenses' }"
              >
                證照管理
              </router-link>
              <router-link
                to="/contracts"
                class="nav-link"
                :class="{ 'active': $route.name === 'contracts' }"
              >
                我的合約
              </router-link>
              <router-link
                to="/profile"
                class="nav-link"
                :class="{ 'active': $route.name === 'profile' }"
              >
                個人資訊
              </router-link>
            </template>

            <!-- User Navigation -->
            <template v-else-if="authStore.isUser">
              <router-link
                to="/dashboard"
                class="nav-link"
                :class="{ 'active': $route.name === 'dashboard' }"
              >
                儀表板
              </router-link>
              <router-link
                to="/caregiver-search"
                class="nav-link"
                :class="{ 'active': $route.name === 'caregiver-search' }"
              >
                搜尋照護者
              </router-link>
              <router-link
                to="/contracts"
                class="nav-link"
                :class="{ 'active': $route.name === 'contracts' }"
              >
                我的合約
              </router-link>
              <router-link
                to="/profile"
                class="nav-link"
                :class="{ 'active': $route.name === 'profile' }"
              >
                個人資訊
              </router-link>
            </template>
          </div>
        </div>

        <!-- Right Side: User Info & Actions -->
        <div class="flex items-center space-x-4">
          <template v-if="authStore.isAuthenticated">
            <!-- Desktop: Show email -->
            <span class="hidden md:block text-gray-700 text-sm">
              {{ authStore.user?.email || '使用者' }}
            </span>
            <!-- Mobile: Show user icon -->
            <div class="md:hidden relative">
              <button
                @click="showMobileMenu = !showMobileMenu"
                class="p-2 rounded-md text-gray-700 hover:text-primary-600 focus:outline-none"
                aria-label="選單"
              >
                <svg
                  class="h-6 w-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    v-if="!showMobileMenu"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 6h16M4 12h16M4 18h16"
                  />
                  <path
                    v-else
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
            <!-- Desktop: Profile & Logout -->
            <div class="hidden md:flex items-center space-x-2">
              <button
                @click="handleLogout"
                class="text-gray-500 hover:text-red-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              >
                <i class="fa-solid fa-right-from-bracket mr-1"></i> 登出
              </button>
            </div>
          </template>
          <template v-else>
            <router-link
              to="/login"
              class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium"
            >
              登入
            </router-link>
            <router-link to="/register" class="btn-primary">
              註冊
            </router-link>
          </template>
        </div>
      </div>

      <!-- Mobile Menu -->
      <div
        v-if="authStore.isAuthenticated && showMobileMenu"
        class="md:hidden border-t border-gray-200"
      >
        <div class="px-2 pt-2 pb-3 space-y-1">
          <!-- Admin Mobile Navigation -->
          <template v-if="authStore.isAdmin">
            <router-link
              to="/dashboard"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'dashboard' }"
            >
              儀表板
            </router-link>
            <router-link
              to="/admin"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'admin-dashboard' }"
            >
              管理後台
            </router-link>
            <router-link
              to="/profile"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'profile' }"
            >
              個人資訊
            </router-link>
          </template>

          <!-- Caregiver Mobile Navigation -->
          <template v-else-if="authStore.isCaregiver">
            <router-link
              to="/dashboard"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'dashboard' }"
            >
              儀表板
            </router-link>
            <router-link
              to="/caregiver-profile-setup"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'caregiver-profile-setup' }"
            >
              我的檔案
            </router-link>
            <router-link
              to="/caregiver-availability"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'caregiver-availability' }"
            >
              服務時間
            </router-link>
            <router-link
              to="/caregiver-licenses"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'caregiver-licenses' }"
            >
              證照管理
            </router-link>
            <router-link
              to="/contracts"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'contracts' }"
            >
              我的合約
            </router-link>
            <router-link
              to="/profile"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'profile' }"
            >
              個人資訊
            </router-link>
          </template>

          <!-- User Mobile Navigation -->
          <template v-else-if="authStore.isUser">
            <router-link
              to="/dashboard"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'dashboard' }"
            >
              儀表板
            </router-link>
            <router-link
              to="/caregiver-search"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'caregiver-search' }"
            >
              搜尋照護者
            </router-link>
            <router-link
              to="/contracts"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'contracts' }"
            >
              我的合約
            </router-link>
            <router-link
              to="/profile"
              @click="showMobileMenu = false"
              class="mobile-nav-link"
              :class="{ 'active': $route.name === 'profile' }"
            >
              個人資訊
            </router-link>
          </template>

          <!-- Mobile Logout -->
          <div class="pt-4 border-t border-gray-200">
            <div class="px-3 py-2 text-sm text-gray-700">
              {{ authStore.user?.email || '使用者' }}
            </div>
            <button
              @click="handleLogout"
              class="w-full text-left px-3 py-2 text-sm text-gray-500 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
            >
              <i class="fa-solid fa-right-from-bracket mr-2"></i> 登出
            </button>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const showMobileMenu = ref(false)

// Close mobile menu when route changes
watch(() => route.path, () => {
  showMobileMenu.value = false
})

const handleLogout = () => {
  authStore.logout()
  showMobileMenu.value = false
  router.push('/login')
}
</script>

<style scoped>
.nav-link {
  @apply text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors;
}

.nav-link.active {
  @apply text-primary-600 border-b-2 border-primary-600;
}

.mobile-nav-link {
  @apply block px-3 py-2 text-base font-medium text-gray-700 hover:text-primary-600 hover:bg-gray-50 rounded-md transition-colors;
}

.mobile-nav-link.active {
  @apply text-primary-600 bg-primary-50;
}
</style>