package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"hiroyoshii/go-aas-proxy/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		defer cancel()
		s := <-signals
		log.Println("hello", s.String())

	}()

	wg := &sync.WaitGroup{}
	server, err := server.NewServer(ctx)
	if err != nil {
		panic(err)
	}

	if err := server.Serve(); err != nil {
		panic(err)
	}
	defer server.GracefulStop()
	log.Println("server started")
	wg.Add(1)
	go func() {
		<-ctx.Done()
		wg.Done()
	}()

	wg.Wait()
}
