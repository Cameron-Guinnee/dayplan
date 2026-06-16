# dayplan

A day-planning web app for tracking tasks (with duration, deadline, and
priority) alongside fixed time blocks (meetings, appointments, etc.) on a
visual daily timeline.

This started as a personal tool to solve my own daily planning problem, and
doubles as a way to get hands-on with Go.

## Why this exists

Most to-do apps treat every task the same and leave you to work out when
things actually fit into your day. dayplan captures the information that
matters for that — duration, deadline, priority — and lays your fixed
commitments out on a visual timeline so you can see what your day actually
looks like at a glance.

**Where this is headed:** the next milestone is an automatic scheduler that
takes your tasks and fixed time blocks and proposes when each task should
happen, using an earliest-deadline-first algorithm. Right now tasks are
tracked but not yet auto-scheduled — see [Roadmap](#roadmap).

## Tech stack

**Backend**
- Go
- [chi](https://github.com/go-chi/chi) for routing, with `chi/middleware` for logging and panic recovery
- SQLite via [`modernc.org/sqlite`](https://gitlab.com/cznic/sqlite) (pure Go, no CGo dependency)

**Frontend**
- Vue 3 (Composition API)
- Vite
- Tailwind CSS

## Features

- Create tasks with a title, estimated duration, deadline, and priority
- Mark tasks complete
- Add fixed time blocks (meetings, appointments, etc.) for a given day
- Visual daily timeline of time blocks, with day-by-day navigation
- REST API backing both, with a clean separation between models, storage, and HTTP handlers

## Getting started

### Prerequisites
- [Go](https://go.dev/dl/) 1.21+
- [Node.js](https://nodejs.org/) 18+

### Backend

```bash
git clone https://github.com/Cameron-Guinnee/dayplan.git
cd dayplan
go run cmd/server/main.go
```

The API starts on `http://localhost:8080`. A SQLite database file is created
automatically on first run.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

The dev server will print the local URL (typically `http://localhost:5173`).

## API

| Method | Endpoint                  | Description                                  |
|--------|----------------------------|-----------------------------------------------|
| GET    | `/health`                 | Health check                                  |
| POST   | `/tasks`                  | Create a task                                 |
| GET    | `/tasks`                  | List incomplete tasks                         |
| PATCH  | `/tasks/{id}/complete`    | Mark a task complete                          |
| POST   | `/time-blocks`            | Create a time block                           |
| GET    | `/time-blocks?date=YYYY-MM-DD` | List time blocks for a given day (defaults to today) |

## Project structure

```
dayplan/
├── cmd/
│   └── server/         # Go entry point, route wiring, HTTP handlers
├── internal/
│   ├── models/          # Task, TimeBlock types
│   └── store/            # SQLite-backed persistence
├── frontend/
│   └── src/
│       ├── components/  # TaskList, TaskForm, TimeBlockForm, TimelineView
│       ├── api.js        # Fetch wrapper for the Go API
│       └── App.vue        # Layout and state
├── go.mod
└── go.sum
```

## Roadmap

- [ ] Scheduling algorithm (`internal/scheduler`) — earliest-deadline-first, fitting tasks into the gaps between fixed time blocks
- [ ] Render scheduled tasks on the timeline alongside time blocks
- [ ] Priority-weighted scheduling on top of EDF
- [ ] Drag-and-drop manual adjustment of the generated schedule
- [ ] Recurring tasks
- [ ] Export schedule to `.ics` / calendar apps

## License

MIT
