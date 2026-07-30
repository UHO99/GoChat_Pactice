package api

import (
	"gochat/db/store"
	"gochat/step6_room_persistDB/servers/hub"
	"net/http"
)

type Options struct {
	Addr string
}

type Server struct {
	httpServer *http.Server
	hub        *hub.Hub
	store      *store.Store
}

func New(h *hub.Hub, st *store.Store, opts Options) *Server {
	s := &Server{
		hub:   h,
		store: st,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", s.handleCreateUser)
	mux.HandleFunc("POST /rooms", s.handleCreateRoom)
	mux.HandleFunc("GET /rooms", s.handleListRooms)
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
