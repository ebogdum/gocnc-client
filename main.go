package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Config struct {
	Server struct {
		Protocol string `toml:"protocol"`
		Hostname string `toml:"hostname"`
		Port     int    `toml:"port"`
	}
	Files []struct {
		Path        string      `toml:"path"`
		RequestPath string      `toml:"req"`
		Mode        os.FileMode `toml:"mode"`
	}
	Template struct {
		Enabled bool `toml:"enabled"`
	}
}

var conf Config

func main() {
	var err error
	var serial []byte

	_, err = toml.DecodeFile("config.toml", &conf)
	check(err)

	//if _, err := os.Stat("/sys/firmware/devicetree/base/serial-number"); os.IsExist(err) {
	serial, err = os.ReadFile("/sys/firmware/devicetree/base/serial-number")
	check(err)
	serial = bytes.Trim(serial, "\x00")
	//}

	protocol := conf.Server.Protocol
	hostname := conf.Server.Hostname
	port := strconv.Itoa(conf.Server.Port)
	baseURL := protocol + "://" + hostname + ":" + port
	baseRequestURL := fmt.Sprintf("%s/%s", baseURL, serial)

	newPath := filepath.Join("/", "etc", "nebula.d")
	if _, err = os.Stat(newPath); os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		check(err)
	}

	url := baseRequestURL + "/make"
	res, err := http.Get(url)
	check(err)

	switch res.StatusCode {
	case http.StatusOK:
		log.Println(res.StatusCode)

		for _, file := range conf.Files {
			err = downloadFile(filepath.Clean(file.Path), baseRequestURL+file.RequestPath, file.Mode)
			check(err)
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

func check(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Println("Error: " + msg[0] + " -- " + e.Error())
		} else {
			log.Println("Error: " + e.Error())
		}
		panic(e)
	}
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

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	// Write the body to file
	_, wErr := io.Copy(out, resp.Body)

	err = os.Chmod(filepath, mode)
	check(err)

	log.Println("downloaded: " + filepath)

	return wErr
}
