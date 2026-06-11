<script setup lang="ts">
defineProps<{
  source: string
  destination: string
}>()

const emit = defineEmits<{
  newSearch: []
}>()

const routes = [
  {
    name: 'Route 1',
    label: 'Recommended',
    color: 'text-[#2f9e52]',
    border: 'border-[#2f9e52]',
    time: '45 min',
    distance: '28.5 km',
    aqi: 62,
    quality: 'Good',
    note: 'Cleaner air along the route with moderate traffic.',
  },
  {
    name: 'Route 2',
    label: '',
    color: 'text-orange-500',
    border: 'border-slate-200',
    time: '48 min',
    distance: '30.2 km',
    aqi: 98,
    quality: 'Moderate',
    note: 'Slightly longer, with acceptable air quality in most areas.',
  },
  {
    name: 'Route 3',
    label: '',
    color: 'text-red-500',
    border: 'border-slate-200',
    time: '40 min',
    distance: '26.1 km',
    aqi: 156,
    quality: 'Unhealthy',
    note: 'Faster, but air quality is poor across several route segments.',
  },
]
</script>

<template>
  <main class="grid flex-1 grid-cols-1 overflow-hidden lg:grid-cols-[320px_1fr]">
    <aside class="border-b border-slate-200 bg-white px-6 py-6 lg:border-b-0 lg:border-r">
      <div class="flex items-start justify-between gap-4">
        <div>
          <h1 class="text-[24px] font-extrabold leading-tight text-black">Recommended Routes</h1>
          <p class="mt-2 text-sm leading-6 text-slate-600">From {{ source }} to {{ destination }}</p>
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
          :key="route.name"
          class="rounded-lg border bg-white p-4 shadow-[0_6px_18px_rgba(15,23,42,0.06)]"
          :class="route.border"
        >
          <div class="flex items-center justify-between gap-3">
            <h2 class="text-base font-extrabold" :class="route.color">
              {{ route.name }} <span v-if="route.label">({{ route.label }})</span>
            </h2>
            <span
              class="rounded-full px-3 py-1 text-xs font-bold text-white"
              :class="
                route.aqi < 80 ? 'bg-[#2f9e52]' : route.aqi < 120 ? 'bg-orange-400' : 'bg-red-500'
              "
            >
              {{ route.aqi }}
            </span>
          </div>

          <div class="mt-5 grid grid-cols-3 gap-3 text-sm font-semibold text-slate-800">
            <span>{{ route.time }}</span>
            <span>{{ route.distance }}</span>
            <span>{{ route.quality }}</span>
          </div>

          <p class="mt-4 text-sm leading-6 text-slate-600">{{ route.note }}</p>
        </article>
      </div>
    </aside>

    <section class="relative min-h-[500px] bg-[#edf3ee]">
      <div
        class="absolute inset-0 bg-[linear-gradient(90deg,rgba(148,163,184,0.18)_1px,transparent_1px),linear-gradient(rgba(148,163,184,0.18)_1px,transparent_1px)] bg-[size:56px_56px]"
      />
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_72%_22%,rgba(186,230,253,0.65),transparent_22%),radial-gradient(circle_at_64%_58%,rgba(187,247,208,0.55),transparent_28%)]" />

      <svg viewBox="0 0 720 460" class="absolute inset-0 h-full w-full" fill="none" preserveAspectRatio="none" aria-hidden="true">
        <path d="M70 145 C160 105 210 70 288 95 S388 155 468 130 610 92 665 62" stroke="#2f9e52" stroke-linecap="round" stroke-width="7" />
        <path d="M70 145 C155 180 230 172 300 215 S440 250 510 292 616 300 665 350" stroke="#fb923c" stroke-linecap="round" stroke-width="7" />
        <path d="M70 145 C125 240 218 225 274 300 S398 330 468 342 585 376 665 350" stroke="#ef4444" stroke-linecap="round" stroke-width="7" />
      </svg>

      <div class="absolute left-[8%] top-[26%]">
        <div class="flex h-12 w-12 items-center justify-center rounded-full bg-[#2f9e52] text-white shadow-lg shadow-emerald-700/20">
          <svg viewBox="0 0 24 24" width="24" height="24" class="h-6 w-6" fill="none" aria-hidden="true">
            <path d="M12 21s7-5.1 7-11a7 7 0 1 0-14 0c0 5.9 7 11 7 11Z" stroke="currentColor" stroke-width="2" />
            <circle cx="12" cy="10" r="2.4" fill="currentColor" />
          </svg>
        </div>
      </div>

      <div class="absolute bottom-[20%] right-[7%]">
        <div class="flex h-12 w-12 items-center justify-center rounded-full bg-red-500 text-white shadow-lg shadow-red-700/20">
          <svg viewBox="0 0 24 24" width="24" height="24" class="h-6 w-6" fill="none" aria-hidden="true">
            <path d="M12 21s7-5.1 7-11a7 7 0 1 0-14 0c0 5.9 7 11 7 11Z" stroke="currentColor" stroke-width="2" />
            <circle cx="12" cy="10" r="2.4" fill="currentColor" />
          </svg>
        </div>
      </div>

      <div class="absolute right-5 top-5 overflow-hidden rounded-md border border-slate-200 bg-white shadow-[0_6px_18px_rgba(15,23,42,0.10)]">
        <button type="button" class="block h-10 w-10 border-b border-slate-200 text-xl font-bold text-slate-700">+</button>
        <button type="button" class="block h-10 w-10 text-xl font-bold text-slate-700">-</button>
      </div>

      <div
        class="absolute bottom-6 left-6 right-6 rounded-lg border border-slate-200 bg-white/95 px-5 py-4 shadow-[0_10px_30px_rgba(15,23,42,0.12)] backdrop-blur"
      >
        <div class="flex items-center gap-4">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-emerald-50 text-[#2f9e52]">
            <svg viewBox="0 0 24 24" width="22" height="22" class="h-[22px] w-[22px]" fill="none" aria-hidden="true">
              <path d="M5 19 19 5M19 5c-5.9.2-10.5 2.9-12.4 7-1 2.2-.9 4.4.1 5.7 1.3 1 3.5 1.1 5.7.1 4.1-1.9 6.8-6.5 7-12.4Z" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
            </svg>
          </div>
          <div>
            <h2 class="text-sm font-extrabold text-slate-950">You have selected Route 1</h2>
            <p class="mt-1 text-sm text-slate-600">This route has the best air quality with moderate travel time.</p>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
