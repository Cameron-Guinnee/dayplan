<template>
  <div class="flex flex-col gap-2">
    <TaskItem
      v-for="task in tasks"
      :key="task.id"
      :task="task"
      @complete="$emit('complete', $event)"
    />
    <p v-if="tasks.length === 0" class="text-sm text-gray-400 text-center py-4">
      No pending tasks
    </p>
  </div>
</template>

<script setup>
import TaskItem from './TaskItem.vue'

defineProps({
  tasks: { type: Array, required: true },
})
defineEmits(['complete'])
</script>
