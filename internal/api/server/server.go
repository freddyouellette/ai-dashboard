package server

import (
	"net/http"
)

type Server struct {
	handler *http.Handler
	port    int
}

func NewServer(handler *http.Handler, port int) *Server {
	return &Server{
		handler: handler,
		port:    port,
	}
}

func (s *Server) Start() error {
	http.Handle("/", *s.handler)
	return http.ListenAndServe(":8080", nil)
}
