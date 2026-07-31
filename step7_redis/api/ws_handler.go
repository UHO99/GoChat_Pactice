package api

import (
	"errors"
	"gochat/step7_redis/servers/hub"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/jackc/pgx/v5"
)

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	roomName := r.URL.Query().Get("room")
	if roomName == "" || nickname == "" {
		http.Error(w, "room and nickname query empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	user, err := s.store.GetUser(ctx, nickname)
	if err != nil {
		http.Error(w, "user get fail", http.StatusNotFound)
		return
	}

	room, err := s.store.GetRoomByName(ctx, roomName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "cannot look up room", http.StatusInternalServerError)
		return
	}

	rm, err := s.hub.Room(ctx, room.ID, room.Name)
	if err != nil {
		http.Error(w, "cannot join room", http.StatusInternalServerError)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Cannot Accept Server : ", err)
		return
	}
	defer conn.CloseNow()

	client := hub.NewClient(conn, rm, nickname, s.store, user.ID)

	client.Run(ctx)

	s.hub.ReleaselfEmpty(rm)
}
