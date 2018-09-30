package main

import (
	"context"
	"log"

	"github.com/coreos/go-systemd/daemon"
	"github.com/waltton/talktask/handler"
	"github.com/waltton/talktask/manager"
	"github.com/waltton/talktask/ws"
)

func main() {
	log.Println("Starting talktasks...")

	sm := manager.New(context.Background())

	runWebServer := ws.New(
		sm.Context,
		ws.Config{
			Host:             "0.0.0.0",
			Port:             "4040",
			UseSystemdSocket: false,
		},
		handler.New(),
	)

	sm.CheckSigToQuit()
	sm.RunServiceFunc(runWebServer)

	// Notify that is ready when running under systemd, not necessary systemd socket
	daemon.SdNotify(false, "READY=1")

	sm.WaitServices()
}
