package config

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

type SSHKey struct{ ssh.Signer }

func (k *SSHKey) UnmarshalJSON(raw []byte) error {
	data := &struct {
		Path string `json:"path,omitempty"`
		Raw  string `json:"key,omitempty"`
	}{}
	if err := json.Unmarshal(raw, data); err != nil {
		return err
	}
	if data.Raw == "" {
		if data.Path == "" {
			return errors.New("path and key cannot be empty at the same time")
		}
		if rk, err := ioutil.ReadFile(data.Path); err != nil {
			return err
		} else {
			data.Raw = string(rk)
		}
	}
	if signer, err := ssh.ParsePrivateKey([]byte(data.Raw)); err != nil {
		return err
	} else {
		k.Signer = signer
	}
	return nil
}
