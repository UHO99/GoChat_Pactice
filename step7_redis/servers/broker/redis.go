package broker

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	client *redis.Client
}

func NewRedisBroker(client *redis.Client) *RedisBroker {
	return &RedisBroker{client: client}
}

func (b *RedisBroker) Publish(ctx context.Context, roomName string, payload []byte) error {
	return b.client.Publish(ctx, "room:"+roomName, payload).Err()
}

func (b *RedisBroker) Subscribe(ctx context.Context, roomName string) (<-chan []byte, func(), error) {
	ps := b.client.Subscribe(ctx, "room:"+roomName)
	if _, err := ps.Receive(ctx); err != nil {
		return nil, nil, err
	}

	out := make(chan []byte)
	go func() {
		defer close(out)
		for msg := range ps.Channel() {
			out <- []byte(msg.Payload)
		}
	}()

	cancel := func() { ps.Close() }
	return out, cancel, nil
}
