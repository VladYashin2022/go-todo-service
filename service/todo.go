package service

import (
	"cli_todo/model"
	"cli_todo/storage"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var TasksId int

var AllTasks []model.Task //Срез всех задач

var ErrNotExist error = errors.New("Не существует задачи с выбранным ID.")

func FindMaxID(tasks []model.Task) int {
	var maxID int
	for _, v := range tasks {
		if maxID < v.ID {
			maxID = v.ID
		}
	}
	return maxID
}

// Функция для создания задачи
func CreateTask(s, t string) (model.Task, error) {
	taskName, err := model.ReadNameTask(s)
	if err != nil {
		return model.Task{}, err
	}
	taskDate, err := model.ReadDateTask(t)
	if err != nil {
		return model.Task{}, err
	}
	TasksId += 1

	var task model.Task = model.Task{Name: taskName, Date: taskDate, ID: TasksId}

	AllTasks = append(AllTasks, task)

	err = storage.JsonUpdate(AllTasks)
	if err != nil {
		return model.Task{}, fmt.Errorf("create task update: %w", err)
	}

	return task, nil
}

// функция для чтения отдельной задачи по ID
func ReadTask(id int, t []model.Task) ([]byte, error) {
	for _, v := range t {
		if int(v.ID) == id {
			result, err := json.Marshal(v.ID)
			if err != nil {
				return []byte{}, err
			}
			return result, nil
		}
	}
	return []byte{}, ErrNotExist
}

// функция ищет возвращает объект model.Task
func FindTask(id int, t []model.Task) (model.Task, error) {
	for _, v := range t {
		if int(v.ID) == id {
			return v, nil
		}
	}
	return model.Task{}, ErrNotExist
}

// функция возвращает json по id
func FindTaskJson(id int, t []model.Task) ([]byte, error) {
	for _, v := range t {
		if int(v.ID) == id {
			result, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}
	return nil, ErrNotExist
}

// функция для чтения всех задач
func ReadAllTasks(a []model.Task) string {
	if len(a) == 0 {
		return "В программе пока нет задач для чтения."
	}

	var builder strings.Builder

	for _, v := range a {
		builder.WriteString(stringTask(v) + "\n")
	}
	result := builder.String()
	return result
}

// функция для внесения новых данных в задачу
// if choose 1 - update only name
// if choose 2 - update only date
// if choose 3 - updaye name and date
func UpdateTask(choose, id int, s, d string, t *[]model.Task) error {
	var err error
	switch choose {
	case 1:
		err = UpdateName(id, s, t)
	case 2:
		err = UpdateDate(id, d, t)
	case 3:
		err = UpdateAllTask(id, s, d, t)
	}

	if err != nil {
		return err
	}

	err = storage.JsonUpdate(*t) //обновляем json если нет ошибок. можно потом перенести в main
	return err
}

// функция для обновления имени
func UpdateName(id int, s string, t *[]model.Task) error {
	var err error
	for i := range *t {
		if (*t)[i].ID == id {
			(*t)[i].Name, err = model.ReadNameTask(s)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrNotExist
}

// функция для обновления даты
func UpdateDate(id int, s string, t *[]model.Task) error {
	var err error
	for i := range *t {
		if (*t)[i].ID == id {
			(*t)[i].Date, err = model.ReadDateTask(s)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrNotExist
}

// функция для обновления и имени,  и даты в задаче
func UpdateAllTask(id int, s, d string, t *[]model.Task) error {
	var err error
	for i := range *t {
		if (*t)[i].ID == id {
			(*t)[i].Name, err = model.ReadNameTask(s)
			if err != nil {
				return err
			}

			(*t)[i].Date, err = model.ReadDateTask(d)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrNotExist
}

// функция для удаления задачи из общего среза и storage
func DeleteTask(id int, t *[]model.Task) error {
	for i := range *t {
		if (*t)[i].ID == id {
			(*t) = slices.Delete((*t), i, i+1)
			err := storage.JsonUpdate(*t)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("Не существет задачи с выбранным ID.")
}

// функция для преобразования Задачи в строку
func stringTask(t model.Task) string {
	formatID := strconv.FormatInt(int64(t.ID), 10) //Форматируем int32 в строку
	var builder strings.Builder

	builder.WriteString(formatID + " ")
	builder.WriteString(t.Name + " ")
	builder.WriteString(t.Date.Format("02-01-2006 15:04"))

	result := builder.String()
	return result
}
