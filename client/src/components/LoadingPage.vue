<script setup lang="ts">
import { useStore } from '../store/store'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

const emit = defineEmits<{
  complete: []
}>()

const store = useStore()
const locations = computed(() => store.getLocation)
const source = computed(() => locations.value.source)
const destination = computed(() => locations.value.destination)
const isLoaded = computed(() => Boolean(store.getRecommendedRoute))

const progress = ref(18)
let timer: number | undefined
const waitingProgress = 84

const stopTimer = () => {
  if (timer) {
    window.clearInterval(timer)
    timer = undefined
  }
}

const startProgress = (target: number, step: number, delay: number, onComplete?: () => void) => {
  stopTimer()
  timer = window.setInterval(() => {
    progress.value = Math.min(progress.value + step, target)

    if (progress.value >= target) {
      stopTimer()
      onComplete?.()
    }
  }, delay)
}

const completeLoading = () => {
  if (progress.value >= 100) {
    emit('complete')
    return
  }

  startProgress(100, 4, 60, () => {
    window.setTimeout(() => emit('complete'), 150)
  })
}

onMounted(() => {
  progress.value = 18
  startProgress(waitingProgress, 2, 140)

  if (isLoaded.value) {
    completeLoading()
  }
})

watch(
  isLoaded,
  (isLoaded) => {
    if (isLoaded) {
      completeLoading()
    }
  }
)

onUnmounted(() => {
  stopTimer()
})
</script>

<template>
  <main class="flex flex-1 flex-col items-center justify-center px-6 py-10 text-center">
    <div
      class="relative h-36 w-36 overflow-hidden rounded-full border border-slate-200 bg-[linear-gradient(135deg,#edf2f7_0%,#f8fafc_42%,#dbeafe_43%,#f8fafc_48%,#ecfdf5_100%)] shadow-[0_8px_26px_rgba(15,23,42,0.08)]"
    >
      <svg viewBox="0 0 144 144" width="144" height="144" class="h-full w-full" fill="none" aria-hidden="true">
        <path d="M20 104 48 78 70 92 98 66 124 82" stroke="#2f9e52" stroke-linecap="round" stroke-linejoin="round" stroke-width="5" />
        <circle cx="28" cy="96" r="5" fill="#2f9e52" />
        <circle cx="48" cy="78" r="5" fill="#2f9e52" />
        <circle cx="70" cy="92" r="5" fill="#2f9e52" />
        <circle cx="98" cy="66" r="5" fill="#2f9e52" />
        <circle cx="114" cy="76" r="18" fill="#2f9e52" />
        <path d="M107 76h14m0 0-6-6m6 6-6 6" stroke="white" stroke-linecap="round" stroke-linejoin="round" stroke-width="3" />
      </svg>
    </div>

    <h1 class="mt-7 text-[28px] font-extrabold leading-tight text-black sm:text-[34px]">
      Finding the best routes for you...
    </h1>
    <p class="mt-4 max-w-md text-[16px] leading-7 text-slate-700">
      Comparing AQI, travel time, and distance between {{ source?.label || 'your source' }} and {{ destination?.label || 'your destination' }}.
    </p>

    <div class="mt-10 flex w-full max-w-[520px] items-center gap-5">
      <div class="h-2 flex-1 overflow-hidden rounded-full bg-slate-200">
        <div
          class="h-full rounded-full bg-[#2f9e52] transition-all duration-150"
          :style="{ width: `${progress}%` }"
        />
      </div>
      <span class="w-11 text-left text-sm font-bold text-slate-700">{{ progress }}%</span>
    </div>
  </main>
</template>
