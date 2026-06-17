package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Cameron-Guinnee/dayplan/internal/models"
	"github.com/Cameron-Guinnee/dayplan/internal/scheduler"
	"github.com/Cameron-Guinnee/dayplan/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	db, err := store.New("dayplan.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "dayplan is running")
	})

	r.Post("/tasks", createTask(db))
	r.Get("/tasks", listTasks(db))
	r.Patch("/tasks/{id}/complete", completeTask(db))

	r.Post("/time-blocks", createTimeBlock(db))
	r.Get("/time-blocks", listTimeBlocks(db))

	r.Get("/schedule", getSchedule(db))

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// createTask decodes a Task from the request body and persists it.
// The store fills in the generated ID before we encode the response.
func createTask(db *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t models.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := db.CreateTask(&t); err != nil {
			http.Error(w, "failed to create task", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	}
}

func listTasks(db *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := db.GetTasks()
		if err != nil {
			http.Error(w, "failed to list tasks", http.StatusInternalServerError)
			return
		}
		// Return an empty array rather than null when there are no tasks.
		if tasks == nil {
			tasks = []models.Task{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func completeTask(db *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "invalid task id", http.StatusBadRequest)
			return
		}
		if err := db.CompleteTask(id); err != nil {
			http.Error(w, "failed to complete task", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func createTimeBlock(db *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tb models.TimeBlock
		if err := json.NewDecoder(r.Body).Decode(&tb); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := db.CreateTimeBlock(&tb); err != nil {
			http.Error(w, "failed to create time block", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tb)
	}
}

// getSchedule runs the EDF scheduler for a given day and returns the result.
//
// Query params:
//
//	date      – YYYY-MM-DD, defaults to today
//	day_start – HH:MM, start of the schedulable window, defaults to 09:00
//	day_end   – HH:MM, end of the schedulable window, defaults to 17:00
//
// Response shape:
//
//	{ "scheduled": [...], "unscheduled": [...] }
func getSchedule(db *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := time.Now()
		if dateStr := r.URL.Query().Get("date"); dateStr != "" {
			var err error
			date, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		}

		dayStart, err := parseDayTime(date, r.URL.Query().Get("day_start"), 9, 0)
		if err != nil {
			http.Error(w, "invalid day_start, use HH:MM", http.StatusBadRequest)
			return
		}
		dayEnd, err := parseDayTime(date, r.URL.Query().Get("day_end"), 17, 0)
		if err != nil {
			http.Error(w, "invalid day_end, use HH:MM", http.StatusBadRequest)
			return
		}
		if !dayEnd.After(dayStart) {
			http.Error(w, "day_end must be after day_start", http.StatusBadRequest)
			return
		}

		tasks, err := db.GetTasks()
		if err != nil {
			http.Error(w, "failed to load tasks", http.StatusInternalServerError)
			return
		}
		blocks, err := db.GetTimeBlocks(date)
		if err != nil {
			http.Error(w, "failed to load time blocks", http.StatusInternalServerError)
			return
		}

		scheduled, unscheduled := scheduler.Schedule(tasks, blocks, dayStart, dayEnd)

		// Prefer empty arrays over null in the JSON output.
		if scheduled == nil {
			scheduled = []scheduler.ScheduledTask{}
		}
		if unscheduled == nil {
			unscheduled = []models.Task{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Scheduled   []scheduler.ScheduledTask `json:"scheduled"`
			Unscheduled []models.Task             `json:"unscheduled"`
		}{scheduled, unscheduled})
	}
}

// parseDayTime builds a time.Time from date + an optional "HH:MM" string.
// If the string is empty, defaultH and defaultM are used instead.
func parseDayTime(date time.Time, s string, defaultH, defaultM int) (time.Time, error) {
	h, m := defaultH, defaultM
	if s != "" {
		t, err := time.Parse("15:04", s)
		if err != nil {
			return time.Time{}, err
		}
		h, m = t.Hour(), t.Minute()
	}
	return time.Date(date.Year(), date.Month(), date.Day(), h, m, 0, 0, date.Location()), nil
}

// listTimeBlocks returns time blocks for a given day.
// The date query param uses YYYY-MM-DD format; defaults to today when omitted.
func listTimeBlocks(db *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := time.Now()
		if dateStr := r.URL.Query().Get("date"); dateStr != "" {
			var err error
			date, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		}
		blocks, err := db.GetTimeBlocks(date)
		if err != nil {
			http.Error(w, "failed to list time blocks", http.StatusInternalServerError)
			return
		}
		if blocks == nil {
			blocks = []models.TimeBlock{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blocks)
	}
}
