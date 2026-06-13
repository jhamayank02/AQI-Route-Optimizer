<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useDebounceFn, useEventListener } from '@vueuse/core'
import { searchLocation } from '../services/location.service';
import type { Location } from '../services/location.service';
import {useStore} from '../store/store'

const store = useStore()

const emit = defineEmits<{
  findRoutes: []
}>()

const source = ref('')
const destination = ref('')

const sourceOptions = ref<Location[]>([])
const destinationOptions = ref<Location[]>([])
const activeDropdown = ref<'source' | 'destination' | null>(null)
const searchBoxRef = ref<HTMLElement | null>(null)
const selectedSource = ref<Location | null>(null)
const selectedDestination = ref<Location | null>(null)

const debounceSearch = useDebounceFn(
  (type: 'source' | 'destination', value: string)=>{
    searchLocationHandler(type, value)
  },
  500
)

const hasSource = computed(() => source.value.trim().length > 0)
const canFindRoutes = computed(() => hasSource.value && destination.value.trim().length > 0)
const shouldShowSourceOptions = computed(() => activeDropdown.value === 'source' && sourceOptions.value.length > 0)
const shouldShowDestinationOptions = computed(
  () => activeDropdown.value === 'destination' && destinationOptions.value.length > 0
)

// If source input is empty, clear destination input
watch(source, (value) => {
  if (!value.trim()) {
    destination.value = ''
    sourceOptions.value = []
    destinationOptions.value = []
    activeDropdown.value = null
    selectedSource.value = null
    selectedDestination.value = null
    return
  }

  if (selectedSource.value?.label !== value) {
    selectedSource.value = null
  }
})

watch(destination, (value) => {
  if (!value.trim()) {
    destinationOptions.value = []
    selectedDestination.value = null
    if (activeDropdown.value === 'destination') {
      activeDropdown.value = null
    }
    return
  }

  if (selectedDestination.value?.label !== value) {
    selectedDestination.value = null
  }
})

// Handler to swap locations (source <-> destination and vice versa)
const swapLocations = () => {
  if (!canFindRoutes.value) {
    return
  }

  const currentSource = source.value
  source.value = destination.value
  destination.value = currentSource
}

const findRoutes = () => {
  if (!canFindRoutes.value) {
    return
  }

  store.setSource(selectedSource.value as Location)
  store.setDestination(selectedDestination.value as Location)

  closeDropdowns()
  emit('findRoutes')
}

const closeDropdowns = () => {
  activeDropdown.value = null
}

const selectLocation = (type: 'source' | 'destination', option: Location) => {
  if (type === 'source') {
    source.value = option.label
    sourceOptions.value = []
    selectedSource.value = option
  } else {
    destination.value = option.label
    destinationOptions.value = []
    selectedDestination.value = option
  }

  activeDropdown.value = null
}

const validateSelection = (type: 'source' | 'destination') => {
  window.setTimeout(() => {
    if (type === 'source' && source.value.trim() && !selectedSource.value) {
      source.value = ''
      sourceOptions.value = []
    }

    if (type === 'destination' && destination.value.trim() && !selectedDestination.value) {
      destination.value = ''
      destinationOptions.value = []
    }
  }, 150)
}

const searchLocationHandler = async (type: 'source' | 'destination', value: string) => {
  const query = value.trim()

  if (!query) {
    if (type === 'source') {
      sourceOptions.value = []
    } else {
      destinationOptions.value = []
    }
    activeDropdown.value = null
    return
  }

  try {
    const response = await searchLocation(query)
    const options = response.data || []

    if (type === 'source') {
      sourceOptions.value = options
    } else {
      destinationOptions.value = options
    }

    activeDropdown.value = options.length > 0 ? type : null
  } catch(error) {
    console.log(error)
  }
}

useEventListener(document, 'click', (event) => {
  if (!searchBoxRef.value?.contains(event.target as Node)) {
    closeDropdowns()
  }
})
</script>

<template>
  <main class="flex flex-1 flex-col items-center px-6 pb-8 pt-9">
    <div class="text-center">
      <h1 class="text-[28px] font-extrabold leading-tight text-black sm:text-[34px]">
        Find a <span class="text-[#2f9e52]">healthier</span> route
      </h1>
      <p class="mt-4 max-w-md text-[16px] leading-7 text-slate-800">
        We analyze air quality along multiple routes and suggest the best one for you.
      </p>
    </div>

    <form
      ref="searchBoxRef"
      class="relative mt-8 w-full max-w-[610px] rounded-lg border border-slate-200 bg-white px-5 py-5 shadow-[0_8px_26px_rgba(15,23,42,0.08)]"
      @submit.prevent="findRoutes"
    >
      <div class="space-y-5">
        <div class="relative">
          <label for="source" class="mb-2 block text-sm font-bold text-slate-950">From</label>
          <div
            class="flex h-14 items-center gap-3 rounded-md border border-slate-200 bg-white px-4 transition focus-within:border-[#2f9e52] focus-within:shadow-[0_0_0_3px_rgba(47,158,82,0.12)]"
          >
            <input
              id="source"
              v-model="source"
              type="text"
              placeholder="Enter starting location"
              class="min-w-0 flex-1 border-0 bg-transparent p-0 text-[16px] font-medium text-slate-900 outline-none placeholder:text-slate-500 focus:ring-0"
              @focus="activeDropdown = sourceOptions.length ? 'source' : null"
              @blur="validateSelection('source')"
              @input="debounceSearch('source', source)"
            />
            <svg
              viewBox="0 0 24 24"
              width="24"
              height="24"
              class="h-6 w-6 shrink-0 text-[#2f9e52]"
              fill="none"
              aria-hidden="true"
            >
              <circle cx="12" cy="12" r="8" stroke="currentColor" stroke-width="2" />
              <circle cx="12" cy="12" r="2.5" fill="currentColor" />
            </svg>
          </div>

          <div
            v-if="shouldShowSourceOptions"
            class="absolute z-20 mt-2 max-h-72 w-full overflow-y-auto overflow-x-hidden rounded-md border border-slate-200 bg-white shadow-[0_12px_30px_rgba(15,23,42,0.12)]"
          >
            <button
              v-for="option in sourceOptions"
              :key="`source-${option.lat}-${option.lng}`"
              type="button"
              class="block w-full border-b border-slate-100 px-4 py-3 text-left transition last:border-b-0 hover:bg-slate-50"
              @click="selectLocation('source', option)"
            >
              <span class="block text-sm font-semibold text-slate-900">{{ option.name }}</span>
              <span class="block text-xs text-slate-500">{{ option.label }}</span>
            </button>
          </div>
        </div>

        <div class="relative">
          <label for="destination" class="mb-2 block text-sm font-bold text-slate-950">To</label>
          <div
            class="flex h-14 items-center gap-3 rounded-md border bg-white px-4 transition"
            :class="
              hasSource
                ? 'border-slate-200 focus-within:border-red-400 focus-within:shadow-[0_0_0_3px_rgba(239,68,68,0.10)]'
                : 'border-slate-200 bg-slate-50 opacity-75'
            "
          >
            <input
              id="destination"
              v-model="destination"
              type="text"
              placeholder="Enter destination"
              :disabled="!hasSource"
              class="min-w-0 flex-1 border-0 bg-transparent p-0 text-[16px] font-medium text-slate-900 outline-none placeholder:text-slate-500 disabled:cursor-not-allowed disabled:text-slate-400 focus:ring-0"
              @focus="activeDropdown = destinationOptions.length ? 'destination' : null"
              @blur="validateSelection('destination')"
              @input="debounceSearch('destination', destination)"
            />
            <svg
              viewBox="0 0 24 24"
              width="24"
              height="24"
              class="h-6 w-6 shrink-0 text-red-500"
              fill="none"
              aria-hidden="true"
            >
              <path
                d="M12 21s7-5.1 7-11a7 7 0 1 0-14 0c0 5.9 7 11 7 11Z"
                stroke="currentColor"
                stroke-width="2"
              />
              <circle cx="12" cy="10" r="2.4" fill="currentColor" />
            </svg>
          </div>

          <div
            v-if="shouldShowDestinationOptions"
            class="absolute z-20 mt-2 max-h-72 w-full overflow-y-auto overflow-x-hidden rounded-md border border-slate-200 bg-white shadow-[0_12px_30px_rgba(15,23,42,0.12)]"
          >
            <button
              v-for="option in destinationOptions"
              :key="`destination-${option.lat}-${option.lng}`"
              type="button"
              class="block w-full border-b border-slate-100 px-4 py-3 text-left transition last:border-b-0 hover:bg-slate-50"
              @click="selectLocation('destination', option)"
            >
              <span class="block text-sm font-semibold text-slate-900">{{ option.name }}</span>
              <span class="block text-xs text-slate-500">{{ option.label }}</span>
            </button>
          </div>
        </div>
      </div>

      <button
        type="button"
        class="absolute right-[-45px] top-1/2 hidden h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full border border-slate-200 bg-white text-slate-700 shadow-[0_5px_16px_rgba(15,23,42,0.14)] transition hover:border-[#2f9e52] hover:text-[#2f9e52] disabled:cursor-not-allowed disabled:opacity-45 sm:flex"
        :disabled="!canFindRoutes"
        aria-label="Swap source and destination"
        @click="swapLocations"
      >
        <svg viewBox="0 0 24 24" width="24" height="24" class="h-6 w-6" fill="none" aria-hidden="true">
          <path
            d="M8 4v14m0 0-4-4m4 4 4-4M16 20V6m0 0-4 4m4-4 4 4"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
          />
        </svg>
      </button>

      <button
        type="submit"
        class="mt-6 h-14 w-full rounded-md bg-[#2f9e52] text-[16px] font-bold text-white shadow-[0_8px_20px_rgba(47,158,82,0.22)] transition hover:bg-[#278846] disabled:cursor-not-allowed disabled:bg-slate-300 disabled:shadow-none"
        :disabled="!canFindRoutes"
      >
        Find Routes
      </button>
    </form>

    <div
      class="mt-auto grid w-full max-w-[860px] grid-cols-1 gap-4 border-t border-slate-200 pt-6 text-sm font-medium text-slate-700 sm:grid-cols-3 sm:divide-x sm:divide-slate-200"
    >
      <div class="flex items-center justify-center gap-3">
        <svg viewBox="0 0 24 24" width="24" height="24" class="h-6 w-6 text-[#2f9e52]" fill="none" aria-hidden="true">
          <path d="M19 4c-6.2.4-11.1 3.5-13.2 8.1-1.1 2.5-1 5.1.2 6.5 1.4 1.2 4 1.3 6.5.2 4.6-2.1 7.7-7 8.1-13.2.1-1-.6-1.7-1.6-1.6Z" stroke="currentColor" stroke-width="2" />
          <path d="M5 19 19 5" stroke="currentColor" stroke-linecap="round" stroke-width="2" />
        </svg>
        <span>Real-time AQI Data</span>
      </div>
      <div class="flex items-center justify-center gap-3 sm:px-6">
        <svg viewBox="0 0 24 24" width="24" height="24" class="h-6 w-6 text-[#2f9e52]" fill="none" aria-hidden="true">
          <path d="M4 8c2.5-3 5-3 7.5 0s5 3 8.5 0M4 16c2.5-3 5-3 7.5 0s5 3 8.5 0" stroke="currentColor" stroke-linecap="round" stroke-width="2" />
        </svg>
        <span>Multiple Route Options</span>
      </div>
      <div class="flex items-center justify-center gap-3">
        <svg viewBox="0 0 24 24" width="24" height="24" class="h-6 w-6 text-[#2f9e52]" fill="none" aria-hidden="true">
          <path d="M20.3 5.7a5.2 5.2 0 0 0-7.4 0l-.9.9-.9-.9a5.2 5.2 0 1 0-7.4 7.4l8.3 8.3 8.3-8.3a5.2 5.2 0 0 0 0-7.4Z" stroke="currentColor" stroke-linejoin="round" stroke-width="2" />
        </svg>
        <span>Healthier Choices</span>
      </div>
    </div>
  </main>
</template>
