package templates

import (
	"os"
	"path/filepath"
)

//{
//      "version": "1",
//      "module": "templates",
//      "template": [
//        {
//          "src": "path to file or URL //this will ignore content -- if both are present, source takes precedence",
//          "content": "content of file //this will ignore source -- if both are present, source takes precedence",
//          "var_file": "path to file or URL",
//          "variables": [{
//            "key": "variable_name",
//            "value": "variable_value"
//          }],
//          "dest": "path to file OR URL",
//          "mode": "mode of file",
//          "owner": "owner of file",
//          "group": "group of file",
//          "modification_time": "modification time of file",
//          "access_time": "access time of file"
//        }
//      ]
//    }

type Templates struct {
	Version  int `json:"version"`
	Template []Template
}

type Template struct {
	Source           string      `json:"src"`
	Content          string      `json:"content"`
	Destination      string      `json:"dest"`
	VarFile          string      `json:"var_file"`
	Variables        []Variables `json:"variables"`
	Mode             os.FileMode `json:"mode"`
	Owner            string      `json:"owner"`
	Group            string      `json:"group"`
	ModificationTime string      `json:"modification_time"`
	AccessTime       string      `json:"access_time"`
}

type Variables struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (f *Template) Src() string {
	return filepath.Clean(f.Source)
}

func (f *Template) Dst() string {
	return filepath.Clean(f.Destination)
}
