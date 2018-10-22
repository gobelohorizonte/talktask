package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/coreos/go-systemd/daemon"
	"github.com/pkg/errors"
	"github.com/waltton/talktask/acd"
	"github.com/waltton/talktask/handler"
	"github.com/waltton/talktask/manager"
	"github.com/waltton/talktask/ws"
)

func main() {
	log.Println("Starting talktasks...")

	sm := manager.New(context.Background())

	acdPoolSize, err := strconv.Atoi(os.Getenv("GOTALK_ACD_POOL_SIZE"))
	if err != nil {
		panic(errors.Wrap(err, "got an invalid value for 'GOTALK_ACD_POOL_SIZE'"))
	}

	acdQueueSize, err := strconv.Atoi(os.Getenv("GOTALK_ACD_QUEUEL_SIZE"))
	if err != nil {
		panic(errors.Wrap(err, "got an invalid value for 'GOTALK_ACD_QUEUEL_SIZE'"))
	}

	jobs := make(chan acd.Job, acdQueueSize)
	runACD := acd.New(sm.Context, acdPoolSize, jobs)

	runWebServer := ws.New(
		sm.Context,
		ws.Config{
			Host:             os.Getenv("GOTALK_HOST"),
			Port:             os.Getenv("GOTALK_PORT"),
			UseSystemdSocket: os.Getenv("GOTALK_USE_SYSTEMD_SOCKET") == "true",
		},
		handler.New(jobs),
	)

	sm.CheckSigToQuit()
	sm.RunServiceFunc(runWebServer)
	sm.RunServiceFunc(runACD)

	// Notify that is ready when running under systemd, not necessarily systemd socket
	daemon.SdNotify(false, "READY=1")

	sm.WaitServices()
}
