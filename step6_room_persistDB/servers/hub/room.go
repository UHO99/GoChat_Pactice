package hub

import (
	"log"
	"sync"
)

type Room struct {
	name    string
	mu      sync.Mutex
	clients map[*Client]struct{}
}

func newRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[*Client]struct{}),
	}
}

func (r *Room) Register(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

func (r *Room) Unregister(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}

func (r *Room) Broadcast(payload []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for c := range r.clients {
		log.Printf("[ROOM : %s / USER : %s] : %s", c.room.name, c.nickname, payload)
		c.send <- payload
	}

	return nil
}
