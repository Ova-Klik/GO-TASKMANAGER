# TaskManager — Milestone 1

Minimal TaskManager REST API implementing Milestone 1: in-memory GET/POST task endpoints using Gin.

Run:

```bash
go mod tidy
go run main.go
```

The server listens on port 8080.

Examples:

Create a task:

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"title":"Buy milk","description":"2 liters"}' \
  http://localhost:8080/api/v1/tasks
```

List tasks:

```bash
curl http://localhost:8080/api/v1/tasks
```

Get a task by id:

```bash
curl http://localhost:8080/api/v1/tasks/1
```
