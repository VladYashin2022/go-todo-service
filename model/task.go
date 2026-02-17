package model

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	Name string
	Date time.Time
	ID   int
}

// func (t Task) MarshalJSON() ([]byte, error) {
// 	timeLayout := "02-01-2006 15:04"
// 	var newDate string
// 	newDate = t.Date.Format(timeLayout)

// 	customStruct := map[string]any{
// 		"name": t.Name,
// 		"date": newDate,
// 		"id":   t.ID,
// 	}
// 	return json.Marshal(customStruct)
// }

func ReadNameTask(s string) (string, error) { //убрать reader, перенести в main
	taskName := strings.TrimSuffix(s, "\n")
	return taskName, nil
}

func ReadDateTask(s string) (time.Time, error) { //убрать reader, перенести в main
	timeLayout := "02-01-2006 15:04"

	s = strings.TrimSuffix(s, "\n")

	taskData, errParse := time.Parse(timeLayout, s)
	if errParse != nil {
		return time.Time{}, fmt.Errorf("read date task parse: %w", errParse)
	}
	return taskData, nil
}
