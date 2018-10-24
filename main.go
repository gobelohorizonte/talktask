package main

import (
	"context"
	"log"

	"github.com/coreos/go-systemd/daemon"
	"github.com/pkg/errors"

	"github.com/waltton/talktask/acd"
	"github.com/waltton/talktask/config"
	"github.com/waltton/talktask/handler"
	"github.com/waltton/talktask/manager"
	"github.com/waltton/talktask/ws"
)

func main() {
	sm := manager.New(context.Background())

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not load config"))
	}

	jobs := make(chan acd.Job, cfg.Options.ACD.QueueSize)
	runACD := acd.New(sm.Context, cfg.Options.ACD.PoolSize, jobs)

	runWebServer := ws.New(sm.Context, cfg.Server, handler.New(jobs))

	sm.CheckSigToQuit()
	sm.RunServiceFunc(runWebServer)
	sm.RunServiceFunc(runACD)

	daemon.SdNotify(false, "READY=1") // Notify that is ready when running under systemd, not necessarily systemd socket

	sm.WaitServices()
}
