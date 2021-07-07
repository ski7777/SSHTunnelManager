package config

import "strings"
import "github.com/thoas/go-funk"

type Endpoint struct {
	Addr   string `json:"addr,omitempty"`
	Type   string `json:"type,omitempty"`
	Remote string `json:"remote,omitempty"`
}

func (e Endpoint) String() string {
	return strings.Join(funk.FilterString([]string{e.Type, e.Remote, e.Addr}, func(s string) bool { return s != "" }), ":")
}
