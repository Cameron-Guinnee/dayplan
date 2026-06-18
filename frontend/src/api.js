const BASE = ''

// toDateStr formats a Date as YYYY-MM-DD using local time, not UTC.
// Using toISOString() would shift to UTC first, sending the wrong date
// for users in negative-offset timezones late in the evening.
function toDateStr(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

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
  const dateStr = toDateStr(date)
  const res = await fetch(`${BASE}/time-blocks?date=${dateStr}`)
  if (!res.ok) throw new Error('Failed to fetch time blocks')
  return res.json()
}

export async function getSchedule(date) {
  const dateStr = toDateStr(date)
  const res = await fetch(`${BASE}/schedule?date=${dateStr}`)
  if (!res.ok) throw new Error('Failed to fetch schedule')
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
