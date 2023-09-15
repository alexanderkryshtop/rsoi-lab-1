package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type httpServer struct {
	logger          *zap.SugaredLogger
	server          *http.Server
	shutdownTimeout time.Duration
}

func (a *App) newHTTPServer() (*httpServer, error) {
	mux := chi.NewMux()
	mux.Mount("/", adminHandler())
	return &httpServer{
		logger: a.Logger,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", a.Config.HTTPServer.Port),
			Handler: mux,
		},
		shutdownTimeout: a.Config.Application.ShutdownTimeout,
	}, nil
}

func adminHandler() http.Handler {
	mux := chi.NewMux()
	mux.Get("/ping", pingHandler())
	return mux
}

func pingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK\n"))
	}
}

func (s *httpServer) listenAndServe() error {
	s.logger.Infof("listen for http requests: %s\n", s.server.Addr)
	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		s.logger.Infof("http server stopped: %s\n", s.server.Addr)
		return nil
	}
	s.logger.Infof("stop http server: %s\n", err.Error())
	return err
}

func (s *httpServer) gracefulShutdown(ctx context.Context) func() error {
	return func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		s.logger.Infof("stopping http server: %s\n", s.server.Addr)
		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Infof("failed to shutdown http server: %+v\n", err)
		}
		return ctx.Err()
	}
}

func (s *httpServer) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	group.Go(s.listenAndServe)
	group.Go(s.gracefulShutdown(ctx))
	return group.Wait()
}
