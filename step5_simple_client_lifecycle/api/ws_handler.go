package api

import (
	"gochat/step5_simple_client_lifecycle/servers/hub"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	roomName := r.URL.Query().Get("room")
	if roomName == "" || nickname == "" {
		http.Error(w, "room and nickname query empty", http.StatusInternalServerError)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Fatal("Cannot Accept Server : ", err)
		return
	}
	defer conn.CloseNow()

	client := hub.NewClient(conn, s.hub.Room(roomName), nickname)

	ctx := r.Context()

	client.Run(ctx)
}
