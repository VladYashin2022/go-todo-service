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