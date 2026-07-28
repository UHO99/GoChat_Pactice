package api

import (
	"encoding/json"
	"gochat/step4_simple_room/servers/messages"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general"
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Fatal("Cannot Accept Server : ", err)
		return
	}
	defer conn.CloseNow()

	room := s.hub.Room(roomName)

	room.Register(conn)
	defer room.Unregister(conn)

	ctx := r.Context()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			log.Fatal("Cannot Connection Failed : ", err)
			return
		}

		// messages.Incoming으로 Unmarshal-역직렬화
		var in messages.Incoming
		if err := json.Unmarshal(data, &in); err != nil {
			log.Fatal("역직렬화 실패", err)
			continue
		}

		// messages.Outgoing으로 marshal-직렬화
		out := messages.Outgoing{
			Nickname: nickname,
			Content:  in.Content,
			SentAt:   time.Now(),
		}
		payload, err := json.Marshal(out)
		if err != nil {
			log.Fatal("직렬화 실패", err)
			continue
		}

		log.Printf("[%s] %s", nickname, in.Content)

		if err := room.Broadcast(ctx, payload); err != nil {
			log.Fatal("broadcast error : ", err)
			return
		}
	}
}
