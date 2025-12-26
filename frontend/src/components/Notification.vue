<template>
  <teleport to="body">
    <transition name="notification">
      <div
        v-if="visible"
        :class="[
          'fixed top-4 right-4 p-4 rounded-lg shadow-lg z-50 text-white',
          type === 'error' ? 'bg-red-500' : type === 'success' ? 'bg-green-500' : 'bg-blue-500'
        ]"
      >
        {{ message }}
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'

const props = defineProps({
  message: {
    type: String,
    required: true,
  },
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['info', 'success', 'error'].includes(value),
  },
  duration: {
    type: Number,
    default: 3000,
  },
})

const visible = ref(false)
let timeout = null

onMounted(() => {
  visible.value = true
  if (props.duration > 0) {
    timeout = setTimeout(() => {
      visible.value = false
    }, props.duration)
  }
})

watch(() => props.message, () => {
  visible.value = true
  if (timeout) clearTimeout(timeout)
  if (props.duration > 0) {
    timeout = setTimeout(() => {
      visible.value = false
    }, props.duration)
  }
})
</script>

<style scoped>
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>