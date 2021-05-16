package main

import (
	concurrency "concurrency/src"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	pipeline := &concurrency.PipeLine{
		Group: &errgroup.Group{},
	}
	pipeline.Init()

	pipeline.RunServer("app", ":3001", func(mux *http.ServeMux, svcName string) {
		mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw, "respone from %s", svcName)
		})
	}, func(svcName string) {
		fmt.Printf("%s cleaning when shutdown \n", svcName)
	})

	pipeline.RunServer("agent", ":3002", func(mux *http.ServeMux, svcName string) {
		mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw, "respone from %s", svcName)
		})
	}, func(svcName string) {
		fmt.Printf("%s cleaning when shutdown \n", svcName)
	})

	pipeline.Go(func() error {

		duration := 5 * time.Second
		time.Sleep(duration)
		return fmt.Errorf("err after %s seconds", duration)
	})

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigs:
		pipeline.Shutdown()
		err := pipeline.Wait(func() {
			fmt.Printf("exit with signal %s\n", sig)
		})
		if err != nil {
			fmt.Print(err.Error())
		}
	case <-pipeline.Context.Done():
		err := pipeline.Wait(func() {
			fmt.Printf("exit...\n")
		})
		if err != nil {
			fmt.Print(err.Error())
		}

	}

}
