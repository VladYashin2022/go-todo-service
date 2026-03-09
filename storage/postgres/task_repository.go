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
		RETURNING id;`,
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

	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrTaskNotFound
	}

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetAllTasks() ([]model.Task, error) {
	var tasks []model.Task

	rows, err := r.db.Query(
		`SELECT id, name, date
		FROM tasks
		ORDER BY id;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Name, &task.Date)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil

}

func (r *TaskRepository) DeleteTask(id int) error {
	result, err := r.db.Exec(
		`DELETE FROM tasks
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) UpdateTask(id int, name string, date time.Time) (model.Task, error) {
	var task model.Task

	err := r.db.QueryRow(
		`UPDATE tasks
		SET name = $1, date = $2
		WHERE id = $3
		RETURNING id, name, date;`,
		name,
		date,
		id,
	).Scan(&task.ID, &task.Name, &task.Date)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) PatchTask(
	id int,
	name *string,
	date *time.Time,
) (model.Task, error) {

	var task model.Task

	err := r.db.QueryRow(
		`UPDATE tasks
		SET name = COALESCE($1, name),
			date = COALESCE($2, date)
		WHERE id = $3
		RETURNING id, name, date;`,
		name,
		date,
		id,
	).Scan(&task.ID, &task.Name, &task.Date)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}
