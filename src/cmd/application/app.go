package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"rsoi-lab-1/cmd/configuration"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Config *configuration.Config
	Logger *zap.SugaredLogger
}

func New(cfg *configuration.Config, logger *zap.SugaredLogger) *App {
	return &App{
		Config: cfg,
		Logger: logger,
	}
}

func (a *App) Run() error {
	server, err := a.newHTTPServer()
	if err != nil {
		return fmt.Errorf("new httpServer: %w", err)
	}
	wg, ctx := errgroup.WithContext(context.Background())
	wg.Go(func() error {
		return server.Run(ctx)
	})
	signalsCh := subscribeSiglans()
	wg.Go(func() error {
		select {
		case sig := <-signalsCh:
			a.Logger.Infof("app received external signal for shutdown %v\n", sig)
			return fmt.Errorf("shutdown by external signal %v", sig)
		case <-ctx.Done():
			return nil
		}
	})
	return wg.Wait()
}

func subscribeSiglans() <-chan os.Signal {
	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh, syscall.SIGINT, syscall.SIGTERM)
	return signalsCh
}
