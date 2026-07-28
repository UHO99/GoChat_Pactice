package api

import (
	"gochat/step3_simple_json_format/servers"
	"net/http"
)

type Options struct {
	Addr string
}

type Server struct {
	httpServer *http.Server
	hub        *servers.Hub
}

func New(opts Options) *Server {
	s := &Server{hub: servers.NewHub()}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws", s.handleWS)

	s.httpServer = &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
	}

	return s
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}
