// Package scheduler implements a greedy Earliest-Deadline-First algorithm
// that fits tasks into the free gaps of a day around fixed time blocks.
//
// It is intentionally free of DB and HTTP imports so it can be unit-tested
// in isolation — call Schedule with plain data, get a plain slice back.
package scheduler

import (
	"sort"
	"time"

	"github.com/Cameron-Guinnee/dayplan/internal/models"
)

// ScheduledTask is a task that has been assigned a concrete start and end time.
type ScheduledTask struct {
	Task      models.Task
	StartTime time.Time
	EndTime   time.Time
}

// freeWindow is an open time interval within the day that tasks can fill.
type freeWindow struct {
	start time.Time
	end   time.Time
}

// Schedule fits tasks into the free gaps of dayStart–dayEnd around the given
// fixed time blocks.
//
// Tasks are ordered by earliest deadline first; equal deadlines break on
// descending priority (3 = high scheduled before 1 = low). Tasks that do not
// fit in any remaining window before their deadline are skipped and returned
// in unscheduled.
func Schedule(
	tasks []models.Task,
	blocks []models.TimeBlock,
	dayStart, dayEnd time.Time,
) (scheduled []ScheduledTask, unscheduled []models.Task) {
	free := freeWindows(blocks, dayStart, dayEnd)
	candidates := pendingTasks(tasks)

	// EDF ordering; ties broken by descending priority.
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Deadline.Equal(candidates[j].Deadline) {
			return candidates[i].Priority > candidates[j].Priority
		}
		return candidates[i].Deadline.Before(candidates[j].Deadline)
	})

	// wi tracks how far into the current free window we have consumed.
	wi := 0
	windowCursor := make([]time.Time, len(free))
	for i, w := range free {
		windowCursor[i] = w.start
	}

	for _, task := range candidates {
		dur := time.Duration(task.Duration) * time.Minute
		placed := false

		for wi < len(free) {
			windowEnd := free[wi].end
			available := windowEnd.Sub(windowCursor[wi])

			if available <= 0 {
				wi++
				continue
			}

			start := windowCursor[wi]
			end := start.Add(dur)

			if end.After(windowEnd) {
				// Task does not fit in the remainder of this window; advance.
				wi++
				continue
			}

			// Respect the deadline: the task must finish by its deadline.
			if end.After(task.Deadline) {
				break // later windows are even further out; skip task
			}

			scheduled = append(scheduled, ScheduledTask{
				Task:      task,
				StartTime: start,
				EndTime:   end,
			})
			windowCursor[wi] = end
			placed = true
			break
		}

		if !placed {
			unscheduled = append(unscheduled, task)
		}
	}

	return scheduled, unscheduled
}

// freeWindows returns the gaps within [dayStart, dayEnd] not covered by any
// time block. Overlapping or adjacent blocks are merged before gap extraction.
func freeWindows(blocks []models.TimeBlock, dayStart, dayEnd time.Time) []freeWindow {
	// Clamp blocks to the day boundary and sort by start.
	type interval struct{ start, end time.Time }
	var intervals []interval
	for _, b := range blocks {
		s := maxTime(b.StartTime, dayStart)
		e := minTime(b.EndTime, dayEnd)
		if s.Before(e) {
			intervals = append(intervals, interval{s, e})
		}
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start.Before(intervals[j].start)
	})

	// Merge overlapping intervals.
	merged := make([]interval, 0, len(intervals))
	for _, iv := range intervals {
		if len(merged) == 0 || iv.start.After(merged[len(merged)-1].end) {
			merged = append(merged, iv)
		} else if iv.end.After(merged[len(merged)-1].end) {
			merged[len(merged)-1].end = iv.end
		}
	}

	// Collect gaps.
	var windows []freeWindow
	cursor := dayStart
	for _, iv := range merged {
		if cursor.Before(iv.start) {
			windows = append(windows, freeWindow{cursor, iv.start})
		}
		cursor = iv.end
	}
	if cursor.Before(dayEnd) {
		windows = append(windows, freeWindow{cursor, dayEnd})
	}
	return windows
}

// pendingTasks filters out completed tasks.
func pendingTasks(tasks []models.Task) []models.Task {
	out := make([]models.Task, 0, len(tasks))
	for _, t := range tasks {
		if !t.Completed {
			out = append(out, t)
		}
	}
	return out
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
