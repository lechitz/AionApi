package models

type Password struct {
	New     string `json:"new,omitempty"`
	Current string `json:"current,omitempty"`
}
