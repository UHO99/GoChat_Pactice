package api

import (
	servers "gochat/step6_room_persistDB/servers/hub"
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
