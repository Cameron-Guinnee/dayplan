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

The backend runs an earliest-deadline-first scheduler that fits tasks into
the gaps between your fixed commitments, and the frontend renders the result
on a visual daily timeline — see [Roadmap](#roadmap) for what's next.

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
- EDF scheduler fits tasks into the free gaps around your fixed commitments
- Visual daily timeline showing both fixed blocks and scheduled tasks, with day-by-day navigation
- Unscheduled tasks surfaced in an amber banner when the day is too full
- REST API backing all of the above, with a clean separation between models, storage, scheduler, and HTTP handlers

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
| GET    | `/schedule?date=YYYY-MM-DD&day_start=HH:MM&day_end=HH:MM` | Run the EDF scheduler; returns `scheduled` and `unscheduled` arrays |

## Project structure

```
dayplan/
├── cmd/
│   └── server/         # Go entry point, route wiring, HTTP handlers
├── internal/
│   ├── models/          # Task, TimeBlock types
│   ├── scheduler/       # EDF scheduling algorithm (pure Go, no DB/HTTP deps)
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

- [x] Scheduling algorithm (`internal/scheduler`) — earliest-deadline-first, fitting tasks into the gaps between fixed time blocks
- [x] Render scheduled tasks on the timeline alongside time blocks
- [ ] Delete tasks and time blocks
- [ ] Input validation (required fields, sensible bounds)
- [ ] Priority-weighted scheduling on top of EDF
- [ ] Drag-and-drop manual adjustment of the generated schedule
- [ ] Recurring tasks
- [ ] Export schedule to `.ics` / calendar apps

## License

MIT
