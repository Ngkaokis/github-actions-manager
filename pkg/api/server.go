package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oursky/github-actions-manager/pkg/github"
	"github.com/oursky/github-actions-manager/pkg/utils/httputil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	logger   *zap.Logger
	server   *http.Server
	regToken *github.RegistrationTokenStore
}

func NewServer(logger *zap.Logger, config *Config, target github.Target, gatherer prometheus.Gatherer) *Server {
	logger = logger.Named("api")

	mux := http.NewServeMux()
	server := &Server{
		logger: logger,
		server: &http.Server{
			Addr:         config.GetAddr(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      mux,
			ErrorLog:     zap.NewStdLog(logger),
		},
		regToken: github.NewRegistrationTokenStore(logger, target),
	}

	mux.Handle("/metrics", promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{
		ErrorLog: zap.NewStdLog(logger.Named("prom")),
	}))

	apiMux := http.NewServeMux()
	mux.Handle("/api/", httputil.UseKeyAuth(config.SharedKeys, apiMux))

	apiMux.HandleFunc("/api/token", server.token)

	return server
}

func (s *Server) Start(ctx context.Context, g *errgroup.Group) error {
	g.Go(func() error {
		go func() {
			<-ctx.Done()

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			s.server.Shutdown(shutdownCtx)
		}()

		s.logger.Info("starting server", zap.String("addr", s.server.Addr))
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to run server: %w", err)
		}
		return nil
	})
	return nil
}
