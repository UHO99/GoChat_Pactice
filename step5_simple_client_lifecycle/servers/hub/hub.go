package hub

import (
	"log"
	"sync"
)

type Hub struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]*Room)}
}

func (h *Hub) Room(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	r, ok := h.rooms[name]
	if !ok {
		r = newRoom(name)
		h.rooms[name] = r
		log.Println("방 생성 : " + name)
	}

	log.Println("방 입장 : " + name)
	return r
}
