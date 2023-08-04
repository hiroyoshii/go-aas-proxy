package main

import (
	"context"
	"hiroyoshii/go-aas-proxy/internal/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env"
)

type config struct {
	Port string `env:"HTTP_PORT" envDefault:":8080"`
}

func main() {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		log.Printf("%+v\n", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		defer cancel()
		s := <-signals
		log.Printf("terminated by %s signal\n", s.String())

	}()

	wg := &sync.WaitGroup{}
	server, err := server.NewServer(ctx)
	if err != nil {
		panic(err)
	}

	if err := server.Start(cfg.Port); err != nil {
		panic(err)
	}
	defer server.Shutdown(ctx)
	log.Println("server started")
	wg.Add(1)
	go func() {
		<-ctx.Done()
		wg.Done()
	}()

	wg.Wait()
}
