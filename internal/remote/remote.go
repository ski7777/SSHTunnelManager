package remote

import (
	"github.com/ski7777/SSHTunnelManager/internal/config"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

const (
	StateUp = iota
	StateDown
)

type Remote struct {
	Name      string
	Config    config.Remote
	Keys      []config.SSHKey
	sshconfig ssh.ClientConfig
	Client    *ssh.Client
	stop      bool
}

func (r *Remote) genConfig() {
	var auth []ssh.AuthMethod
	var keys []ssh.Signer
	for _, k := range r.Keys {
		keys = append(keys, k.Signer)
	}
	if len(keys) != 0 {
		auth = append(auth, ssh.PublicKeys(keys...))
	}
	if r.Config.Password != "" {
		auth = append(auth, ssh.Password(r.Config.Password))
	}
	r.sshconfig = ssh.ClientConfig{
		User:            r.Config.User,
		Auth:            auth,
		Timeout:         15 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func (r *Remote) Start(cb func(rn string, state int)) {
	r.genConfig()
remoteloop:
	for !r.stop {
		var err error
		r.Client, err = ssh.Dial("tcp", r.Config.Addr, &r.sshconfig)
		if err != nil {
			log.Println(err)
			continue remoteloop
		}
		go cb(r.Name, StateUp)
		_ = r.Client.Wait()
		go cb(r.Name, StateDown)
	}
}
