package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/Cameron-Guinnee/dayplan/internal/store"
)

func main() {
    db, err := store.New("dayplan.db")
    if err != nil {
        log.Fatalf("failed to open database: %v", err)
    }

    r := chi.NewRouter()

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "dayplan is running")
    })

    _ = db // we'll use this shortly

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", r)
}