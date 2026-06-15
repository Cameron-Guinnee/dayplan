package store

import (
    "database/sql"
    "time"

    "github.com/Cameron-Guinnee/dayplan/internal/models"
    _ "modernc.org/sqlite"
)

type Store struct {
    db *sql.DB
}

func New(dbPath string) (*Store, error) {
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    s := &Store{db: db}
    if err := s.migrate(); err != nil {
        return nil, err
    }

    return s, nil
}

func (s *Store) migrate() error {
    _, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id          INTEGER PRIMARY KEY AUTOINCREMENT,
            title       TEXT NOT NULL,
            duration    INTEGER NOT NULL,
            deadline    DATETIME NOT NULL,
            priority    INTEGER NOT NULL DEFAULT 2,
            completed   BOOLEAN NOT NULL DEFAULT 0,
            created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS time_blocks (
            id          INTEGER PRIMARY KEY AUTOINCREMENT,
            title       TEXT NOT NULL,
            start_time  DATETIME NOT NULL,
            end_time    DATETIME NOT NULL
        );
    `)
    return err
}

func (s *Store) CreateTask(t *models.Task) error {
    res, err := s.db.Exec(
        `INSERT INTO tasks (title, duration, deadline, priority) VALUES (?, ?, ?, ?)`,
        t.Title, t.Duration, t.Deadline, t.Priority,
    )
    if err != nil {
        return err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return err
    }
    t.ID = int(id)
    return nil
}

func (s *Store) GetTasks() ([]models.Task, error) {
    rows, err := s.db.Query(
        `SELECT id, title, duration, deadline, priority, completed, created_at FROM tasks WHERE completed = 0`,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []models.Task
    for rows.Next() {
        var t models.Task
        err := rows.Scan(&t.ID, &t.Title, &t.Duration, &t.Deadline, &t.Priority, &t.Completed, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }
    return tasks, nil
}

func (s *Store) CompleteTask(id int) error {
    _, err := s.db.Exec(`UPDATE tasks SET completed = 1 WHERE id = ?`, id)
    return err
}

func (s *Store) CreateTimeBlock(tb *models.TimeBlock) error {
    res, err := s.db.Exec(
        `INSERT INTO time_blocks (title, start_time, end_time) VALUES (?, ?, ?)`,
        tb.Title, tb.StartTime, tb.EndTime,
    )
    if err != nil {
        return err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return err
    }
    tb.ID = int(id)
    return nil
}

func (s *Store) GetTimeBlocks(date time.Time) ([]models.TimeBlock, error) {
    start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    end := start.Add(24 * time.Hour)

    rows, err := s.db.Query(
        `SELECT id, title, start_time, end_time FROM time_blocks WHERE start_time >= ? AND start_time < ?`,
        start, end,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var blocks []models.TimeBlock
    for rows.Next() {
        var tb models.TimeBlock
        err := rows.Scan(&tb.ID, &tb.Title, &tb.StartTime, &tb.EndTime)
        if err != nil {
            return nil, err
        }
        blocks = append(blocks, tb)
    }
    return blocks, nil
}