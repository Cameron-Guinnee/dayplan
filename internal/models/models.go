package models

import "time"

type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Duration    int       `json:"duration"`   // in minutes
    Deadline    time.Time `json:"deadline"`
    Priority    int       `json:"priority"`   // 1 (low) to 3 (high)
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}

type TimeBlock struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
}