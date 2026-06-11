<template>
  <main class="flex flex-1 flex-col px-6 pb-8 pt-9">
    <section class="mx-auto w-full max-w-[820px]">
      <p class="text-sm font-bold uppercase text-[#2f9e52]">About</p>
      <h1 class="mt-3 max-w-2xl text-[30px] font-extrabold leading-tight text-black sm:text-[38px]">
        Cleaner route planning powered by live map and air-quality data.
      </h1>
      <p class="mt-4 max-w-3xl text-[16px] leading-7 text-slate-700">
        AQI Route Optimizer helps users compare the air quality along a trip, not just the fastest
        path. It checks AQI across route coordinates, compares that air-quality exposure with time
        and distance, then highlights the route that is healthier for the journey.
      </p>
    </section>

    <section class="mx-auto mt-9 grid w-full max-w-[820px] gap-4 sm:grid-cols-3">
      <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-[0_6px_18px_rgba(15,23,42,0.06)]">
        <div class="mb-4 flex h-10 w-10 items-center justify-center rounded-full bg-emerald-50 text-[#2f9e52]">
          <svg viewBox="0 0 24 24" width="22" height="22" class="h-[22px] w-[22px]" fill="none" aria-hidden="true">
            <path d="M5 19 19 5M19 5c-5.9.2-10.5 2.9-12.4 7-1 2.2-.9 4.4.1 5.7 1.3 1 3.5 1.1 5.7.1 4.1-1.9 6.8-6.5 7-12.4Z" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
          </svg>
        </div>
        <h2 class="text-lg font-bold text-slate-950">AQI-aware routing</h2>
        <p class="mt-2 text-sm leading-6 text-slate-600">
          Samples air quality across the route path instead of relying on a single city-level value.
        </p>
      </article>

      <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-[0_6px_18px_rgba(15,23,42,0.06)]">
        <div class="mb-4 flex h-10 w-10 items-center justify-center rounded-full bg-emerald-50 text-[#2f9e52]">
          <svg viewBox="0 0 24 24" width="22" height="22" class="h-[22px] w-[22px]" fill="none" aria-hidden="true">
            <path d="M4 7h10a4 4 0 0 1 0 8H8m0 0 3-3m-3 3 3 3" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
          </svg>
        </div>
        <h2 class="text-lg font-bold text-slate-950">Time and distance</h2>
        <p class="mt-2 text-sm leading-6 text-slate-600">
          Compares travel duration and distance so a cleaner route still makes practical sense.
        </p>
      </article>

      <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-[0_6px_18px_rgba(15,23,42,0.06)]">
        <div class="mb-4 flex h-10 w-10 items-center justify-center rounded-full bg-emerald-50 text-[#2f9e52]">
          <svg viewBox="0 0 24 24" width="22" height="22" class="h-[22px] w-[22px]" fill="none" aria-hidden="true">
            <path d="M5 13h4l2-7 4 12 2-5h2" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
          </svg>
        </div>
        <h2 class="text-lg font-bold text-slate-950">Clear recommendation</h2>
        <p class="mt-2 text-sm leading-6 text-slate-600">
          Blends AQI, time, and distance into an easy route suggestion you can act on.
        </p>
      </article>
    </section>

    <section class="mx-auto mt-auto w-full max-w-[820px] border-t border-slate-200 pt-6">
      <div class="grid gap-4 text-sm text-slate-700 sm:grid-cols-2">
        <p>
          <span class="font-bold text-slate-950">What we compare:</span>
          AQI readings, route distance, and estimated travel time.
        </p>
        <p>
          <span class="font-bold text-slate-950">What you get:</span>
          A route choice that balances cleaner air with a realistic trip.
        </p>
      </div>
    </section>
  </main>
</template>
