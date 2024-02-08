package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	err := error(nil)

	fmt.Println("Hello, World!")

	dat, _ := os.ReadFile("/sys_host/firmware/devicetree/base/serial-number")
	dat = bytes.Trim(dat, "\x00")
	fmt.Println(string(dat))

	newPath := filepath.Join(".", "/", "etc", "/", "nebula.d")
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		check(err)
	}

	url := fmt.Sprintf("http://localhost:3333/%s/make", dat)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	switch res.StatusCode {
	case http.StatusOK:
		fmt.Println(res.StatusCode)

		fileUrlCA := fmt.Sprintf("http://localhost:3333/%s/ca", dat)
		err = downloadFile(string(dat)+".ccrt", fileUrlCA)
		check(err)
		fmt.Println("Downloaded_CA: " + fileUrlCA)

		fileUrlCert := fmt.Sprintf("http://localhost:3333/%s/cert", dat)
		err = downloadFile(string(dat)+".crt", fileUrlCert)
		check(err)
		fmt.Println("Downloaded_Cert: " + fileUrlCert)

		fileUrlKey := fmt.Sprintf("http://localhost:3333/%s/key", dat)
		err = downloadFile(string(dat)+".key", fileUrlKey)
		check(err)
		fmt.Println("Downloaded_Key: " + fileUrlKey)

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
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
