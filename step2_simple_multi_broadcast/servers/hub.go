package servers

import (
	"context"
	"sync"

	"github.com/coder/websocket"
)

type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]struct{})}
}

func (h *Hub) Register(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = struct{}{}
}

func (h *Hub) Unregister(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
}

func (h *Hub) Broadcast(ctx context.Context, payload []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.clients {
		if err := c.Write(context.Background(), websocket.MessageText, payload); err != nil {
			return err
		}
	}

	return nil
}
