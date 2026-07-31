package broker

import "context"

type Broker interface {
	Publish(ctx context.Context, roomName string, payload []byte) error
	Subscribe(ctx context.Context, roomName string) (<-chan []byte, func(), error)
}
