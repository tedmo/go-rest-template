package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tedmo/go-rest-template/internal/app"
)

func (s *Server) HandleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := app.NewLogger(ctx)

		var req app.CreateUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error(err.Error())
			s.badRequest(w)
			return
		}

		user, err := s.UserService.CreateUser(ctx, &req)
		if err != nil {
			logger.Error(err.Error())
			s.internalError(w)
			return
		}

		s.created(w, user)
	}
}

func (s *Server) HandleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := app.NewLogger(ctx)

		id, err := s.pathValueInt64(r, "id")
		if err != nil {
			logger.Error(err.Error())
			s.badRequest(w)
			return
		}

		user, err := s.UserService.FindUserByID(ctx, id)
		if errors.Is(err, app.ErrNotFound) {
			logger.Error(err.Error())
			s.notFound(w)
			return
		}
		if err != nil {
			logger.Error(err.Error())
			s.internalError(w)
			return
		}

		s.ok(w, user)
	}
}

func (s *Server) HandleGetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := app.NewLogger(ctx)

		users, err := s.UserService.FindUsers(ctx)
		if err != nil {
			logger.Error(err.Error())
			s.internalError(w)
			return
		}

		s.ok(w, users)
	}
}
