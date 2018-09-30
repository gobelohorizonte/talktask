package ws

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	"github.com/pkg/errors"
)

// Config contains parameters to the webserver
type Config struct {
	Host             string
	Port             string
	UseSystemdSocket bool
}

type server struct {
	ctx context.Context
	cfg Config
	srv *http.Server
}

// New creates a new server instance
func New(ctx context.Context, cfg Config, handler http.Handler) func() error {
	srv := &http.Server{
		Handler:     handler,
		IdleTimeout: time.Second * 60,
	}

	ws := &server{ctx: ctx, cfg: cfg, srv: srv}

	return ws.run
}

func getListner(cfg Config) (listner net.Listener, err error) {
	if cfg.UseSystemdSocket {
		listeners, err := activation.Listeners()
		if err != nil {
			return nil, errors.Wrap(err, "could not get listners from systemd")
		}

		if len(listeners) != 1 {
			return nil, errors.Wrap(err, "got an unexpected number of socket activation")
		}

		listner = listeners[0]
		log.Printf("Listening on systemd socket: %s\n", listner.Addr())
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if listner, err = net.Listen("tcp", addr); err != nil {
		return nil, errors.Wrap(err, "fail to get a tcp listner")
	}
	log.Printf("Listening on: %s\n", addr)

	return
}

func (s *server) run() error {
	listner, err := getListner(s.cfg)

	if err != nil {
		return errors.Wrap(err, "fail to get any listner")
	}

	// Notify that is ready when running under systemd, not necessary systemd socket
	daemon.SdNotify(false, "READY=1")

	done := make(chan error)

	go func() {
		<-s.ctx.Done()
		s.srv.Shutdown(s.ctx)
		done <- s.ctx.Err()
	}()

	go func() {
		done <- s.srv.Serve(listner)
	}()

	return <-done
}
