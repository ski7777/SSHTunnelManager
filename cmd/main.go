package main

import (
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/ski7777/SSHTunnelManager/internal/config"
	"github.com/ski7777/SSHTunnelManager/internal/connection"
	"github.com/ski7777/SSHTunnelManager/internal/logger"
	"github.com/ski7777/SSHTunnelManager/internal/remote"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	parser := argparse.NewParser("SSHTunnelManager", "tunnels all your services to other machines")
	conffile := parser.File("c", "config", os.O_RDONLY, 0, &argparse.Options{
		Required: true,
	})
	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	l := logger.GetLogger()
	l.Infow("Starting up")
	confbytes, err := ioutil.ReadAll(conffile)
	if err != nil {
		l.Fatal("Failed reading the configuration file", err)
	}
	_ = conffile.Close()
	conf := config.Config{}
	if err := json.Unmarshal(confbytes, &conf); err != nil {
		l.Fatal("Failed parsing the configuration file", err)
	}

	var connections []*connection.Connection
	remotes := map[string]*remote.Remote{}

	for _, c := range conf.Connections {
		for _, d := range c.Destinations {
			connections = append(connections, &connection.Connection{
				Source:      c.Source,
				Destination: d,
				RemoteGetter: func(r string) *ssh.Client {
					return remotes[r].Client
				},
				Logger: l.With("src", c.Source, "dst", d),
			})
		}
	}

	for rn, rc := range conf.Remotes {
		remotes[rn] = &remote.Remote{
			Name:   rn,
			Config: rc,
			Keys:   conf.SSHKeys,
			Logger: l.With("name", rn),
		}
		go remotes[rn].Start(func(rnf string, state int) {
			for _, c := range connections {
				go c.RemoteCallback(rnf, state)
			}
		})
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	l.Infow("Shutting down")
	for _, r := range remotes {
		r.Stop = true
		_ = r.Client.Close()
	}
}
