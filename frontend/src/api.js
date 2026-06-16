const BASE = ''

export async function getTasks() {
  const res = await fetch(`${BASE}/tasks`)
  if (!res.ok) throw new Error('Failed to fetch tasks')
  return res.json()
}

export async function createTask(task) {
  const res = await fetch(`${BASE}/tasks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(task),
  })
  if (!res.ok) throw new Error('Failed to create task')
  return res.json()
}

export async function completeTask(id) {
  const res = await fetch(`${BASE}/tasks/${id}/complete`, { method: 'PATCH' })
  if (!res.ok) throw new Error('Failed to complete task')
}

export async function getTimeBlocks(date) {
  const dateStr = date.toISOString().slice(0, 10)
  const res = await fetch(`${BASE}/time-blocks?date=${dateStr}`)
  if (!res.ok) throw new Error('Failed to fetch time blocks')
  return res.json()
}

export async function createTimeBlock(block) {
  const res = await fetch(`${BASE}/time-blocks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(block),
  })
  if (!res.ok) throw new Error('Failed to create time block')
  return res.json()
}
