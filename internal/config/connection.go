package config

type Connection struct {
	Source       Endpoint   `json:"source,omitempty"`
	Destinations []Endpoint `json:"destinations,omitempty"`
}
