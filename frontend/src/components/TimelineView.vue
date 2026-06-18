<template>
  <div class="flex flex-col gap-4">
    <!-- Grid: 06:00–22:00, each hour = 64px -->
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

      <!-- Fixed time blocks (blue) -->
      <div
        v-for="block in timeBlocks"
        :key="`tb-${block.id}`"
        :style="rangeStyle(block.start_time, block.end_time)"
        class="absolute left-14 right-0 rounded bg-blue-100 border border-blue-300 px-2 py-1 overflow-hidden"
        :title="`${block.title}: ${fmt(block.start_time)} – ${fmt(block.end_time)}`"
      >
        <p class="text-xs font-semibold text-blue-900 truncate">{{ block.title }}</p>
        <p class="text-xs text-blue-600">{{ fmt(block.start_time) }} – {{ fmt(block.end_time) }}</p>
      </div>

      <!-- Scheduled tasks (green) -->
      <div
        v-for="item in scheduledTasks"
        :key="`st-${item.task.id}`"
        :style="rangeStyle(item.start_time, item.end_time)"
        class="absolute left-14 right-0 rounded bg-emerald-100 border border-emerald-300 px-2 py-1 overflow-hidden"
        :title="`${item.task.title}: ${fmt(item.start_time)} – ${fmt(item.end_time)}`"
      >
        <p class="text-xs font-semibold text-emerald-900 truncate">{{ item.task.title }}</p>
        <p class="text-xs text-emerald-600">{{ fmt(item.start_time) }} – {{ fmt(item.end_time) }}</p>
      </div>
    </div>

    <!-- Unscheduled tasks -->
    <div v-if="unscheduled.length > 0" class="rounded border border-amber-200 bg-amber-50 px-4 py-3">
      <p class="text-xs font-semibold text-amber-700 uppercase tracking-wide mb-2">
        Could not be scheduled ({{ unscheduled.length }})
      </p>
      <ul class="flex flex-col gap-1">
        <li
          v-for="task in unscheduled"
          :key="task.id"
          class="text-xs text-amber-800"
        >
          {{ task.title }} — {{ task.duration }} min, due {{ fmtDeadline(task.deadline) }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  timeBlocks:     { type: Array, required: true },
  scheduledTasks: { type: Array, default: () => [] },
  unscheduled:    { type: Array, default: () => [] },
})

const START_HOUR = 6
const END_HOUR   = 22
const HOUR_PX    = 64

const totalHeight = (END_HOUR - START_HOUR) * HOUR_PX
const hours = Array.from({ length: END_HOUR - START_HOUR + 1 }, (_, i) => START_HOUR + i)

function minutesFromMidnight(isoString) {
  const d = new Date(isoString)
  return d.getHours() * 60 + d.getMinutes()
}

function rangeStyle(startIso, endIso) {
  const startMin = minutesFromMidnight(startIso)
  const endMin   = minutesFromMidnight(endIso)
  const top    = (startMin - START_HOUR * 60) * (HOUR_PX / 60)
  const height = Math.max((endMin - startMin) * (HOUR_PX / 60), 20)
  return { top: `${top}px`, height: `${height}px` }
}

function fmt(isoString) {
  return new Date(isoString).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })
}

function fmtDeadline(isoString) {
  return new Date(isoString).toLocaleString(undefined, {
    month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}
</script>
