package api

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	s.hub.Register(conn)
	defer s.hub.Unregister(conn)

	ctx := r.Context()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		log.Printf("[%s] %s", r.RemoteAddr, string(data))

		if err := s.hub.Broadcast(ctx, data); err != nil {
			log.Fatal("broadcast error : ", err)
			return
		}
	}
}
