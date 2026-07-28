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

	ctx := r.Context()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		log.Println(string(data))

		if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
			return
		}
	}
}
