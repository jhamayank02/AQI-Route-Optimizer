<script setup lang="ts">
import { useStore } from './store/store'

const store = useStore()

const navItems = [
  { link: '/', label: 'Home' },
  { link: '/about', label: 'About' },
  { link: '/how-it-works', label: 'How it works' }
]

const handleNavClick = (link: string) => {
  if (link !== '/') {
    return
  }

  store.clearLocation()
  store.clearRecommendedRoute()
}

</script>

<template>
  <section class="h-screen overflow-hidden bg-white text-slate-950 sm:bg-[#f3f6f8] sm:px-5 sm:py-7">
    <div
      class="mx-auto flex h-full min-h-0 w-full flex-col overflow-hidden bg-white sm:max-w-[960px] sm:rounded-lg sm:border sm:border-slate-200 sm:shadow-[0_12px_36px_rgba(15,23,42,0.10)]">
      <header class="flex shrink-0 items-center justify-between border-b border-slate-200 px-4 py-4 sm:px-8 sm:py-5">
        <RouterLink to="/" class="flex min-w-0 items-center gap-3 text-left text-base font-bold sm:text-lg" @click="handleNavClick('/')">
          <span class="relative flex h-7 w-7 items-center justify-center">
            <svg viewBox="0 0 32 32" width="28" height="28" class="h-7 w-7" fill="none" aria-hidden="true">
              <path
                d="M25.3 6.7c-7.7.5-13.8 4.4-16.4 10.1-1.4 3.1-1.2 6.3.2 8.1 1.8 1.4 5 1.6 8.1.2 5.7-2.6 9.6-8.7 10.1-16.4.1-1.2-.8-2.1-2-2Z"
                class="fill-emerald-50 stroke-emerald-600" stroke-width="2" />
              <path d="M7 25 24 8M11.5 20.5l-4-4M16 16l-4-4" class="stroke-orange-400" stroke-linecap="round"
                stroke-width="2" />
            </svg>
          </span>
          <span class="min-w-0 truncate">AQI <span class="text-[#2f9e52]">Route</span> Optimizer</span>
        </RouterLink>

        <nav class="hidden items-center gap-9 text-sm font-medium text-slate-800 sm:flex">
          <RouterLink v-for="item in navItems" :key="item.link" :to="item.link" @click="handleNavClick(item.link)">
            {{ item.label }}
          </RouterLink>
        </nav>
      </header>
      <div class="min-h-0 flex-1 overflow-y-auto overflow-x-hidden">
        <RouterView />
      </div>
    </div>
  </section>
</template>
