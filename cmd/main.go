package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"nipple/internal/config"
	server "nipple/internal/http"
	"nipple/internal/logger"
	"nipple/internal/provider"
	"nipple/internal/rcon"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		log.Fatalf("could not load config: %s", err)
	}

	lg := logger.New(cfg.Logger)

	rconClient, err := rcon.NewClient(cfg.RCON, lg)
	if err != nil {
		lg.Fatalf("could not establish rcon connection: %s", err)
	}

	prov := provider.New(cfg, rconClient, lg)
	server := server.New(prov)

	errChan := make(chan error)
	go func() {
		err := server.Run(ctx)
		if !errors.Is(err, http.ErrServerClosed) {
			lg.Errorf("server.Run: %s", err)
		}
		errChan <- err
	}()
	lg.Infof("http server started at port %s", server.Info())

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		lg.Infof("Server error: %s", err)
	case sig := <-exitChan:
		lg.Infof("Received shutdown signal: %s", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		lg.Errorf("Application shutdown error: %s", err)
		return
	}
	lg.Infof("server stopped")

	if err := rconClient.Close(); err != nil {
		lg.Errorf("RCON client close error: %s", err)
	}
	lg.Infof("RCON client closed")
}
