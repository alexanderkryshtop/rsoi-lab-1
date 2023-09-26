package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"rsoi-lab-1/cmd/configuration"
	"rsoi-lab-1/internal/handlers"
	"rsoi-lab-1/internal/repository"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Config *configuration.Config
	Logger *zap.SugaredLogger
	Pool   *pgxpool.Pool
}

func (a *App) constructHandler() *handlers.Handler {
	repo := repository.NewPersonRepository(a.Pool)
	handler := handlers.NewHandler(repo)
	return handler
}

func New(cfg *configuration.Config, logger *zap.SugaredLogger, pool *pgxpool.Pool) *App {
	return &App{
		Config: cfg,
		Logger: logger,
		Pool:   pool,
	}
}

func (a *App) Run() error {
	handler := a.constructHandler()

	server, err := a.newHTTPServer(handler)
	if err != nil {
		return fmt.Errorf("new httpServer: %w", err)
	}

	wg, ctx := errgroup.WithContext(context.Background())
	wg.Go(func() error {
		return server.Run(ctx)
	})

	signalsCh := subscribeSignals()
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

func subscribeSignals() <-chan os.Signal {
	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh, syscall.SIGINT, syscall.SIGTERM)
	return signalsCh
}
