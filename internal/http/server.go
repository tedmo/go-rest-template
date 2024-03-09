package http

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tedmo/go-rest-template/internal/app"
)

//go:embed docs/openapi.yaml
var OpenAPISpec []byte

type UserService interface {
	CreateUser(ctx context.Context, user *app.CreateUserReq) (*app.User, error)
	FindUserByID(ctx context.Context, id int64) (*app.User, error)
	FindUsers(ctx context.Context) ([]app.User, error)
}

type Server struct {
	UserService UserService
	Port        int
}

func (s *Server) ListenAndServe() error {
	logger := app.NewLogger(context.Background())
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: s.Routes(),
	}

	// This channel is used to receive any errors from our graceful shutdown
	shutdownErrChan := make(chan error)

	go func() {
		// Intercept exit signal to shut down server gracefully
		exitChan := make(chan os.Signal, 1)
		signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-exitChan

		// Allow 20 seconds for graceful shutdown to complete
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		logger.Info("stopping server", slog.String("signal", sig.String()))
		shutdownErrChan <- server.Shutdown(ctx)
	}()

	logger.Info("starting server", slog.Int("port", s.Port))

	// Shutdown() will cause http.ErrServerClosed error, so that specific error will be ignored
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait to receive the return value from our graceful shutdown
	err = <-shutdownErrChan
	if err != nil {
		return err
	}

	logger.Info("server stopped")
	return nil
}

func (s *Server) pathValueInt64(r *http.Request, name string) (int64, error) {
	val := r.PathValue(name)
	int64Val, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64Val, nil
}

func (s *Server) ok(w http.ResponseWriter, v interface{}) {
	s.json(w, http.StatusOK, v)
}

func (s *Server) created(w http.ResponseWriter, v interface{}) {
	s.json(w, http.StatusCreated, v)
}

type Response[T any] struct {
	Data T `json:"data"`
}

func (s *Server) json(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response[any]{Data: v})
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func (s *Server) badRequest(w http.ResponseWriter) {
	s.error(w, http.StatusBadRequest, "bad request")
}

func (s *Server) notFound(w http.ResponseWriter) {
	s.error(w, http.StatusNotFound, "not found")
}

func (s *Server) internalError(w http.ResponseWriter) {
	s.error(w, http.StatusInternalServerError, "unexpected error")
}

func (s *Server) error(w http.ResponseWriter, status int, msg string) {
	s.json(w, status, ErrorResponse{Errors: []string{msg}})
}
