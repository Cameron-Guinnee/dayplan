<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- Top bar -->
    <header class="bg-white border-b border-gray-200 px-6 py-3 flex items-center gap-3">
      <h1 class="text-lg font-semibold text-gray-900 tracking-tight">dayplan</h1>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar: tasks -->
      <aside class="w-80 shrink-0 bg-white border-r border-gray-200 flex flex-col overflow-y-auto">
        <div class="p-4 border-b border-gray-100">
          <h2 class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">Tasks</h2>
          <TaskList :tasks="tasks" @complete="handleComplete" />
        </div>
        <div class="p-4">
          <h2 class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">Add Task</h2>
          <TaskForm @submit="handleCreateTask" />
        </div>
      </aside>

      <!-- Main: date nav + timeline + time block form -->
      <main class="flex-1 flex flex-col overflow-y-auto">
        <!-- Date navigation -->
        <div class="sticky top-0 z-10 bg-white border-b border-gray-200 px-6 py-3 flex items-center gap-4">
          <button @click="shiftDate(-1)" class="btn-ghost">&#8592;</button>
          <span class="text-sm font-medium text-gray-700 w-36 text-center">{{ dateLabel }}</span>
          <button @click="shiftDate(1)" class="btn-ghost">&#8594;</button>
        </div>

        <!-- Timeline -->
        <div class="flex-1 px-6 py-4 overflow-x-hidden">
          <TimelineView :timeBlocks="timeBlocks" />
        </div>

        <!-- Add time block -->
        <div class="border-t border-gray-200 bg-white px-6 py-4">
          <h2 class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">Add Time Block</h2>
          <TimeBlockForm :date="selectedDate" @submit="handleCreateTimeBlock" />
        </div>
      </main>
    </div>

    <!-- Error toast -->
    <div
      v-if="error"
      class="fixed bottom-4 right-4 bg-red-600 text-white text-sm px-4 py-2 rounded shadow-lg"
    >
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import TaskList from './components/TaskList.vue'
import TaskForm from './components/TaskForm.vue'
import TimelineView from './components/TimelineView.vue'
import TimeBlockForm from './components/TimeBlockForm.vue'
import { getTasks, createTask, completeTask, getTimeBlocks, createTimeBlock } from './api.js'

const tasks = ref([])
const timeBlocks = ref([])
const selectedDate = ref(new Date())
const error = ref(null)

const dateLabel = computed(() =>
  selectedDate.value.toLocaleDateString(undefined, { weekday: 'long', month: 'long', day: 'numeric' })
)

function shiftDate(delta) {
  const d = new Date(selectedDate.value)
  d.setDate(d.getDate() + delta)
  selectedDate.value = d
}

function showError(msg) {
  error.value = msg
  setTimeout(() => (error.value = null), 4000)
}

async function loadTasks() {
  try { tasks.value = await getTasks() }
  catch (e) { showError(e.message) }
}

async function loadTimeBlocks() {
  try { timeBlocks.value = await getTimeBlocks(selectedDate.value) }
  catch (e) { showError(e.message) }
}

async function handleCreateTask(task) {
  try {
    const created = await createTask(task)
    tasks.value.push(created)
  } catch (e) { showError(e.message) }
}

async function handleComplete(id) {
  try {
    await completeTask(id)
    tasks.value = tasks.value.filter(t => t.id !== id)
  } catch (e) { showError(e.message) }
}

async function handleCreateTimeBlock(block) {
  try {
    const created = await createTimeBlock(block)
    timeBlocks.value.push(created)
  } catch (e) { showError(e.message) }
}

watch(selectedDate, loadTimeBlocks)
onMounted(() => {
  loadTasks()
  loadTimeBlocks()
})
</script>
