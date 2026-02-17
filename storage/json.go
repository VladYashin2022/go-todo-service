package storage

import (
	"cli_todo/model"
	"encoding/json"
	"fmt"
)

// функция для маршалинга json из среза структуры []Task
func CreateJson(t []model.Task) ([]byte, error) {
	jsonTask, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("create json marshal: %w", err)
	}
	return jsonTask, nil
}

// функция для записи json в файл
func JsonWriter(jsonTask []byte) error {
	file, err := WriteFile()
	if err != nil {
		return fmt.Errorf("json writer open file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(jsonTask); err != nil {
		return fmt.Errorf("json writer write: %w", err)
	}
	return nil
}

// функция для обновления json из AllTasks
func JsonUpdate(t []model.Task) error {
	data, err := CreateJson(t)
	if err != nil {
		return err
	}
	err = JsonWriter(data)
	return err
}

// функция для записи всех задач из файла json. вызывается в начале при каждом запуске программы
func AllTasksWriter() ([]model.Task, error) {
	data, err := ReadFile()
	if err != nil {
		return nil, fmt.Errorf("all task writer read file: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	var allTasksFromJson []model.Task

	err = json.Unmarshal(data, &allTasksFromJson)
	if err != nil {
		return nil, fmt.Errorf("all task writer unmarshal: %w", err)
	}

	return allTasksFromJson, nil
}
