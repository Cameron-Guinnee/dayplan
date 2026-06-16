<template>
  <!-- 06:00–22:00, each hour = 64px -->
  <div class="relative" :style="{ height: `${totalHeight}px` }">
    <!-- Hour grid lines and labels -->
    <div
      v-for="hour in hours"
      :key="hour"
      class="absolute w-full flex items-start"
      :style="{ top: `${(hour - START_HOUR) * HOUR_PX}px` }"
    >
      <span class="w-12 shrink-0 text-xs text-gray-400 -mt-2 select-none">
        {{ String(hour).padStart(2, '0') }}:00
      </span>
      <div class="flex-1 border-t border-gray-100 mt-0" />
    </div>

    <!-- Time blocks -->
    <div
      v-for="block in timeBlocks"
      :key="block.id"
      :style="blockStyle(block)"
      class="absolute left-14 right-0 rounded bg-blue-200 border border-blue-400 px-2 py-1 overflow-hidden"
      :title="`${block.title}: ${fmt(block.start_time)} – ${fmt(block.end_time)}`"
    >
      <p class="text-xs font-medium text-blue-900 truncate">{{ block.title }}</p>
      <p class="text-xs text-blue-700">{{ fmt(block.start_time) }} – {{ fmt(block.end_time) }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  timeBlocks: { type: Array, required: true },
})

const START_HOUR = 6
const END_HOUR = 22
const HOUR_PX = 64

const totalHeight = (END_HOUR - START_HOUR) * HOUR_PX
const hours = Array.from({ length: END_HOUR - START_HOUR + 1 }, (_, i) => START_HOUR + i)

function minutesFromMidnight(isoString) {
  const d = new Date(isoString)
  return d.getHours() * 60 + d.getMinutes()
}

function blockStyle(block) {
  const startMin = minutesFromMidnight(block.start_time)
  const endMin = minutesFromMidnight(block.end_time)
  const top = (startMin - START_HOUR * 60) * (HOUR_PX / 60)
  const height = Math.max((endMin - startMin) * (HOUR_PX / 60), 20)
  return { top: `${top}px`, height: `${height}px` }
}

function fmt(isoString) {
  return new Date(isoString).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })
}
</script>
