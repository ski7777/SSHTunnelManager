package config

type Config struct {
	SSHKeys     []SSHKey          `json:"keys,omitempty"`
	Remotes     map[string]Remote `json:"remotes,omitempty"`
	Connections []Connection      `json:"connections,omitempty"`
}
