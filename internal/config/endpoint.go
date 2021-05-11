package config

type Endpoint struct {
	Addr   string `json:"addr,omitempty"`
	Type   string `json:"type,omitempty"`
	Remote string `json:"remote,omitempty"`
}
