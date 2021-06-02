package connection

import (
	"github.com/ski7777/SSHTunnelManager/internal/config"
	"github.com/ski7777/SSHTunnelManager/internal/remote"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
)

type Connection struct {
	Source           config.Endpoint
	sourceState      bool
	Destination      config.Endpoint
	destinationState bool
	RemoteGetter     func(r string) *ssh.Client
	dialer           func() (net.Conn, error)
	dest             net.Listener
	srcs             []net.Conn
	Logger           *zap.SugaredLogger
}

func (c *Connection) RemoteCallback(r string, state int) {
	var changed bool
	if c.Source.Remote == r {
		c.sourceState = state == remote.StateUp
		changed = true
	}
	if c.Destination.Remote == r {
		c.destinationState = state == remote.StateUp
		changed = true
	}
	if changed {
		if (c.Source.Remote == "" || c.sourceState) && (c.Destination.Remote == "" || c.destinationState) {
			c.connect()
		} else {
			c.disconnect()
		}
	}
}

func (c *Connection) connect() {
	var err error
	c.Logger.Infow("Enabling connection")
	var dialer func(n, addr string) (net.Conn, error)
	if c.Source.Remote == "" {
		dialer = net.Dial
	} else {
		dialer = c.RemoteGetter(c.Source.Remote).Dial
	}
	c.dialer = func() (net.Conn, error) {
		return dialer(c.Source.Type, c.Source.Addr)
	}

	var listener func(n, addr string) (net.Listener, error)
	if c.Destination.Remote == "" {
		listener = net.Listen
	} else {
		r := c.RemoteGetter(c.Destination.Remote)
		if c.Destination.Type == "unix" {
			s, err := r.NewSession()
			if err != nil {
				return
			}
			if err := s.Run("rm -rf " + c.Destination.Addr); err != nil {
				return
			}
		}
		listener = r.Listen
	}
	c.dest, err = listener(c.Destination.Type, c.Destination.Addr)
	if err != nil {
		c.Logger.Warnw("Listening failed", "reason", err)
		c.close()
		return
	}

runner:
	for {
		client, err := c.dest.Accept()
		if err != nil {
			break runner
		}
		go c.handleClient(client)
	}
}

func (c *Connection) disconnect() {
	c.Logger.Infow("Disabling connection")
	c.close()
}

func (c *Connection) close() {
	if c.dest != nil {
		_ = c.dest.Close()
		c.dest = nil
	}
	for _, s := range c.srcs {
		if s != nil {
			go func() {
				_ = s.Close()
			}()
		}
	}
	c.srcs = nil
}

func (c *Connection) handleClient(client net.Conn) {
	l:=c.Logger.With("client",client.RemoteAddr())
	l.Debugw("Connecting")
	r, err := c.dialer()
	if err != nil {
		c.Logger.Warnw("Dialing failed", "reason", err)
		if r != nil {
			if err := r.Close();err!=nil{
				l.Warnw("Failed closing dial-out connection", "reason", err)
			}
		}
		return
	}
	c.srcs = append(c.srcs, r)

	chDone := make(chan bool)

	go func() {
		if _, err := io.Copy(client, r);err!=nil&&err!=io.EOF{
			l.Warnw("Failed transfering data", "reason", err)
		}
		chDone <- true
	}()

	go func() {
		if _, err := io.Copy( r,client);err!=nil&&err!=io.EOF{
			l.Warnw("Failed transfering data", "reason", err)
		}
		chDone <- true
	}()

	<-chDone
	go func(){
		<-chDone
		close(chDone)
	}()

	_ = client.Close()
	if err := client.Close();err!=nil&&err!=io.EOF{
		l.Warnw("Failed closing dial-in connection", "reason", err)
	}
	if err := r.Close();err!=nil&&err!=io.EOF{
		l.Warnw("Failed closing dial-out connection", "reason", err)
	}
	l.Debugw("Disconnected")
}
