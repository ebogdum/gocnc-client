package main

import (
	"bytes"
	"errors"
	"fmt"
	"gocnc/config"
	"gocnc/providers/tasks"
	taskStructs "gocnc/providers/tasks/structs"
	"gocnc/utils"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	TaskTypeFile taskStructs.TaskType = iota
	TaskTypeTemplate
)

// define global variables that might be used in multiple functions
var conf = config.GetConfig()
var err error

func main() {

	// define local main variables that will be used only in main function
	var serial []byte
	var taskList []taskStructs.TaskList
	var res *http.Response
	var jsonTasks []byte

	// read the tasks.json file and parse it into a list of tasks
	// TODO act as a true client and read these from the server only
	// TODO in this case the config should only contain the server details and some tags to help identify the
	// TODO device or type of client (e.g: role, location, etc)

	jsonTasks, err = os.ReadFile("./providers/tasks.json")
	utils.Check(err)
	taskList = tasks.GetTasks(jsonTasks)

	// Example implementation of the tasks received these should also be moved under providers
	for _, task := range taskList {
		switch task.Type {
		case TaskTypeFile:
			log.Println("File Task:", task.File.Mode) // Example action with the File
			// Access other File properties as needed
		case TaskTypeTemplate:
			log.Println("Template Task:", task.Template.Destination) // Example action with the Template
			// Access other Template properties as needed
		}
	}

	// get the serial number of the device, make the if function, if the file exists, read it, else exit
	// TODO this should be a local taks that is run on the device and the serial number should be read from the device
	// TODO these local tasks should be defined in the tasks.json file and should be run by the client prior to the
	// TODO server tasks
	// TODO the server should be able to define the order of the tasks and the client should be able to run them in the
	//if _, err := os.Stat("/sys/firmware/devicetree/base/serial-number"); os.IsExist(err) {
	serial, err = os.ReadFile("/sys/firmware/devicetree/base/serial-number")
	utils.Check(err)
	serial = bytes.Trim(serial, "\x00")
	//}

	baseURL := conf.Server.Protocol + "://" + conf.Server.Hostname + ":" + strconv.Itoa(conf.Server.Port)
	baseRequestURL := fmt.Sprintf("%s/%s", baseURL, serial)

	newPath := filepath.Join("/", "etc", "nebula.d")
	if _, err = os.Stat(newPath); os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		utils.Check(err)
	}

	// TODO add a new task type called API that will allow the client to make requests to the server
	// TODO add names to the tasks so that I can identify them in the logs and implement conditional tasks
	// TODO add a new task type called command that will allow the client to run commands on the device
	// TODO add a new task type called check that will allow to perform checks and evaluate the results
	// TODO -- this should have a standard list of checks that can be performed and the results should be
	// TODO -- evaluated and the tasks should be run based on the results, e.g: file exists, file contains,
	// TODO -- file does not exist, file does not contain, api status code, api response, previous tasks status
	// TODO add a new task type that can get Variables from different places in order to use them in the templates
	
	url := baseRequestURL + "/make"
	res, err = http.Get(url)
	utils.Check(err)

	switch res.StatusCode {
	case http.StatusOK:
		log.Println(res.StatusCode)

		for _, file := range conf.Files {
			err = downloadFile(filepath.Clean(file.Path), baseRequestURL+file.RequestPath, file.Mode)
			utils.Check(err)
		}
		break
	case http.StatusBadRequest:
		log.Println(res.StatusCode)
		log.Println("bad request")
		break
	default:
		log.Println(res.StatusCode)
		log.Println("unknown error")
		break
	}

	runCommand("/usr/bin/systemctl", "daemon-reload")
	runCommand("/usr/bin/systemctl", "enable", "nebula.service")
	runCommand("/usr/bin/systemctl", "restart", "nebula.service")

}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	log.Printf("running: %s", cmd)
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			log.Printf("Exit Status: %d", exitErr.ExitCode())
		}
	}

	stdout, _ := cmd.CombinedOutput()
	log.Printf("output: %s", stdout)
}

func downloadFile(filepath string, url string, mode os.FileMode) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Create the files
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	// Write the body to files
	_, wErr := io.Copy(out, resp.Body)

	err = os.Chmod(filepath, mode)
	utils.Check(err)

	log.Println("downloaded: " + filepath)

	return wErr
}
