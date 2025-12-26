<template>
  <div>
    <Navbar />
    <div class="bg-gray-50 min-h-screen">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
        <header class="mb-8">
          <h2 class="text-3xl font-bold text-gray-900">服務時間管理</h2>
          <p class="text-gray-500 mt-2">設定您的每週固定行程與特殊日期的排程，讓案主更容易預約。</p>
        </header>

        <!-- 每週固定時間 -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div class="px-6 py-4 bg-blue-50 border-b border-blue-100 flex justify-between items-center">
            <div class="flex items-center gap-3">
              <div class="bg-blue-100 p-2 rounded-full text-blue-600">
                <i class="fa-solid fa-calendar-week"></i>
              </div>
              <div>
                <h3 class="text-lg font-bold text-gray-800">每週固定時間</h3>
                <p class="text-xs text-blue-600/80">例行性的服務時間</p>
              </div>
            </div>
            <button
              @click="addAvailabilitySlot"
              class="text-sm bg-white border border-blue-200 text-blue-600 px-3 py-1.5 rounded-lg hover:bg-blue-50 transition-colors shadow-sm"
            >
              <i class="fa-solid fa-plus mr-1"></i> 新增時段
            </button>
          </div>

          <div class="p-6">
            <div v-if="weeklySlots.length === 0" class="text-center py-8 text-gray-400">
              <i class="fa-regular fa-calendar-xmark text-4xl mb-2 opacity-50"></i>
              <p>尚未設定每週時間</p>
            </div>
            <div v-else class="space-y-3">
              <div
                v-for="(slot, index) in weeklySlots"
                :key="slot.id"
                class="slot-card group flex flex-wrap md:flex-nowrap items-center gap-3 p-3 bg-white border border-gray-200 rounded-lg hover:border-blue-300"
              >
                <div class="w-full md:w-auto md:flex-1 grid grid-cols-2 md:grid-cols-4 gap-3 items-center">
                  <div class="col-span-2 md:col-span-1">
                    <label class="block text-xs text-gray-500 md:hidden mb-1">星期</label>
                    <div class="relative">
                      <i class="fa-solid fa-calendar absolute left-2 top-2.5 text-gray-400 text-sm"></i>
                      <select
                        v-model="slot.day_of_week"
                        class="input-mini pl-7 bg-gray-50 focus:bg-white focus:ring-1 focus:ring-blue-500 w-full"
                        required
                      >
                        <option value="">選擇星期</option>
                        <option :value="1">週一</option>
                        <option :value="2">週二</option>
                        <option :value="3">週三</option>
                        <option :value="4">週四</option>
                        <option :value="5">週五</option>
                        <option :value="6">週六</option>
                        <option :value="7">週日</option>
                      </select>
                    </div>
                  </div>
                  <div class="md:col-span-3 flex items-center gap-2">
                    <div class="flex-1">
                      <label class="block text-xs text-gray-500 md:hidden mb-1">開始</label>
                      <input
                        v-model="slot.start_time"
                        type="time"
                        class="input-mini w-full"
                        required
                      />
                    </div>
                    <span class="text-gray-400 mt-5 md:mt-0">
                      <i class="fa-solid fa-arrow-right text-xs"></i>
                    </span>
                    <div class="flex-1">
                      <label class="block text-xs text-gray-500 md:hidden mb-1">結束</label>
                      <input
                        v-model="slot.end_time"
                        type="time"
                        class="input-mini w-full"
                        required
                      />
                    </div>
                  </div>
                </div>
                <div class="w-full md:w-auto flex justify-end md:border-l md:pl-3 border-gray-100">
                  <button
                    @click="removeSlot(index)"
                    class="text-gray-400 hover:text-red-500 hover:bg-red-50 p-2 rounded-full transition-colors"
                    title="刪除"
                  >
                    <i class="fa-solid fa-trash-can"></i>
                  </button>
                </div>
              </div>
            </div>

            <div v-if="weeklyError" class="text-red-600 text-sm mt-4 bg-red-50 p-3 rounded-md border border-red-100">
              {{ weeklyError }}
            </div>

            <div class="flex justify-end space-x-3 mt-6 pt-4 border-t border-gray-100">
              <button
                @click="resetWeekly"
                class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg text-sm"
              >
                取消變更
              </button>
              <button
                @click="saveAvailability"
                :disabled="weeklyLoading"
                class="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2 rounded-lg text-sm font-medium shadow-sm transition-all"
              >
                {{ weeklyLoading ? '儲存中...' : '儲存設定' }}
              </button>
            </div>
          </div>
        </div>

        <!-- 特殊日期設定 -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div class="px-6 py-4 bg-amber-50 border-b border-amber-100 flex justify-between items-center">
            <div class="flex items-center gap-3">
              <div class="bg-amber-100 p-2 rounded-full text-amber-600">
                <i class="fa-solid fa-calendar-day"></i>
              </div>
              <div>
                <h3 class="text-lg font-bold text-gray-800">特殊日期設定</h3>
                <p class="text-xs text-amber-600/80">優先權高於每週時間（如：國定假日、臨時請假）</p>
              </div>
            </div>
            <button
              @click="addSpecialSlot"
              class="text-sm bg-white border border-amber-200 text-amber-600 px-3 py-1.5 rounded-lg hover:bg-amber-50 transition-colors shadow-sm"
            >
              <i class="fa-solid fa-plus mr-1"></i> 新增日期
            </button>
          </div>

          <div class="p-6">
            <div v-if="specialSlots.length === 0" class="text-center py-8 text-gray-400">
              <i class="fa-solid fa-mug-hot text-4xl mb-2 opacity-50"></i>
              <p>無特殊排程</p>
            </div>
            <div v-else class="space-y-3">
              <div
                v-for="(slot, index) in specialSlots"
                :key="slot.id"
                class="slot-card group flex flex-wrap md:flex-nowrap items-center gap-3 p-3 bg-white border border-gray-200 rounded-lg hover:border-amber-300 border-l-4 border-l-amber-400"
              >
                <div class="w-full md:w-auto md:flex-1 grid grid-cols-2 md:grid-cols-10 gap-3 items-center">
                  <div class="col-span-2 md:col-span-3">
                    <label class="block text-xs text-gray-500 md:hidden mb-1">日期</label>
                    <input
                      v-model="slot.date"
                      type="date"
                      class="input-mini bg-gray-50 focus:bg-white w-full"
                      required
                    />
                  </div>
                  <div class="col-span-2 md:col-span-4 flex items-center gap-2">
                    <div class="flex-1">
                      <label class="block text-xs text-gray-500 md:hidden mb-1">開始</label>
                      <input
                        v-model="slot.start_time"
                        type="time"
                        class="input-mini w-full"
                        required
                      />
                    </div>
                    <span class="text-gray-400 mt-5 md:mt-0">
                      <i class="fa-solid fa-arrow-right text-xs"></i>
                    </span>
                    <div class="flex-1">
                      <label class="block text-xs text-gray-500 md:hidden mb-1">結束</label>
                      <input
                        v-model="slot.end_time"
                        type="time"
                        class="input-mini w-full"
                        required
                      />
                    </div>
                  </div>
                  <div class="col-span-2 md:col-span-3">
                    <label class="block text-xs text-gray-500 md:hidden mb-1">狀態</label>
                    <select
                      v-model="slot.is_available"
                      :class="[
                        'input-mini w-full',
                        slot.is_available ? 'text-green-600 font-medium' : 'text-red-500 font-medium'
                      ]"
                      required
                    >
                      <option :value="true">✅ 可服務</option>
                      <option :value="false">⛔ 不可服務</option>
                    </select>
                  </div>
                </div>
                <div class="w-full md:w-auto flex justify-end md:border-l md:pl-3 border-gray-100">
                  <button
                    @click="removeSpecialSlot(index)"
                    class="text-gray-400 hover:text-red-500 hover:bg-red-50 p-2 rounded-full transition-colors"
                    title="刪除"
                  >
                    <i class="fa-solid fa-trash-can"></i>
                  </button>
                </div>
              </div>
            </div>

            <div v-if="specialError" class="text-red-600 text-sm mt-4 bg-red-50 p-3 rounded-md border border-red-100">
              {{ specialError }}
            </div>

            <div class="flex justify-end space-x-3 mt-6 pt-4 border-t border-gray-100">
              <button
                @click="resetSpecial"
                class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg text-sm"
              >
                取消變更
              </button>
              <button
                @click="saveSpecialAvailability"
                :disabled="specialLoading"
                class="bg-amber-500 hover:bg-amber-600 text-white px-5 py-2 rounded-lg text-sm font-medium shadow-sm transition-all"
              >
                {{ specialLoading ? '儲存中...' : '儲存特殊時段' }}
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
import { caregiverAPI } from '@/services/api'
import Navbar from '@/components/Navbar.vue'

const weeklySlots = ref([])
const specialSlots = ref([])
const weeklyError = ref('')
const specialError = ref('')
const weeklyLoading = ref(false)
const specialLoading = ref(false)
let slotIdCounter = 0

const addAvailabilitySlot = () => {
  weeklySlots.value.push({
    id: slotIdCounter++,
    day_of_week: '',
    start_time: '',
    end_time: '',
  })
}

const removeSlot = (index) => {
  weeklySlots.value.splice(index, 1)
}

const addSpecialSlot = () => {
  specialSlots.value.push({
    id: slotIdCounter++,
    date: '',
    start_time: '',
    end_time: '',
    is_available: true,
  })
}

const removeSpecialSlot = (index) => {
  specialSlots.value.splice(index, 1)
}

const saveAvailability = async () => {
  weeklyError.value = ''
  
  const availabilities = weeklySlots.value
    .filter(slot => slot.day_of_week && slot.start_time && slot.end_time)
    .map(slot => {
      if (slot.start_time >= slot.end_time) {
        throw new Error('開始時間必須早於結束時間')
      }
      return {
        day_of_week: parseInt(slot.day_of_week),
        start_time: slot.start_time,
        end_time: slot.end_time,
      }
    })

  weeklyLoading.value = true
  try {
    await caregiverAPI.updateAvailability(availabilities)
    await loadAvailability()
    alert('每週服務時間已更新！')
  } catch (error) {
    weeklyError.value = error.response?.data?.error || error.message || '儲存失敗'
  } finally {
    weeklyLoading.value = false
  }
}

const saveSpecialAvailability = async () => {
  specialError.value = ''
  
  const specialAvailabilities = specialSlots.value
    .filter(slot => slot.date && slot.start_time && slot.end_time)
    .map(slot => {
      if (slot.start_time >= slot.end_time) {
        throw new Error('開始時間必須早於結束時間')
      }
      return {
        date: slot.date,
        start_time: slot.start_time,
        end_time: slot.end_time,
        is_available: slot.is_available,
      }
    })

  specialLoading.value = true
  try {
    await caregiverAPI.updateSpecialAvailability(specialAvailabilities)
    await loadSpecialAvailability()
    alert('特殊時段已更新！')
  } catch (error) {
    specialError.value = error.response?.data?.error || error.message || '儲存失敗'
  } finally {
    specialLoading.value = false
  }
}

const loadAvailability = async () => {
  try {
    const response = await caregiverAPI.getAvailability()
    const availabilities = response.availabilities || response || []
    weeklySlots.value = availabilities.map(avail => ({
      id: slotIdCounter++,
      day_of_week: avail.day_of_week || avail.day_of_week,
      start_time: avail.start_time || '',
      end_time: avail.end_time || '',
    }))
  } catch (error) {
    console.error('載入失敗:', error)
  }
}

const loadSpecialAvailability = async () => {
  try {
    const response = await caregiverAPI.getSpecialAvailability()
    const specialAvailabilities = response.special_availabilities || response || []
    specialSlots.value = specialAvailabilities.map(avail => {
      const dateStr = avail.date ? (typeof avail.date === 'string' ? avail.date.split('T')[0] : avail.date) : ''
      return {
        id: slotIdCounter++,
        date: dateStr,
        start_time: avail.start_time || '',
        end_time: avail.end_time || '',
        is_available: avail.is_available !== undefined ? avail.is_available : true,
      }
    })
  } catch (error) {
    console.error('載入特殊時段失敗:', error)
  }
}

const resetWeekly = () => {
  loadAvailability()
}

const resetSpecial = () => {
  loadSpecialAvailability()
}

onMounted(() => {
  loadAvailability()
  loadSpecialAvailability()
})
</script>

<style scoped>
.input-mini {
  @apply px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.slot-card {
  transition: all 0.2s;
}
</style>