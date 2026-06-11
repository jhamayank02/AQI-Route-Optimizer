<script setup lang="ts">
import { ref } from 'vue'
import AboutPage from './components/AboutPage.vue'
import HowItWorksPage from './components/HowItWorksPage.vue'
import LoadingPage from './components/LoadingPage.vue'
import MapPage from './components/Map.vue'
import SearchBox from './components/SearchBox.vue'

type Page = 'home' | 'about' | 'how' | 'loading' | 'results'
type RouteSearch = {
  source: string
  destination: string
}

const activePage = ref<Page>('home')
const routeSearch = ref<RouteSearch>({
  source: '',
  destination: '',
})

const navItems: { label: string; page: Page }[] = [
  { label: 'Home', page: 'home' },
  { label: 'About', page: 'about' },
  { label: 'How it works', page: 'how' },
]

const startRouteSearch = (search: RouteSearch) => {
  routeSearch.value = search
  activePage.value = 'loading'
}

const showResults = () => {
  activePage.value = 'results'
}

const startNewSearch = () => {
  activePage.value = 'home'
}
</script>

<template>
  <section class="min-h-screen bg-[#f3f6f8] px-5 py-7 text-slate-950">
    <div
      class="mx-auto flex min-h-[640px] max-w-[960px] flex-col overflow-hidden rounded-lg border border-slate-200 bg-white shadow-[0_12px_36px_rgba(15,23,42,0.10)]"
    >
      <header class="flex items-center justify-between border-b border-slate-200 px-8 py-5">
        <button
          type="button"
          class="flex items-center gap-3 text-left text-lg font-bold"
          @click="activePage = 'home'"
        >
          <span class="relative flex h-7 w-7 items-center justify-center">
            <svg
              viewBox="0 0 32 32"
              width="28"
              height="28"
              class="h-7 w-7"
              fill="none"
              aria-hidden="true"
            >
              <path
                d="M25.3 6.7c-7.7.5-13.8 4.4-16.4 10.1-1.4 3.1-1.2 6.3.2 8.1 1.8 1.4 5 1.6 8.1.2 5.7-2.6 9.6-8.7 10.1-16.4.1-1.2-.8-2.1-2-2Z"
                class="fill-emerald-50 stroke-emerald-600"
                stroke-width="2"
              />
              <path
                d="M7 25 24 8M11.5 20.5l-4-4M16 16l-4-4"
                class="stroke-orange-400"
                stroke-linecap="round"
                stroke-width="2"
              />
            </svg>
          </span>
          <span>AQI <span class="text-[#2f9e52]">Route</span> Optimizer</span>
        </button>

        <nav class="hidden items-center gap-9 text-sm font-medium text-slate-800 sm:flex">
          <button
            v-for="item in navItems"
            :key="item.page"
            type="button"
            class="transition hover:text-[#2f9e52]"
            :class="activePage === item.page ? 'text-[#2f9e52]' : 'text-slate-800'"
            @click="activePage = item.page"
          >
            {{ item.label }}
          </button>
        </nav>

        <button
          v-if="activePage === 'results'"
          type="button"
          class="hidden rounded-md border border-[#2f9e52] px-4 py-2 text-sm font-bold text-[#2f9e52] transition hover:bg-emerald-50 sm:block"
          @click="startNewSearch"
        >
          New Search
        </button>
      </header>

      <nav class="flex border-b border-slate-200 px-4 py-3 text-sm font-medium sm:hidden">
        <button
          v-for="item in navItems"
          :key="`mobile-${item.page}`"
          type="button"
          class="flex-1 rounded-md px-3 py-2 transition"
          :class="
            activePage === item.page
              ? 'bg-emerald-50 text-[#2f9e52]'
              : 'text-slate-700 hover:bg-slate-50'
          "
          @click="activePage = item.page"
        >
          {{ item.label }}
        </button>
      </nav>

      <SearchBox v-if="activePage === 'home'" @find-routes="startRouteSearch" />
      <AboutPage v-else-if="activePage === 'about'" />
      <HowItWorksPage v-else-if="activePage === 'how'" />
      <LoadingPage
        v-else-if="activePage === 'loading'"
        :source="routeSearch.source"
        :destination="routeSearch.destination"
        @complete="showResults"
      />
      <MapPage
        v-else
        :source="routeSearch.source"
        :destination="routeSearch.destination"
        @new-search="startNewSearch"
      />
    </div>
  </section>
</template>
