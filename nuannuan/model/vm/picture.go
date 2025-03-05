package vm

type Picture struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	MyFile File   `json:"file"`
}

type File struct {
	Uid                string `json:"uid"`
	LastModified       string `json:"lastModified"`
	Name               string `json:"name"`
	Size               string `json:"size"`
	Type               string `json:"type"`
	WebkitRelativePath string `json:"webkitRelativePath"`
}
