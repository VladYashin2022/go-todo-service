package postgres

import (
	"cli_todo/model"
	"database/sql"
	"errors"
	"time"
)

var ErrTaskNotFound error = errors.New("task not found")

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(database *sql.DB) *TaskRepository {
	return &TaskRepository{db: database}
}

func (r *TaskRepository) CreateTask(name string, date time.Time) (model.Task, error) {
	var id int

	err := r.db.QueryRow(
		`INSERT INTO tasks (name, date)
		VALUES ($1, $2)
		RETURNING id`,
		name,
		date,
	).Scan(&id) // Sсan сканирует RETURNING id

	if err != nil {
		return model.Task{}, err
	}

	return model.Task{
		ID:   id,
		Name: name,
		Date: date,
	}, nil
}

func (r *TaskRepository) GetTaskByID(id int) (model.Task, error) {
	var task model.Task

	err := r.db.QueryRow(
		`SELECT id, name, date
		FROM tasks
		WHERE id = $1;`,
		id,
	).Scan(&task.ID, &task.Name, &task.Date)

	if err == sql.ErrNoRows {
		return model.Task{}, ErrTaskNotFound
	}

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}
