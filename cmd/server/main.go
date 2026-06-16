package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Cameron-Guinnee/dayplan/internal/models"
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
