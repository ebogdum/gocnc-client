package tasks

import (
	"encoding/json"
	structs "gocnc/providers/tasks/structs"
	"gocnc/utils"
	"log"
)

const (
	TaskTypeFile structs.TaskType = iota
	TaskTypeTemplate
)

func GetTasks(jsonData []byte) []structs.TaskList {
	var err error
	var tasks structs.Tasks
	var orderedTasks []structs.TaskList

	err = json.Unmarshal(jsonData, &tasks)
	utils.Check(err)

	for _, t := range tasks.Task {
		switch t.Module {
		case "files":
			var file structs.File
			err = json.Unmarshal(t.Task, &file)
			utils.Check(err)
			orderedTasks = append(orderedTasks, structs.TaskList{Type: TaskTypeFile, File: &file})
		case "templates":
			var template structs.Template
			err = json.Unmarshal(t.Task, &template)
			utils.Check(err)
			orderedTasks = append(orderedTasks, structs.TaskList{Type: TaskTypeTemplate, Template: &template})
		default:
			log.Println("Unknown module type:", t.Module)
		}
	}

	return orderedTasks

	// Example usage of the ordered tasks
	//for _, task := range orderedTasks {
	//	switch task.Type {
	//	case TaskTypeFile:
	//		log.Println("File Task:", task.File.Mode) // Example action with the File
	//		// Access other File properties as needed
	//	case TaskTypeTemplate:
	//		log.Println("Template Task:", task.Template.Destination) // Example action with the Template
	//		// Access other Template properties as needed
	//	}
	//}
}

//
//func main() {
//
//	var taskList []structs.TaskList
//
//	f, err := os.ReadFile("../tasks.json")
//	utils.Check(err)
//	taskList = GetTasks(f)
//
//	//Example usage of the ordered tasks
//	for _, task := range taskList {
//		switch task.Type {
//		case TaskTypeFile:
//			log.Println("File Task:", task.File.Mode) // Example action with the File
//			// Access other File properties as needed
//		case TaskTypeTemplate:
//			log.Println("Template Task:", task.Template.Destination) // Example action with the Template
//			// Access other Template properties as needed
//		}
//	}
//}
