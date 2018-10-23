package ws

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/coreos/go-systemd/activation"
	"github.com/pkg/errors"

	"github.com/waltton/talktask/config"
)

type server struct {
	ctx context.Context
	cfg config.Server
	srv *http.Server
}

// New creates a new server instance
func New(ctx context.Context, cfg config.Server, handler http.Handler) func() error {
	srv := &http.Server{
		Handler:     handler,
		IdleTimeout: time.Second * 60,
	}

	ws := &server{ctx: ctx, cfg: cfg, srv: srv}

	return ws.run
}

func getListner(cfg config.Server) (net.Listener, error) {
	if cfg.UseSystemdSocket {
		listeners, err := activation.Listeners()
		if err != nil {
			return nil, errors.Wrap(err, "could not get listners from systemd")
		}

		if len(listeners) != 1 {
			return nil, errors.Wrap(err, "got an unexpected number of socket activation")
		}

		log.Printf("Listening on systemd socket: %s\n", listeners[0].Addr())
		return listeners[0], nil
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	listner, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get a tcp listner")
	}
	log.Printf("Listening on: %s\n", addr)

	return listner, err
}

func (s *server) run() error {
	listner, err := getListner(s.cfg)

	if err != nil {
		return errors.Wrap(err, "fail to get any listner")
	}

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
