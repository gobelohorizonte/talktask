package manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
)

// ServiceManager enable the execution of multiple services and do an graceful shutdown when one of them fail
type ServiceManager struct {
	Context       context.Context
	contextCancel context.CancelFunc
	wg            sync.WaitGroup
}

// ServiceFunc is the signature of a serice function
type ServiceFunc func() error

// New return a new instance of ServiceManager
func New(ctx context.Context) (sm *ServiceManager) {
	sm = &ServiceManager{}
	sm.Context, sm.contextCancel = context.WithCancel(ctx)
	return
}

// RunServiceFunc run a managed service function
func (sm *ServiceManager) RunServiceFunc(sf ServiceFunc) {
	sm.wg.Add(1)

	go func() {
		err := sf()

		if err != context.Canceled {
			log.Println("Shutting down services; error:", err)
			sm.contextCancel()
		}

		sm.wg.Done()
	}()
}

// WaitServices to finish
func (sm *ServiceManager) WaitServices() {
	sm.wg.Wait()
}

// CheckSigToQuit checks for a closing signal and then close the manager context
func (sm *ServiceManager) CheckSigToQuit() {
	sm.RunServiceFunc(func() error {
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		select {
		case sig := <-sigquit:
			return fmt.Errorf("caught sig: %+v", sig)
		case <-sm.Context.Done():
			return sm.Context.Err()
		}
	})
}
