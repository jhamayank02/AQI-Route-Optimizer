<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { LCircleMarker, LMap, LPolyline, LPopup, LTileLayer } from '@vue-leaflet/vue-leaflet'
import {useStore} from '../store/store'
import type { RouteAnalysis } from '../services/route.service'
import { capitalizeFirstLetter } from '../utils/utils'
import { useToast } from 'vue-toastification'
import { useRouter } from 'vue-router'

const store = useStore()
const toast = useToast()
const router = useRouter()
const locations = computed(() => store.getLocation)
const recommendedRoute = computed(() => store.getRecommendedRoute)
const source = computed(() => locations.value.source)
const destination = computed(() => locations.value.destination)

const getAQIClass = (aqi: number) => {
  if (aqi < 80) return 'bg-[#2f9e52]'
  if (aqi < 120) return 'bg-orange-400'
  return 'bg-red-500'
}

const getRouteStyle = (route: RouteAnalysis) => {
  if (route.is_selected) {
    return {
      border: 'border-[#2f9e52]',
      color: 'text-[#2f9e52]'
    }
  }

  if (route.average_aqi < 80) {
    return {
      border: 'border-emerald-200',
      color: 'text-[#2f9e52]'
    }
  }

  if (route.average_aqi < 120) {
    return {
      border: 'border-orange-200',
      color: 'text-orange-500'
    }
  }

  return {
    border: 'border-red-200',
    color: 'text-red-500'
  }
}

const formatDistance = (distanceKm: number) => `${distanceKm.toFixed(1)} km`
const formatDuration = (durationMinutes: number) => `${Math.round(durationMinutes)} min`
const formatAQI = (aqi: number) => Math.round(aqi)

const routes = computed(() =>
  (recommendedRoute.value?.all_routes || []).map((route, index) => {
    const style = getRouteStyle(route)

    return {
      index,
      label: `Route ${index + 1}`,
      isSelected: route.is_selected,
      averageAQI: formatAQI(route.average_aqi),
      maxAQI: formatAQI(route.max_aqi),
      duration: formatDuration(route.route.duration_minutes),
      distance: formatDistance(route.route.distance_km),
      recommendation: route.recommendation,
      samplingStrategy: route.sampling_strategy,
      selectionReason: route.selection_reason,
      border: style.border,
      color: style.color,
      aqiClass: getAQIClass(route.average_aqi)
    }
  })
)

type LatLngTuple = [number, number]

const routeColors = ['#2f9e52', '#fb923c', '#ef4444', '#2563eb', '#7c3aed']

const defaultCenter: LatLngTuple = [20.5937, 78.9629]

const mapCenter = computed<LatLngTuple>(() => {
  if (source.value) {
    return [source.value.lat, source.value.lng]
  }

  return defaultCenter
})

const mapBounds = computed<LatLngTuple[] | undefined>(() => {
  const routeCoordinates =
    recommendedRoute.value?.all_routes.flatMap((route) =>
      route.route.coordinates.map((coordinate) => [coordinate.lat, coordinate.lng] as LatLngTuple)
    ) || []

  const endpointCoordinates = [source.value, destination.value]
    .filter(Boolean)
    .map((location) => [location!.lat, location!.lng] as LatLngTuple)

  const bounds = [...routeCoordinates, ...endpointCoordinates]

  return bounds.length > 1 ? bounds : undefined
})

const routePolylines = computed(() =>
  (recommendedRoute.value?.all_routes || [])
    .map((route, index) => ({
      index,
      label: `Route ${index + 1}`,
      isSelected: route.is_selected,
      color: routeColors[index % routeColors.length],
      latLngs: route.route.coordinates.map((coordinate) => [coordinate.lat, coordinate.lng] as LatLngTuple)
    }))
    .filter((route) => route.latLngs.length > 1)
)

const selectedAqiSamples = computed(() =>
  recommendedRoute.value?.selected_route.aqi_samples || []
)

const selectedRouteLabel = computed(() => {
  const selectedIndex = recommendedRoute.value?.selected_route_index
  return typeof selectedIndex === 'number' ? `Route ${selectedIndex + 1}` : 'a route'
})

const selectedRouteNote = computed(() =>
  recommendedRoute.value?.selected_route.selection_reason ||
  recommendedRoute.value?.selected_route.recommendation ||
  'This route has the best balance for the selected journey.'
)

const emit = defineEmits<{
  newSearch: []
}>()

onMounted(() => {
  const hasRecommendedRoute =
    Boolean(recommendedRoute.value?.selected_route) &&
    (recommendedRoute.value?.all_routes?.length || 0) > 0

  if (hasRecommendedRoute) {
    return
  }

  store.clearRecommendedRoute()
  toast.error('No recommended route found. Please search again.')
  emit('newSearch')
  router.replace({ name: 'home' })
})

</script>

<template>
  <main class="grid min-h-full grid-cols-1 lg:grid-cols-[320px_1fr]">
    <aside class="border-b border-slate-200 bg-white px-6 py-6 lg:border-b-0 lg:border-r">
      <div class="flex items-start justify-between gap-4">
        <div>
          <h1 class="text-[24px] font-extrabold leading-tight text-black">Recommended Routes</h1>
          <p class="mt-2 text-sm leading-6 text-slate-600">From {{ source?.label }} to {{ destination?.label }}</p>
        </div>
        <button
          type="button"
          class="rounded-md border border-[#2f9e52] px-3 py-2 text-xs font-bold text-[#2f9e52] transition hover:bg-emerald-50 lg:hidden"
          @click="emit('newSearch')"
        >
          New
        </button>
      </div>

      <div class="mt-6 space-y-4">
        <article
          v-for="route in routes"
          :key="route.index"
          class="rounded-lg border bg-white p-4 shadow-[0_6px_18px_rgba(15,23,42,0.06)]"
          :class="route.border"
        >
          <div class="flex items-center justify-between gap-3">
            <h2 class="text-base font-extrabold" :class="route.color">
              {{ route.label }} <span v-if="route.isSelected">(Selected)</span>
            </h2>
            <span
              class="rounded-full px-3 py-1 text-xs font-bold text-white"
              :class="route.aqiClass"
            >
              AQI {{ route.averageAQI }}
            </span>
          </div>

          <div class="mt-5 grid grid-cols-3 gap-3 text-sm font-semibold text-slate-800">
            <span>{{ route.duration }}</span>
            <span>{{ route.distance }}</span>
            <span>Max AQI {{ route.maxAQI }}</span>
          </div>

          <p class="mt-4 text-sm leading-6 text-slate-600">
            {{ capitalizeFirstLetter(route.selectionReason || route.recommendation) }}
          </p>
        </article>
      </div>
    </aside>

    <section class="relative min-h-[500px] bg-[#edf3ee] lg:min-h-full">
      <LMap
        class="absolute inset-0 z-0 h-full w-full"
        :zoom="12"
        :center="mapCenter"
        :bounds="mapBounds"
        :use-global-leaflet="false"
      >
        <LTileLayer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          layer-type="base"
          name="OpenStreetMap"
          attribution="&copy; OpenStreetMap contributors"
        />

        <LPolyline
          v-for="route in routePolylines"
          :key="route.index"
          :lat-lngs="route.latLngs"
          :color="route.color"
          :weight="route.isSelected ? 7 : 5"
          :opacity="route.isSelected ? 0.95 : 0.55"
        >
          <LPopup>{{ route.label }}{{ route.isSelected ? ' (Selected)' : '' }}</LPopup>
        </LPolyline>

        <LCircleMarker
          v-if="source"
          :lat-lng="[source.lat, source.lng]"
          :radius="9"
          color="#047857"
          fill-color="#2f9e52"
          :fill-opacity="1"
        >
          <LPopup>Source: {{ source.label }}</LPopup>
        </LCircleMarker>

        <LCircleMarker
          v-if="destination"
          :lat-lng="[destination.lat, destination.lng]"
          :radius="9"
          color="#b91c1c"
          fill-color="#ef4444"
          :fill-opacity="1"
        >
          <LPopup>Destination: {{ destination.label }}</LPopup>
        </LCircleMarker>

        <LCircleMarker
          v-for="sample in selectedAqiSamples"
          :key="`${sample.lat}-${sample.lng}-${sample.aqi}`"
          :lat-lng="[sample.lat, sample.lng]"
          :radius="5"
          color="#ffffff"
          :weight="1"
          :fill-color="sample.aqi < 80 ? '#2f9e52' : sample.aqi < 120 ? '#fb923c' : '#ef4444'"
          :fill-opacity="0.9"
        >
          <LPopup>AQI {{ formatAQI(sample.aqi) }}</LPopup>
        </LCircleMarker>
      </LMap>

      <div
        class="absolute bottom-6 left-6 right-6 z-[400] rounded-lg border border-slate-200 bg-white/95 px-5 py-4 shadow-[0_10px_30px_rgba(15,23,42,0.12)] backdrop-blur"
      >
        <div class="flex items-center gap-4">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-emerald-50 text-[#2f9e52]">
            <svg viewBox="0 0 24 24" width="22" height="22" class="h-[22px] w-[22px]" fill="none" aria-hidden="true">
              <path d="M5 19 19 5M19 5c-5.9.2-10.5 2.9-12.4 7-1 2.2-.9 4.4.1 5.7 1.3 1 3.5 1.1 5.7.1 4.1-1.9 6.8-6.5 7-12.4Z" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
            </svg>
          </div>
          <div>
            <h2 class="text-sm font-extrabold text-slate-950">You have selected {{ selectedRouteLabel }}</h2>
            <p class="mt-1 text-sm text-slate-600">{{ capitalizeFirstLetter(selectedRouteNote || '') }}</p>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
