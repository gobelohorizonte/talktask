package main

import (
	"context"
	"log"
	"os"

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
			Host:             os.Getenv("GOTALK_HOST"),
			Port:             os.Getenv("GOTALK_PORT"),
			UseSystemdSocket: os.Getenv("GOTALK_USE_SYSTEMD_SOCKET") == "true",
		},
		handler.New(),
	)

	sm.CheckSigToQuit()
	sm.RunServiceFunc(runWebServer)

	// Notify that is ready when running under systemd, not necessary systemd socket
	daemon.SdNotify(false, "READY=1")

	sm.WaitServices()
}
