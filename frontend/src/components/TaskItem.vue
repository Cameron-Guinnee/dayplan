<template>
  <div class="flex items-start gap-3 p-3 bg-white rounded-lg border border-gray-200">
    <div class="flex-1 min-w-0">
      <p class="text-sm font-medium text-gray-900 truncate">{{ task.title }}</p>
      <div class="flex items-center gap-2 mt-1">
        <span class="text-xs text-gray-500">{{ task.duration }}m</span>
        <span :class="priorityClass" class="text-xs font-medium px-1.5 py-0.5 rounded">
          {{ priorityLabel }}
        </span>
        <span class="text-xs text-gray-500">due {{ deadlineLabel }}</span>
      </div>
    </div>
    <button
      @click="$emit('complete', task.id)"
      class="shrink-0 w-5 h-5 mt-0.5 rounded border-2 border-gray-300 hover:border-green-500 hover:bg-green-50 transition-colors"
      title="Mark complete"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  task: { type: Object, required: true },
})
defineEmits(['complete'])

const priorityLabel = computed(() => ({ 1: 'Low', 2: 'Med', 3: 'High' }[props.task.priority] ?? 'Med'))

const priorityClass = computed(() => ({
  1: 'bg-gray-100 text-gray-600',
  2: 'bg-yellow-100 text-yellow-700',
  3: 'bg-red-100 text-red-700',
}[props.task.priority] ?? 'bg-gray-100 text-gray-600'))

const deadlineLabel = computed(() => {
  const d = new Date(props.task.deadline)
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
})
</script>
