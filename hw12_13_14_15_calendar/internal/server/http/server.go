package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/app"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
)

type Server struct {
	server  *http.Server
	log     *logger.Logger
	address string
}

func NewServer(cfg *config.Config, log *logger.Logger, calendar app.Calendar) *Server {
	api := API{
		Log:      log,
		Calendar: calendar,
	}

	router := mux.NewRouter()
	router.HandleFunc("/hello", api.Hello).Methods("GET")

	server := http.Server{ // nolint: exhaustivestruct
		Addr:         net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      loggingMiddleware(router, log),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		server:  &server,
		log:     log,
		address: net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info(fmt.Sprintf("http server has been run on %s", s.address))

	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error occurred on attempt to run server: %w", err)
	}

	err = s.Stop(ctx)
	if err != nil {
		return fmt.Errorf("error occurred on attempt to stop server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error on attempt to stop the server: %w", err)
	}

	return nil
}
