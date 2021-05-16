package concurrency

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type PipeLine struct {
	Group   *errgroup.Group
	Context context.Context
	wg      sync.WaitGroup
}

func (pipeline *PipeLine) Init() {
	g, ctx := errgroup.WithContext(context.Background())
	pipeline.Group = g
	pipeline.Context = ctx
}

func (PipeLine *PipeLine) Go(f func() error) {
	PipeLine.Group.Go(f)
}

func (PipeLine *PipeLine) listenAndSrv(svcName string, addr string, handler func(mux *http.ServeMux, svcName string), cleanWhenShutdown func(svcName string)) error {
	mux := http.NewServeMux()

	handler(mux, svcName)

	fmt.Printf("Server %s Runing...", svcName)
	fmt.Println()

	timeout := time.Second * 10
	srv := http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       timeout * 3,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout * 6,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	srv.RegisterOnShutdown(func() {
		cleanWhenShutdown(svcName)
	})

	PipeLine.wg.Add(1)
	go func() {
		defer PipeLine.wg.Done()
		<-PipeLine.Context.Done()
		srv.Shutdown(PipeLine.Context)
	}()

	return srv.ListenAndServe()
}

func (PipeLine *PipeLine) RunServer(svcName string, addr string, handler func(mux *http.ServeMux, svcName string), cleanWhenShutdown func(svcName string)) {
	PipeLine.Go(func() error {
		return PipeLine.listenAndSrv(svcName, addr, handler, cleanWhenShutdown)
	})
}

func (pipeline *PipeLine) Shutdown() {
	pipeline.Group.Go(func() error {
		return errors.New("pipeline shutdown")
	})
}

func (pipeline *PipeLine) Wait(cb func()) error {
	err := pipeline.Group.Wait()
	pipeline.wg.Wait()
	cb()
	return err
}
