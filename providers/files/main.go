package files

import (
	"encoding/json"
	"fmt"
	"gocnc/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

//{
//      "version": "1",
//      "module": "files",
//      "files": [
//        {
//          "type": "file //type of file: absent, directory, file, hard, symbolic, touch, truncate, download, read, upload",
//          "src": "path to file or URL //this will ignore content -- if both are present, source takes precedence",
//          "content": "content of file //this will ignore source -- if both are present, source takes precedence",
//          "dest": "path to file",
//          "use": "copy or move",
//          "mode": "mode of file",
//          "owner": "owner of file",
//          "group": "group of file",
//          "modification_time": "modification time of file",
//          "access_time": "access time of file"
//        },
//        {
//          "type": "touch //type of file: absent, directory, file, hard, symbolic, touch",
//          "dest": "path to file",
//          "mode": "mode of file",
//          "owner": "owner of file",
//          "group": "group of file",
//          "modification_time": "modification time of file",
//          "access_time": "access time of file"
//        }
//      ]
//    }

type Files struct {
	Version int `json:"version"`
	Files   []File
}

type File struct {
	Type             string      `json:"type"`
	Source           string      `json:"src"`
	Content          string      `json:"content"`
	Destination      string      `json:"dest"`
	Use              string      `json:"use"`
	Mode             os.FileMode `json:"mode"`
	Owner            string      `json:"owner"`
	Group            string      `json:"group"`
	ModificationTime string      `json:"modification_time"`
	AccessTime       string      `json:"access_time"`
}

func (f *File) Src() string {
	return filepath.Clean(f.Source)
}

func (f *File) Dst() string {
	return filepath.Clean(f.Destination)
}

var files Files

func GetJSON() {
	resp, err := http.Get("https://reqres.in/api/users?page=2")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		utils.Check(err)
	}(resp.Body)
	body, err := io.ReadAll(resp.Body) // response body is []byte

	var result Files
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	// fmt.Println(PrettyPrint(result))
	//
	//// Loop through the data node for the FirstName
	//for _, rec := range result.Data {
	//	fmt.Println(rec.FirstName)
	//}
}
