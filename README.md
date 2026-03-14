# Go Todo API

REST API для управления задачами, написанный на Go с использованием PostgreSQL и Docker.

## Стек

- Go
- PostgreSQL
- Docker
- net/http
- pgx

## Возможности

- Создание задачи
- Получение списка задач
- Получение задачи по id
- Полное обновление задачи (PUT)
- Частичное обновление (PATCH)
- Удаление задачи

---

# Запуск через Docker

Сначала клонируйте репозиторий:

```bash
git clone https://github.com/VladYashin2022/go-todo-service.git
cd go-todo-service
```

Запустите контейнеры:

```bash
docker compose up --build
```

API будет доступен:

```bash
http://localhost:8080
```

# Docker Hub

Можно скачать готовый образ:

```bash
docker pull vladislavyashin/go-todo-api:latest
```

И запустить:

```bash
docker compose up --build
```

## API

Создать задачу:

```bash
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{"name":"POST test","date":"14-03-2026 11:14"}'
```

Получить все задачи:

```bash
curl http://localhost:8080/tasks
```

Получить задачу по id:

```bash
curl "http://localhost:8080/tasks?id=1"
```

Обновить задачу:

```bash
curl -X PUT "http://localhost:8080/tasks?id=1" \
-H "Content-Type: application/json" \
-d '{"name":"Learn Docker","date":"2026-03-21T19:00:00Z"}'
```

Частичное обновление:

```bash
curl -X PATCH "http://localhost:8080/tasks?id=1" \
-H "Content-Type: application/json" \
-d '{"name":"Master Go"}'
```

Удалить задачу:

```bash
curl -X DELETE "http://localhost:8080/tasks?id=1"
```
