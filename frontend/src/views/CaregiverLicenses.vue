<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
          <div>
            <h2 class="text-3xl font-bold text-gray-900">專業證照管理</h2>
            <p class="text-gray-500 mt-1">上傳並管理您的專業資格證明，提升案主信任度。</p>
          </div>
          <button
            @click="showUploadModal = true"
            class="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-lg shadow-md hover:shadow-lg transition-all flex items-center gap-2 font-medium"
          >
            <i class="fa-solid fa-plus"></i> 上傳新證照
          </button>
        </div>

        <div v-if="loading" class="bg-white rounded-xl p-10 text-center shadow-sm border border-gray-100">
          <div class="inline-block animate-spin rounded-full h-10 w-10 border-4 border-gray-200 border-t-blue-600"></div>
          <p class="mt-4 text-gray-500 font-medium">資料載入中...</p>
        </div>
        <div v-else-if="licenses.length === 0" class="bg-white rounded-xl p-10 text-center shadow-sm border border-gray-100">
          <i class="fa-solid fa-file-circle-question text-4xl text-gray-400 mb-4"></i>
          <p class="text-gray-600">尚無證照</p>
          <button
            @click="showUploadModal = true"
            class="mt-4 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
          >
            立即上傳
          </button>
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="license in licenses"
            :key="license.license_id"
            class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden"
          >
            <div class="p-6">
              <div class="flex justify-between items-start">
                <div class="flex-1">
                  <div class="flex items-center space-x-3 mb-3">
                    <h3 class="text-lg font-semibold text-gray-800">{{ license.name }}</h3>
                    <span
                      :class="{
                        'bg-yellow-100 text-yellow-800': license.status === 'pending',
                        'bg-green-100 text-green-800': license.status === 'approved',
                        'bg-red-100 text-red-800': license.status === 'rejected',
                      }"
                      class="px-3 py-1 rounded-full text-sm font-medium"
                    >
                      {{ getStatusText(license.status) }}
                    </span>
                  </div>
                  <div class="text-gray-600 space-y-1 text-sm">
                    <p><span class="font-medium">發證日期：</span>{{ formatDate(license.issue_date) }}</p>
                    <p><span class="font-medium">有效期限：</span>{{ formatDate(license.expiry_date) }}</p>
                    <p v-if="license.proof_url">
                      <span class="font-medium">證明文件：</span>
                      <button
                        @click="viewLicenseImage(license.proof_url)"
                        class="text-primary-600 hover:underline"
                      >
                        查看
                      </button>
                    </p>
                  </div>
                </div>
                <button
                  @click="handleDelete(license.license_id)"
                  class="text-gray-400 hover:text-red-500 hover:bg-red-50 p-2 rounded-full transition-colors ml-4"
                  title="刪除"
                >
                  <i class="fa-solid fa-trash-can"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 上傳 Modal -->
      <div
        v-if="showUploadModal"
        class="fixed inset-0 z-50 overflow-y-auto"
        @click.self="closeModal"
      >
        <div class="fixed inset-0 bg-gray-600 bg-opacity-50 transition-opacity"></div>

        <div class="flex min-h-full items-center justify-center p-4 text-center sm:p-0">
          <div
            class="relative transform overflow-hidden rounded-2xl bg-white text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg"
            @click.stop
          >
            <div class="bg-gray-50 px-4 py-3 sm:px-6 border-b border-gray-100 flex justify-between items-center">
              <h3 class="text-lg font-semibold leading-6 text-gray-900 flex items-center gap-2">
                <i class="fa-solid fa-file-circle-plus text-blue-600"></i> 上傳證照
              </h3>
              <button
                @click="closeModal"
                class="text-gray-400 hover:text-gray-600 transition-colors"
              >
                <i class="fa-solid fa-xmark text-xl"></i>
              </button>
            </div>

            <div class="px-4 py-5 sm:p-6">
              <form @submit.prevent="handleUpload" class="space-y-5">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1.5">
                    證照名稱 <span class="text-red-500">*</span>
                  </label>
                  <div class="relative">
                    <span class="absolute left-3 top-2.5 text-gray-400">
                      <i class="fa-solid fa-id-card"></i>
                    </span>
                    <input
                      v-model="uploadData.name"
                      type="text"
                      required
                      class="input-field pl-10"
                      placeholder="例如：照顧服務員技術士證"
                    />
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1.5">
                      發證日期 <span class="text-red-500">*</span>
                    </label>
                    <input
                      v-model="uploadData.issue_date"
                      type="date"
                      required
                      class="input-field text-gray-600"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1.5">
                      有效期限 <span class="text-red-500">*</span>
                    </label>
                    <input
                      v-model="uploadData.expiry_date"
                      type="date"
                      required
                      class="input-field text-gray-600"
                    />
                  </div>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1.5">證明文件</label>
                  <div
                    class="relative border border-dashed border-gray-300 rounded-lg p-4 bg-gray-50 hover:bg-gray-100 transition-colors text-center cursor-pointer"
                    @click="() => fileInputRef?.click()"
                  >
                    <input
                      ref="fileInputRef"
                      type="file"
                      accept="image/*,.pdf"
                      class="hidden"
                      @change="handleFileSelect"
                    />
                    <div class="text-gray-500">
                      <i class="fa-solid fa-cloud-arrow-up text-2xl mb-2 text-gray-400"></i>
                      <p class="text-sm font-medium">點擊或拖曳檔案至此</p>
                      <p class="text-xs text-gray-400 mt-1">支援圖片 (JPG, PNG) 或 PDF 格式</p>
                    </div>
                  </div>
                  <div v-if="selectedFile" class="text-xs text-blue-600 mt-2 text-left">
                    <i class="fa-solid fa-check mr-1"></i> 已選擇: {{ selectedFile.name }}
                  </div>
                </div>

                <div v-if="uploadError" class="text-red-600 text-sm bg-red-50 p-3 rounded-md border border-red-100 flex items-center gap-2">
                  <i class="fa-solid fa-circle-exclamation"></i>
                  <span>{{ uploadError }}</span>
                </div>
              </form>
            </div>

            <div class="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6 border-t border-gray-100">
              <button
                @click="handleUpload"
                :disabled="uploading"
                class="w-full inline-flex justify-center rounded-lg bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 sm:ml-3 sm:w-auto transition-colors disabled:opacity-50"
              >
                {{ uploading ? '上傳中...' : '確認上傳' }}
              </button>
              <button
                @click="closeModal"
                class="mt-3 inline-flex w-full justify-center rounded-lg bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto transition-colors"
              >
                取消
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { licenseAPI } from '@/services/api'
import Navbar from '@/components/Navbar.vue'

const fileInputRef = ref(null)

const licenses = ref([])
const loading = ref(false)
const showUploadModal = ref(false)
const uploadData = ref({
  name: '',
  issue_date: '',
  expiry_date: '',
})
const selectedFile = ref(null)
const uploadError = ref('')
const uploading = ref(false)

const formatDate = (dateString) => {
  if (!dateString) return '無資料'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-TW')
}

const getStatusText = (status) => {
  const statusMap = {
    pending: '審核中',
    approved: '已核准',
    rejected: '已拒絕',
  }
  return statusMap[status] || status
}

const handleFileSelect = (event) => {
  selectedFile.value = event.target.files[0]
}

const closeModal = () => {
  showUploadModal.value = false
  uploadData.value = { name: '', issue_date: '', expiry_date: '' }
  selectedFile.value = null
  uploadError.value = ''
}

const handleUpload = async () => {
  if (!selectedFile.value) {
    uploadError.value = '請選擇證明文件'
    return
  }

  uploading.value = true
  uploadError.value = ''

  try {
    const formData = new FormData()
    formData.append('name', uploadData.value.name)
    formData.append('issue_date', uploadData.value.issue_date)
    formData.append('expiry_date', uploadData.value.expiry_date)
    formData.append('proof', selectedFile.value)

    await licenseAPI.uploadLicense(formData)
    closeModal()
    loadLicenses()
    alert('上傳成功')
  } catch (error) {
    uploadError.value = error.response?.data?.error || error.message || '上傳失敗'
  } finally {
    uploading.value = false
  }
}

const handleDelete = async (id) => {
  if (!confirm('確定要刪除此證照嗎？')) return

  try {
    await licenseAPI.deleteLicense(id)
    loadLicenses()
    alert('刪除成功')
  } catch (error) {
    alert('刪除失敗：' + (error.response?.data?.error || error.message))
  }
}

const viewLicenseImage = async (proofUrl) => {
  if (!proofUrl) return

  const token = localStorage.getItem('token')
  if (!token) {
    alert('請先登入')
    return
  }

  try {
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    let imageUrl = proofUrl
    
    if (proofUrl.startsWith('/uploads/licenses/')) {
      const filename = proofUrl.replace('/uploads/licenses/', '')
      imageUrl = `${API_BASE_URL}/caregivers/licenses/image/${filename}`
    } else if (!proofUrl.startsWith('http') && !proofUrl.startsWith(API_BASE_URL)) {
      imageUrl = `${API_BASE_URL}${proofUrl}`
    }

    // 使用 fetch 並添加 Authorization header
    const response = await fetch(imageUrl, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      throw new Error('無法載入圖片')
    }

    const blob = await response.blob()
    const blobUrl = URL.createObjectURL(blob)
    const newWindow = window.open()
    
    if (newWindow) {
      newWindow.document.write(`
        <html>
          <head><title>證照預覽</title></head>
          <body style="margin:0;padding:0;background:#222;display:flex;justify-content:center;align-items:center;height:100vh;">
            <img src="${blobUrl}" style="max-width:95%;max-height:95vh;box-shadow:0 0 20px rgba(0,0,0,0.5);" />
          </body>
        </html>
      `)
    }
  } catch (error) {
    console.error('View license image error:', error)
    alert('無法載入圖片：' + (error.message || '未知錯誤'))
  }
}

const loadLicenses = async () => {
  loading.value = true
  try {
    licenses.value = await licenseAPI.getMyLicenses()
  } catch (error) {
    console.error('Load licenses error:', error)
    alert('載入失敗：' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadLicenses()
})
</script>