<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div v-if="authStore.isUser">
          <div class="mb-8">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h2 class="text-3xl font-bold text-gray-800">
                  歡迎回來，{{
                    authStore.user?.nickname ||
                    authStore.user?.email ||
                    "使用者"
                  }}！
                </h2>
                <p class="text-gray-600 mt-1">以下是您的服務概覽</p>
              </div>
              <div
                v-if="authStore.user?.avatar_url"
                class="w-16 h-16 rounded-full overflow-hidden border-4 border-white shadow-lg"
              >
                <img
                  :src="authStore.user.avatar_url"
                  alt="頭像"
                  class="w-full h-full object-cover"
                />
              </div>
              <div
                v-else
                class="w-16 h-16 rounded-full bg-primary-100 flex items-center justify-center border-4 border-white shadow-lg"
              >
                <span class="text-2xl text-primary-600 font-bold">
                  {{
                    (authStore.user?.nickname ||
                      authStore.user?.email ||
                      "U")[0].toUpperCase()
                  }}
                </span>
              </div>
            </div>
          </div>

          <div class="mb-8">
            <h3 class="text-xl font-semibold text-gray-800 mb-4">快速操作</h3>
            <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
              <router-link
                to="/caregiver-search"
                class="group bg-white rounded-xl shadow-md hover:shadow-xl transition-all p-4 border border-gray-200 hover:border-primary-300"
              >
                <div
                  class="flex flex-col md:flex-row items-center md:space-x-4 text-center md:text-left"
                >
                  <div
                    class="bg-blue-100 group-hover:bg-blue-200 rounded-lg p-3 transition-colors mb-2 md:mb-0"
                  >
                    <i class="fa-solid fa-search text-blue-600 text-xl"></i>
                  </div>
                  <div>
                    <h4
                      class="font-semibold text-gray-800 group-hover:text-primary-600 transition-colors"
                    >
                      搜尋照護者
                    </h4>
                    <p class="text-xs md:text-sm text-gray-600 hidden md:block">
                      尋找合適的照護服務
                    </p>
                  </div>
                </div>
              </router-link>

              <router-link
                to="/contracts"
                class="group bg-white rounded-xl shadow-md hover:shadow-xl transition-all p-4 border border-gray-200 hover:border-primary-300"
              >
                <div
                  class="flex flex-col md:flex-row items-center md:space-x-4 text-center md:text-left"
                >
                  <div
                    class="bg-green-100 group-hover:bg-green-200 rounded-lg p-3 transition-colors mb-2 md:mb-0"
                  >
                    <i
                      class="fa-solid fa-file-contract text-green-600 text-xl"
                    ></i>
                  </div>
                  <div>
                    <h4
                      class="font-semibold text-gray-800 group-hover:text-primary-600 transition-colors"
                    >
                      我的合約
                    </h4>
                    <p class="text-xs md:text-sm text-gray-600 hidden md:block">
                      查看和管理合約
                    </p>
                  </div>
                </div>
              </router-link>

              <router-link
                to="/profile"
                class="group bg-white rounded-xl shadow-md hover:shadow-xl transition-all p-4 border border-gray-200 hover:border-primary-300"
              >
                <div
                  class="flex flex-col md:flex-row items-center md:space-x-4 text-center md:text-left"
                >
                  <div
                    class="bg-purple-100 group-hover:bg-purple-200 rounded-lg p-3 transition-colors mb-2 md:mb-0"
                  >
                    <i class="fa-solid fa-user text-purple-600 text-xl"></i>
                  </div>
                  <div>
                    <h4
                      class="font-semibold text-gray-800 group-hover:text-primary-600 transition-colors"
                    >
                      個人資料
                    </h4>
                    <p class="text-xs md:text-sm text-gray-600 hidden md:block">
                      管理您的帳號設定
                    </p>
                  </div>
                </div>
              </router-link>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div
              class="bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl shadow-lg p-6 text-white"
            >
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-blue-100 text-sm mb-1">進行中的合約</p>
                  <p class="text-3xl font-bold">
                    {{ stats.activeContracts || 0 }}
                  </p>
                </div>
                <div class="bg-white/20 rounded-full p-3">
                  <i class="fa-solid fa-file-contract text-2xl"></i>
                </div>
              </div>
            </div>

            <div
              class="bg-gradient-to-br from-green-500 to-green-600 rounded-xl shadow-lg p-6 text-white"
            >
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-green-100 text-sm mb-1">已完成的合約</p>
                  <p class="text-3xl font-bold">
                    {{ stats.completedContracts || 0 }}
                  </p>
                </div>
                <div class="bg-white/20 rounded-full p-3">
                  <i class="fa-solid fa-check-circle text-2xl"></i>
                </div>
              </div>
            </div>

            <div
              class="bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl shadow-lg p-6 text-white"
            >
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-purple-100 text-sm mb-1">待處理的合約</p>
                  <p class="text-3xl font-bold">
                    {{ stats.pendingContracts || 0 }}
                  </p>
                </div>
                <div class="bg-white/20 rounded-full p-3">
                  <i class="fa-solid fa-clock text-2xl"></i>
                </div>
              </div>
            </div>
          </div>

          <div
            class="bg-white rounded-xl shadow-md border border-gray-200 overflow-hidden"
          >
            <div
              class="px-6 py-4 border-b border-gray-200 flex justify-between items-center"
            >
              <h3 class="text-xl font-semibold text-gray-800">最近的合約</h3>
              <router-link
                to="/contracts"
                class="text-sm text-primary-600 hover:text-primary-700 font-medium"
              >
                查看全部 →
              </router-link>
            </div>
            <div v-if="loadingContracts" class="p-8 text-center">
              <div
                class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
              ></div>
              <p class="mt-2 text-gray-600">載入中...</p>
            </div>
            <div
              v-else-if="recentContracts.length === 0"
              class="p-8 text-center text-gray-500"
            >
              <i
                class="fa-solid fa-file-circle-question text-4xl mb-3 opacity-50"
              ></i>
              <p>目前沒有合約</p>
              <router-link
                to="/caregiver-search"
                class="text-primary-600 hover:text-primary-700 mt-2 inline-block"
              >
                開始搜尋照護者 →
              </router-link>
            </div>
            <div v-else class="divide-y divide-gray-200">
              <div
                v-for="contract in recentContracts"
                :key="contract.contract_id"
                class="p-6 hover:bg-gray-50 transition-colors cursor-pointer"
                @click="$router.push(`/contract/${contract.contract_id}`)"
              >
                <div class="flex items-center justify-between">
                  <div class="flex-1">
                    <div class="flex items-center space-x-3 mb-2">
                      <h4 class="font-semibold text-gray-800">
                        合約 #{{ contract.contract_id }}
                      </h4>
                      <span
                        :class="{
                          'bg-yellow-100 text-yellow-800':
                            contract.status === 'pending',
                          'bg-green-100 text-green-800':
                            contract.status === 'active',
                          'bg-blue-100 text-blue-800':
                            contract.status === 'completed',
                          'bg-red-100 text-red-800':
                            contract.status === 'canceled',
                        }"
                        class="px-2 py-1 rounded-full text-xs font-medium"
                      >
                        {{ getStatusText(contract.status) }}
                      </span>
                    </div>
                    <p class="text-sm text-gray-600">
                      {{ getCaregiverName(contract) }} ·
                      {{ formatDate(contract.start_date) }} -
                      {{ formatDate(contract.end_date) }}
                    </p>
                  </div>
                  <i class="fa-solid fa-chevron-right text-gray-400"></i>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="authStore.isCaregiver">
          <div class="mb-8">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h2 class="text-3xl font-bold text-gray-800">照護者儀表板</h2>
                <p class="text-gray-600 mt-1">管理您的照護服務</p>
              </div>
              <div
                v-if="authStore.user?.avatar_url"
                class="w-16 h-16 rounded-full overflow-hidden border-4 border-white shadow-lg"
              >
                <img
                  :src="authStore.user.avatar_url"
                  alt="頭像"
                  class="w-full h-full object-cover"
                />
              </div>
              <div
                v-else
                class="w-16 h-16 rounded-full bg-primary-100 flex items-center justify-center border-4 border-white shadow-lg"
              >
                <span class="text-2xl text-primary-600 font-bold">
                  {{
                    (authStore.user?.nickname ||
                      authStore.user?.email ||
                      "C")[0].toUpperCase()
                  }}
                </span>
              </div>
            </div>
          </div>

          <div class="mb-8">
            <h3 class="text-xl font-semibold text-gray-800 mb-4">快速操作</h3>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <router-link
                to="/caregiver-profile-setup"
                class="group bg-white rounded-xl shadow-md hover:shadow-xl transition-all p-4 border border-gray-200 hover:border-primary-300"
              >
                <div class="text-center">
                  <div
                    class="bg-blue-100 group-hover:bg-blue-200 rounded-lg p-3 mb-2 transition-colors inline-block"
                  >
                    <i
                      class="fa-solid fa-user-edit text-blue-600 text-xl md:text-2xl"
                    ></i>
                  </div>
                  <h4
                    class="font-semibold text-gray-800 group-hover:text-primary-600 transition-colors mb-1"
                  >
                    我的檔案
                  </h4>
                  <p class="text-xs text-gray-600">管理個人資料</p>
                </div>
              </router-link>

              <router-link
                to="/caregiver-availability"
                class="group bg-white rounded-xl shadow-md hover:shadow-xl transition-all p-4 border border-gray-200 hover:border-primary-300"
              >
                <div class="text-center">
                  <div
                    class="bg-green-100 group-hover:bg-green-200 rounded-lg p-3 mb-2 transition-colors inline-block"
                  >
                    <i
                      class="fa-solid fa-calendar-check text-green-600 text-xl md:text-2xl"
                    ></i>
                  </div>
                  <h4
                    class="font-semibold text-gray-800 group-hover:text-primary-600 transition-colors mb-1"
                  >
                    服務時間
                  </h4>
                  <p class="text-xs text-gray-600">設定可服務時段</p>
                </div>
              </router-link>

              <router-link
                to="/caregiver-licenses"
                class="group bg-white rounded-xl shadow-md hover:shadow-xl transition-all p-4 border border-gray-200 hover:border-primary-300"
              >
                <div class="text-center">
                  <div
                    class="bg-yellow-100 group-hover:bg-yellow-200 rounded-lg p-3 mb-2 transition-colors inline-block"
                  >
                    <i
                      class="fa-solid fa-certificate text-yellow-600 text-xl md:text-2xl"
                    ></i>
                  </div>
                  <h4
                    class="font-semibold text-gray-800 group-hover:text-primary-600 transition-colors mb-1"
                  >
                    證照管理
                  </h4>
                  <p class="text-xs text-gray-600">上傳專業證照</p>
                </div>
              </router-link>

              <router-link
                to="/contracts"
                class="group bg-white rounded-xl shadow-md hover:shadow-xl transition-all p-4 border border-gray-200 hover:border-primary-300"
              >
                <div class="text-center">
                  <div
                    class="bg-purple-100 group-hover:bg-purple-200 rounded-lg p-3 mb-2 transition-colors inline-block"
                  >
                    <i
                      class="fa-solid fa-file-contract text-purple-600 text-xl md:text-2xl"
                    ></i>
                  </div>
                  <h4
                    class="font-semibold text-gray-800 group-hover:text-primary-600 transition-colors mb-1"
                  >
                    我的合約
                  </h4>
                  <p class="text-xs text-gray-600">查看服務合約</p>
                </div>
              </router-link>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div
              class="bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl shadow-lg p-6 text-white"
            >
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-blue-100 text-sm mb-1">進行中的合約</p>
                  <p class="text-3xl font-bold">
                    {{ stats.activeContracts || 0 }}
                  </p>
                </div>
                <div class="bg-white/20 rounded-full p-3">
                  <i class="fa-solid fa-file-contract text-2xl"></i>
                </div>
              </div>
            </div>

            <div
              class="bg-gradient-to-br from-green-500 to-green-600 rounded-xl shadow-lg p-6 text-white"
            >
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-green-100 text-sm mb-1">已完成</p>
                  <p class="text-3xl font-bold">
                    {{ stats.completedContracts || 0 }}
                  </p>
                </div>
                <div class="bg-white/20 rounded-full p-3">
                  <i class="fa-solid fa-check-circle text-2xl"></i>
                </div>
              </div>
            </div>

            <div
              class="bg-gradient-to-br from-yellow-500 to-yellow-600 rounded-xl shadow-lg p-6 text-white"
            >
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-yellow-100 text-sm mb-1">待接受</p>
                  <p class="text-3xl font-bold">
                    {{ stats.pendingContracts || 0 }}
                  </p>
                </div>
                <div class="bg-white/20 rounded-full p-3">
                  <i class="fa-solid fa-clock text-2xl"></i>
                </div>
              </div>
            </div>

            <div
              class="bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl shadow-lg p-6 text-white"
            >
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-purple-100 text-sm mb-1">平均評分</p>
                  <p class="text-3xl font-bold">
                    {{ profileStats.avgRating || "0.0" }}
                  </p>
                </div>
                <div class="bg-white/20 rounded-full p-3">
                  <i class="fa-solid fa-star text-2xl"></i>
                </div>
              </div>
            </div>
          </div>

          <div
            v-if="profileStats.hasProfile"
            class="bg-white rounded-xl shadow-md border border-gray-200 p-6 mb-8"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div class="bg-green-100 rounded-full p-3">
                  <i
                    class="fa-solid fa-check-circle text-green-600 text-xl"
                  ></i>
                </div>
                <div>
                  <h4 class="font-semibold text-gray-800">檔案已完成</h4>
                  <p class="text-sm text-gray-600">
                    {{ profileStats.fullName }} · NT$
                    {{ profileStats.serviceRate }}/小時 · ⭐
                    {{ profileStats.avgRating || "0.0" }}
                  </p>
                </div>
              </div>
              <router-link
                to="/caregiver-profile-setup"
                class="text-primary-600 hover:text-primary-700 text-sm font-medium"
              >
                編輯 →
              </router-link>
            </div>
          </div>
          <div
            v-else
            class="bg-yellow-50 border border-yellow-200 rounded-xl p-6 mb-8"
          >
            <div class="flex items-center space-x-4">
              <div class="bg-yellow-100 rounded-full p-3">
                <i
                  class="fa-solid fa-exclamation-triangle text-yellow-600 text-xl"
                ></i>
              </div>
              <div class="flex-1">
                <h4 class="font-semibold text-gray-800 mb-1">
                  尚未建立照護者檔案
                </h4>
                <p class="text-sm text-gray-600 mb-3">
                  完成檔案設定以開始接受服務邀請
                </p>
                <router-link
                  to="/caregiver-profile-setup"
                  class="btn-primary inline-block"
                >
                  立即建立
                </router-link>
              </div>
            </div>
          </div>

          <div
            class="bg-white rounded-xl shadow-md border border-gray-200 overflow-hidden"
          >
            <div
              class="px-6 py-4 border-b border-gray-200 flex justify-between items-center"
            >
              <h3 class="text-xl font-semibold text-gray-800">最近的合約</h3>
              <router-link
                to="/contracts"
                class="text-sm text-primary-600 hover:text-primary-700 font-medium"
              >
                查看全部 →
              </router-link>
            </div>
            <div v-if="loadingContracts" class="p-8 text-center">
              <div
                class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
              ></div>
              <p class="mt-2 text-gray-600">載入中...</p>
            </div>
            <div
              v-else-if="recentContracts.length === 0"
              class="p-8 text-center text-gray-500"
            >
              <i
                class="fa-solid fa-file-circle-question text-4xl mb-3 opacity-50"
              ></i>
              <p>目前沒有合約</p>
            </div>
            <div v-else class="divide-y divide-gray-200">
              <div
                v-for="contract in recentContracts"
                :key="contract.contract_id"
                class="p-6 hover:bg-gray-50 transition-colors cursor-pointer"
                @click="$router.push(`/contract/${contract.contract_id}`)"
              >
                <div class="flex items-center justify-between">
                  <div class="flex-1">
                    <div class="flex items-center space-x-3 mb-2">
                      <h4 class="font-semibold text-gray-800">
                        合約 #{{ contract.contract_id }}
                      </h4>
                      <span
                        :class="{
                          'bg-yellow-100 text-yellow-800':
                            contract.status === 'pending',
                          'bg-green-100 text-green-800':
                            contract.status === 'active',
                          'bg-blue-100 text-blue-800':
                            contract.status === 'completed',
                          'bg-red-100 text-red-800':
                            contract.status === 'canceled',
                        }"
                        class="px-2 py-1 rounded-full text-xs font-medium"
                      >
                        {{ getStatusText(contract.status) }}
                      </span>
                    </div>
                    <p class="text-sm text-gray-600">
                      {{ getUserName(contract) }} ·
                      {{ formatDate(contract.start_date) }} -
                      {{ formatDate(contract.end_date) }}
                    </p>
                  </div>
                  <i class="fa-solid fa-chevron-right text-gray-400"></i>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { contractAPI, caregiverAPI } from "@/services/api";
import Navbar from "@/components/Navbar.vue";

const router = useRouter();
const authStore = useAuthStore();

const stats = ref({
  activeContracts: 0,
  completedContracts: 0,
  pendingContracts: 0,
});

const profileStats = ref({
  hasProfile: false,
  fullName: "",
  serviceRate: 0,
  avgRating: 0,
});

const recentContracts = ref([]);
const loadingContracts = ref(false);

const formatDate = (dateString) => {
  if (!dateString) return "無資料";
  const date = new Date(dateString);
  return date.toLocaleDateString("zh-TW", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
};

const getStatusText = (status) => {
  const statusMap = {
    pending: "待確認",
    active: "進行中",
    completed: "已完成",
    canceled: "已取消",
  };
  return statusMap[status] || status;
};

const getCaregiverName = (contract) => {
  const caregiver = contract.Caregiver || contract.caregiver || {};
  return caregiver.nickname || caregiver.email || "照護者";
};

const getUserName = (contract) => {
  const user = contract.User || contract.user || {};
  return user.nickname || user.email || "使用者";
};

const loadContracts = async () => {
  loadingContracts.value = true;
  try {
    const contracts = await contractAPI.getContracts();

    // Calculate statistics
    stats.value = {
      activeContracts: contracts.filter((c) => c.status === "active").length,
      completedContracts: contracts.filter((c) => c.status === "completed")
        .length,
      pendingContracts: contracts.filter((c) => c.status === "pending").length,
    };

    // Get recent contracts (last 5, sorted by date)
    recentContracts.value = contracts
      .sort(
        (a, b) =>
          new Date(b.created_at || b.start_date) -
          new Date(a.created_at || a.start_date)
      )
      .slice(0, 5);
  } catch (error) {
    console.error("Load contracts error:", error);
  } finally {
    loadingContracts.value = false;
  }
};

const loadCaregiverProfile = async () => {
  try {
    const profile = await caregiverAPI.getMyProfile();
    if (profile) {
      profileStats.value = {
        hasProfile: true,
        fullName: profile.full_name || "",
        serviceRate: profile.service_rate || 0,
        avgRating: profile.avg_rating || 0,
      };
    }
  } catch (error) {
    // Profile doesn't exist yet
    profileStats.value.hasProfile = false;
  }
};

onMounted(() => {
  // Redirect admin to admin dashboard
  if (authStore.isAdmin) {
    router.push("/admin");
    return;
  }

  // Load data
  loadContracts();
  if (authStore.isCaregiver) {
    loadCaregiverProfile();
  }
});
</script>
