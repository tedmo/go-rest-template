package http

import (
	"net/http"

	"github.com/flowchartsman/swaggerui"
)

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", s.HandleCreateUser())
	mux.HandleFunc("GET /users", s.HandleGetUsers())
	mux.HandleFunc("GET /users/{id}", s.HandleGetUser())
	mux.Handle("/docs/", http.StripPrefix("/docs", swaggerui.Handler(OpenAPISpec)))

	return mux
}
