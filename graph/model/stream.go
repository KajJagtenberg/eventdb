package model

type Stream struct {
	Name   string   `json:"name"`
	Events []string `json:"events"`
	Size   int      `json:"size"`
}
