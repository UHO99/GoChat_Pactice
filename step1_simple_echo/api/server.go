package api

import "net/http"

type Options struct {
	Addr string
}

type Server struct {
	httpServer *http.Server
}

func New(opts Options) *Server {
	s := &Server{}

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
