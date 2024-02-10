package tasks

type File struct {
	Version          int    `json:"version"`
	Type             string `json:"type"`
	Source           string `json:"src"`
	Content          string `json:"content"`
	Destination      string `json:"dest"`
	Use              string `json:"use"`
	Mode             string `json:"mode"`
	Owner            string `json:"owner"`
	Group            string `json:"group"`
	ModificationTime string `json:"modification_time"`
	AccessTime       string `json:"access_time"`
}
