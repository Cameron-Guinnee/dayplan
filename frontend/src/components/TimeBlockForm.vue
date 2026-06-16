<template>
  <form @submit.prevent="submit" class="flex flex-col gap-3">
    <input
      v-model="title"
      required
      placeholder="Time block title"
      class="input"
    />
    <div class="flex gap-2">
      <input v-model="startTime" required type="time" class="input w-1/2" />
      <input v-model="endTime" required type="time" class="input w-1/2" />
    </div>
    <button type="submit" class="btn-primary">Add Time Block</button>
  </form>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  date: { type: Date, required: true },
})
const emit = defineEmits(['submit'])

const title = ref('')
const startTime = ref('')
const endTime = ref('')

function toISO(timeStr) {
  const [h, m] = timeStr.split(':').map(Number)
  const d = new Date(props.date)
  d.setHours(h, m, 0, 0)
  return d.toISOString()
}

function submit() {
  emit('submit', {
    title: title.value,
    start_time: toISO(startTime.value),
    end_time: toISO(endTime.value),
  })
  title.value = ''
  startTime.value = ''
  endTime.value = ''
}
</script>
