package hub

import (
	"context"
	"gochat/step7_redis/servers/broker"
	"sync"
)

type Hub struct {
	mu     sync.Mutex
	rooms  map[string]*Room
	broker broker.Broker
}

func NewHub(b broker.Broker) *Hub {
	return &Hub{rooms: make(map[string]*Room), broker: b}
}

func (h *Hub) Room(ctx context.Context, id int64, name string) (*Room, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if r, ok := h.rooms[name]; ok {
		return r, nil
	}

	ch, cancel, err := h.broker.Subscribe(ctx, name)
	if err != nil {
		return nil, err
	}

	r := newRoom(id, name, h.broker, cancel)
	h.rooms[name] = r

	go func() {
		for payload := range ch {
			r.deliverLocal(payload)
		}
	}()

	return r, nil
}
