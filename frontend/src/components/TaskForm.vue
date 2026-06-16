<template>
  <form @submit.prevent="submit" class="flex flex-col gap-3">
    <input
      v-model="title"
      required
      placeholder="Task title"
      class="input"
    />
    <div class="flex gap-2">
      <input
        v-model.number="duration"
        required
        type="number"
        min="1"
        placeholder="Duration (min)"
        class="input w-1/2"
      />
      <select v-model.number="priority" class="input w-1/2">
        <option value="1">Low</option>
        <option value="2">Medium</option>
        <option value="3">High</option>
      </select>
    </div>
    <input
      v-model="deadline"
      required
      type="datetime-local"
      class="input"
    />
    <button type="submit" class="btn-primary">Add Task</button>
  </form>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['submit'])

const title = ref('')
const duration = ref(30)
const priority = ref(2)
const deadline = ref('')

function submit() {
  emit('submit', {
    title: title.value,
    duration: duration.value,
    priority: priority.value,
    // Convert local datetime string to an ISO 8601 string the Go backend can parse.
    deadline: new Date(deadline.value).toISOString(),
  })
  title.value = ''
  duration.value = 30
  priority.value = 2
  deadline.value = ''
}
</script>
