<template>
  <main class="flex flex-1 flex-col px-6 pb-8 pt-9">
    <section class="mx-auto w-full max-w-[820px] text-center">
      <p class="text-sm font-bold uppercase text-[#2f9e52]">How it works</p>
      <h1 class="mt-3 text-[30px] font-extrabold leading-tight text-black sm:text-[38px]">
        From two points to one cleaner route recommendation.
      </h1>
      <p class="mx-auto mt-4 max-w-2xl text-[16px] leading-7 text-slate-700">
        We look at possible route coordinates, gather AQI values along the way, then compare air
        quality with travel time and distance to suggest a healthier choice.
      </p>
    </section>

    <section class="mx-auto mt-9 w-full max-w-[760px]">
      <div class="space-y-4">
        <article class="flex gap-4 rounded-lg border border-slate-200 bg-white p-5 shadow-[0_6px_18px_rgba(15,23,42,0.06)]">
          <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-[#2f9e52] text-sm font-bold text-white">
            1
          </div>
          <div>
            <h2 class="text-lg font-bold text-slate-950">Search source and destination</h2>
            <p class="mt-1 text-sm leading-6 text-slate-600">
              Start by entering where the trip begins and where you want to go.
            </p>
          </div>
        </article>

        <article class="flex gap-4 rounded-lg border border-slate-200 bg-white p-5 shadow-[0_6px_18px_rgba(15,23,42,0.06)]">
          <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-[#2f9e52] text-sm font-bold text-white">
            2
          </div>
          <div>
            <h2 class="text-lg font-bold text-slate-950">Compare possible routes</h2>
            <p class="mt-1 text-sm leading-6 text-slate-600">
              Each route is broken into coordinates so we can understand the path, distance, and time.
            </p>
          </div>
        </article>

        <article class="flex gap-4 rounded-lg border border-slate-200 bg-white p-5 shadow-[0_6px_18px_rgba(15,23,42,0.06)]">
          <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-[#2f9e52] text-sm font-bold text-white">
            3
          </div>
          <div>
            <h2 class="text-lg font-bold text-slate-950">Sample AQI along the path</h2>
            <p class="mt-1 text-sm leading-6 text-slate-600">
              AQI values are gathered across those coordinates instead of using only one city-level reading.
            </p>
          </div>
        </article>

        <article class="flex gap-4 rounded-lg border border-slate-200 bg-white p-5 shadow-[0_6px_18px_rgba(15,23,42,0.06)]">
          <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-[#2f9e52] text-sm font-bold text-white">
            4
          </div>
          <div>
            <h2 class="text-lg font-bold text-slate-950">Return the recommendation</h2>
            <p class="mt-1 text-sm leading-6 text-slate-600">
              The final choice balances cleaner air, shorter distance, and reasonable travel time.
            </p>
          </div>
        </article>
      </div>
    </section>

    <section class="mx-auto mt-auto w-full max-w-[820px] border-t border-slate-200 pt-6">
      <div class="grid gap-3 text-sm font-semibold text-slate-700 sm:grid-cols-3">
        <div class="rounded-md bg-emerald-50 px-4 py-3 text-center text-[#2f9e52]">
          AQI along route
        </div>
        <div class="rounded-md bg-emerald-50 px-4 py-3 text-center text-[#2f9e52]">
          Travel time
        </div>
        <div class="rounded-md bg-emerald-50 px-4 py-3 text-center text-[#2f9e52]">
          Distance
        </div>
      </div>
    </section>
  </main>
</template>
