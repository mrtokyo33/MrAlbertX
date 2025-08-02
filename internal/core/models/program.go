package models

type Program struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Filename string   `json:"filename"`
	Keywords []string `json:"keywords"`
	Tags     []string `json:"tags"`
}
