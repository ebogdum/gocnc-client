package tasks

type Template struct {
	Version          int         `json:"version"`
	Source           string      `json:"src"`
	Content          string      `json:"content"`
	Destination      string      `json:"dest"`
	VarFile          string      `json:"var_file"`
	Variables        []Variables `json:"variables"`
	Mode             string      `json:"mode"`
	Owner            string      `json:"owner"`
	Group            string      `json:"group"`
	ModificationTime string      `json:"modification_time"`
	AccessTime       string      `json:"access_time"`
}

type Variables struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
