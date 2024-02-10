package tasks

import "encoding/json"

type Tasks struct {
	Version int    `json:"version"`
	Task    []Task `json:"tasks"`
}

type Task struct {
	Module string `json:"module"`
	Task   json.RawMessage
}

type TaskType int

type TaskList struct {
	Type     TaskType
	File     *File
	Template *Template
}
