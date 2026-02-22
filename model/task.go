package model

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
	ID   int       `json:"id"`
}

func ReadNameTask(s string) (string, error) {
	taskName := strings.TrimSuffix(s, "\n")
	return taskName, nil
}

func ReadDateTask(s string) (time.Time, error) {
	timeLayout := "02-01-2006 15:04"

	s = strings.TrimSuffix(s, "\n")

	taskData, errParse := time.Parse(timeLayout, s)
	if errParse != nil {
		return time.Time{}, fmt.Errorf("read date task parse: %w", errParse)
	}
	return taskData, nil
}
