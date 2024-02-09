package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	var err error
	var serial []byte

	if _, err := os.Stat("/sys/firmware/devicetree/base/serial-number"); os.IsNotExist(err) {
		serial, err = os.ReadFile("/sys/firmware/devicetree/base/serial-number")
		check(err)
		serial = bytes.Trim(serial, "\x00")
	}

	newPath := filepath.Join("/", "etc", "nebula.d")
	if _, err = os.Stat(newPath); os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		check(err)
	}

	url := fmt.Sprintf("http://lnxcode.org:3333/%s/make", serial)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	switch res.StatusCode {
	case http.StatusOK:
		fmt.Println(res.StatusCode)

		fileUrlCA := fmt.Sprintf("http://lnxcode.org:3333/%s/ca", serial)
		err = downloadFile(filepath.Join("/", "etc", "nebula.d", string(serial)+".ccrt"), fileUrlCA)
		check(err)
		fmt.Println("Downloaded_CA: " + fileUrlCA)

		fileUrlCert := fmt.Sprintf("http://lnxcode.org:3333/%s/cert", serial)
		err = downloadFile(filepath.Join("/", "etc", "nebula.d", string(serial)+".crt"), fileUrlCert)
		check(err)
		fmt.Println("Downloaded_Cert: " + fileUrlCert)

		fileUrlKey := fmt.Sprintf("http://lnxcode.org:3333/%s/key", serial)
		err = downloadFile(filepath.Join("/", "etc", "nebula.d", string(serial)+".key"), fileUrlKey)
		check(err)
		fmt.Println("Downloaded_Key: " + fileUrlKey)

		fileUrlConfig := fmt.Sprintf("http://lnxcode.org:3333/%s/config", serial)
		err = downloadFile(filepath.Join("/", "etc", "nebula.d", "config.yml"), fileUrlConfig)
		check(err)
		fmt.Println("Downloaded_Config: " + fileUrlConfig)

		fileUrlService := fmt.Sprintf("http://lnxcode.org:3333/%s/service", serial)
		err = downloadFile(filepath.Join("/", "lib", "systemd", "system", "nebula.service"), fileUrlService)
		check(err)
		fmt.Println("Downloaded_Service: " + fileUrlService)

		break
	case http.StatusBadRequest:
		fmt.Println(res.StatusCode)
		fmt.Println("bad request")
		break
	default:
		fmt.Println(res.StatusCode)
		fmt.Println("unknown error")
		break
	}

	runCommand("/usr/bin/systemctl", "daemon-reload")
	runCommand("/usr/bin/systemctl", "enable", "nebula.service")
	runCommand("/usr/bin/systemctl", "restart", "nebula.service")

}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	fmt.Println(cmd)
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
	fmt.Println(stdout)
}

func check(e error) {
	if e != nil {
		fmt.Println("Error: " + e.Error())
		panic(e)
	}
}

func downloadFile(filepath string, url string) error {

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
	_, err = io.Copy(out, resp.Body)
	return err
}
