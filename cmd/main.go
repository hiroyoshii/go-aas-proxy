package main

import (
	"context"
	"hiroyoshii/go-aas-proxy/internal/server"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
)

type config struct {
	Port       string `env:"HTTP_PORT" envDefault:":8081"`
	TLSEnabled bool   `env:"TlsEnabled" envDefault:"false"`
}

func main() {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer cancel()

	server, err := server.NewServer(ctx)
	if err != nil {
		panic(err)
	}
	go func() {
		slog.Info("starting server")
		if cfg.TLSEnabled {
			if err := server.StartAutoTLS(cfg.Port); err != http.ErrServerClosed {
				panic(err)
			}
		} else {
			if err := server.Start(cfg.Port); err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	s := <-signals
	slog.Info("terminated by %s signal\n", s.String())
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
