package config

type Remote struct {
	Addr     string `json:"addr,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}
