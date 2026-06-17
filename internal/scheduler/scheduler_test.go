package scheduler

import (
	"testing"
	"time"

	"github.com/Cameron-Guinnee/dayplan/internal/models"
)

// day returns a time on an arbitrary reference date at h:m.
func day(h, m int) time.Time {
	return time.Date(2024, 1, 1, h, m, 0, 0, time.UTC)
}

func TestSchedule_BasicEDF(t *testing.T) {
	// Two tasks that fit comfortably in a free day; earlier deadline goes first.
	tasks := []models.Task{
		{ID: 1, Title: "Later", Duration: 60, Deadline: day(18, 0), Priority: 2},
		{ID: 2, Title: "Sooner", Duration: 30, Deadline: day(12, 0), Priority: 2},
	}

	sched, unsched := Schedule(tasks, nil, day(9, 0), day(17, 0))

	if len(unsched) != 0 {
		t.Fatalf("expected no unscheduled tasks, got %v", unsched)
	}
	if len(sched) != 2 {
		t.Fatalf("expected 2 scheduled tasks, got %d", len(sched))
	}
	if sched[0].Task.ID != 2 {
		t.Errorf("expected task 2 (Sooner) first, got task %d", sched[0].Task.ID)
	}
}

func TestSchedule_PriorityTiebreak(t *testing.T) {
	// Same deadline — higher priority should be scheduled first.
	deadline := day(17, 0)
	tasks := []models.Task{
		{ID: 1, Title: "Low", Duration: 30, Deadline: deadline, Priority: 1},
		{ID: 2, Title: "High", Duration: 30, Deadline: deadline, Priority: 3},
	}

	sched, _ := Schedule(tasks, nil, day(9, 0), day(17, 0))

	if len(sched) < 2 {
		t.Fatalf("expected both tasks scheduled, got %d", len(sched))
	}
	if sched[0].Task.ID != 2 {
		t.Errorf("expected high-priority task first, got task %d", sched[0].Task.ID)
	}
}

func TestSchedule_SkipsFixedBlocks(t *testing.T) {
	// A lunch block from 12:00–13:00; task must not overlap it.
	blocks := []models.TimeBlock{
		{ID: 1, Title: "Lunch", StartTime: day(12, 0), EndTime: day(13, 0)},
	}
	tasks := []models.Task{
		{ID: 1, Title: "Work", Duration: 90, Deadline: day(18, 0), Priority: 2},
	}

	sched, unsched := Schedule(tasks, blocks, day(9, 0), day(17, 0))

	if len(unsched) != 0 {
		t.Fatalf("unexpected unscheduled: %v", unsched)
	}
	if len(sched) != 1 {
		t.Fatalf("expected 1 scheduled task, got %d", len(sched))
	}
	st := sched[0]
	// Start must be before lunch or at/after lunch end.
	overlaps := st.StartTime.Before(day(13, 0)) && st.EndTime.After(day(12, 0))
	if overlaps {
		t.Errorf("scheduled task overlaps lunch block: %v–%v", st.StartTime, st.EndTime)
	}
}

func TestSchedule_UnscheduledWhenNoFit(t *testing.T) {
	// A 4-hour task with a deadline before there is any room for it.
	tasks := []models.Task{
		{ID: 1, Title: "Impossible", Duration: 240, Deadline: day(10, 0), Priority: 2},
	}

	_, unsched := Schedule(tasks, nil, day(9, 0), day(17, 0))

	// The task needs 4 h but its deadline is 10:00 — only 1 h of window before deadline.
	if len(unsched) != 1 {
		t.Errorf("expected task to be unscheduled, got scheduled count != 0")
	}
}

func TestSchedule_ConsecutiveTasks(t *testing.T) {
	// Two tasks should be placed back-to-back, not overlap.
	tasks := []models.Task{
		{ID: 1, Title: "A", Duration: 60, Deadline: day(17, 0), Priority: 2},
		{ID: 2, Title: "B", Duration: 60, Deadline: day(17, 0), Priority: 1},
	}

	sched, unsched := Schedule(tasks, nil, day(9, 0), day(17, 0))

	if len(unsched) != 0 {
		t.Fatalf("unexpected unscheduled: %v", unsched)
	}
	if len(sched) != 2 {
		t.Fatalf("expected 2 scheduled, got %d", len(sched))
	}
	// The second task must start no earlier than the first ends.
	if sched[1].StartTime.Before(sched[0].EndTime) {
		t.Errorf("tasks overlap: first ends %v, second starts %v",
			sched[0].EndTime, sched[1].StartTime)
	}
}

func TestFreeWindows_MergesOverlappingBlocks(t *testing.T) {
	blocks := []models.TimeBlock{
		{StartTime: day(10, 0), EndTime: day(12, 0)},
		{StartTime: day(11, 0), EndTime: day(13, 0)}, // overlaps previous
	}

	windows := freeWindows(blocks, day(9, 0), day(17, 0))

	// Expect: 09:00–10:00 and 13:00–17:00
	if len(windows) != 2 {
		t.Fatalf("expected 2 free windows, got %d", len(windows))
	}
	if !windows[0].start.Equal(day(9, 0)) || !windows[0].end.Equal(day(10, 0)) {
		t.Errorf("unexpected first window: %v–%v", windows[0].start, windows[0].end)
	}
	if !windows[1].start.Equal(day(13, 0)) || !windows[1].end.Equal(day(17, 0)) {
		t.Errorf("unexpected second window: %v–%v", windows[1].start, windows[1].end)
	}
}
