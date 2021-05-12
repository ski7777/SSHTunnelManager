package remote

import (
	"github.com/ski7777/SSHTunnelManager/internal/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
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
	Stop      bool
	Logger *zap.SugaredLogger
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
	r.Logger.Infow("Initialized remote")
remoteloop:
	for !r.Stop {
		var err error
		r.Logger.Infow("Connecting remote",)
		r.Client, err = ssh.Dial("tcp", r.Config.Addr, &r.sshconfig)
		if err != nil {
			r.Logger.Infow("Failed connecting remote", "reason",err)
			continue remoteloop
		}
		r.Logger.Infow("Connected remote",)
		go cb(r.Name, StateUp)
		err = r.Client.Wait()
		r.Logger.Infow("Disconnected remote", "reason", err)
		go cb(r.Name, StateDown)
	}
}
