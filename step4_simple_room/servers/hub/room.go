package hub

import (
	"context"
	"sync"

	"github.com/coder/websocket"
)

type Room struct {
	name    string
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func newRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (r *Room) Register(c *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

func (r *Room) Unregister(c *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}

func (r *Room) Broadcast(ctx context.Context, payload []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for c := range r.clients {
		if err := c.Write(ctx, websocket.MessageText, payload); err != nil {
			return err
		}
	}

	return nil
}
