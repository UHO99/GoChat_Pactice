package hub

import (
	"context"
	"gochat/step7_redis/servers/broker"
	"sync"
)

type Room struct {
	id      int64
	name    string
	mu      sync.RWMutex
	clients map[*Client]struct{}
	broker  broker.Broker
	cancel  func()
}

func newRoom(id int64, name string, b broker.Broker, cancel func()) *Room {
	return &Room{
		id:      id,
		name:    name,
		clients: make(map[*Client]struct{}),
		broker:  b,
		cancel:  cancel,
	}
}

func (r *Room) Publish(ctx context.Context, payload []byte) error {
	return r.broker.Publish(ctx, r.name, payload)
}

func (r *Room) deliverLocal(payload []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for c := range r.clients {
		c.send <- payload
	}
}

func (r *Room) Register(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

func (r *Room) Unregister(c *Client) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
	return len(r.clients) == 0
}

func (h *Hub) ReleaselfEmpty(r *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r.mu.RLock()
	empty := len(r.clients) == 0
	r.mu.RUnlock()

	if empty {
		delete(h.rooms, r.name)
		r.cancel()
	}
}
