package hub

import (
	"context"
	"encoding/json"
	"errors"
	"gochat/step5_simple_client_lifecycle/servers/messages"
	"io"
	"log"
	"net"
	"time"

	"github.com/coder/websocket"
)

const (
	sendBuffer   = 16
	writeTimeout = 5 * time.Second
)

type Client struct {
	conn     *websocket.Conn
	room     *Room
	nickname string
	send     chan []byte
}

func NewClient(conn *websocket.Conn, room *Room, nickname string) *Client {
	return &Client{
		conn:     conn,
		room:     room,
		nickname: nickname,
		send:     make(chan []byte, sendBuffer),
	}
}

func (c *Client) Run(ctx context.Context) {
	c.room.Register(c)
	c.room.Broadcast(c.event(messages.TypeJoin, ""))

	done := make(chan struct{})
	go c.writePump(ctx, done)
	c.readPump(ctx)

	c.room.Unregister(c)
	close(c.send)
	<-done

	c.room.Broadcast(c.event(messages.TypeLeave, ""))
}

func (c *Client) readPump(ctx context.Context) {
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			switch {
			case errors.Is(err, io.EOF), errors.Is(err, net.ErrClosed):
				log.Printf("client disconnected : %s", c.nickname)
			case websocket.CloseStatus(err) != -1:
				log.Printf("client close : %v", websocket.CloseStatus(err))
			default:
				log.Printf("read error : %v", err)
			}

			return
		}

		var in messages.Incoming
		if err := json.Unmarshal(data, &in); err != nil || in.Content == "" {
			log.Printf("Cannot Unmarshal : %v", data)
			continue
		}

		c.room.Broadcast(c.event(messages.TypeMessage, in.Content))
	}
}

func (c *Client) writePump(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	for payload := range c.send {
		if err := c.write(ctx, payload); err != nil {
			log.Printf("Cannot Write Buffer : %v", payload)
			return
		}
	}
}

func (c *Client) write(ctx context.Context, payload []byte) error {
	writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	return c.conn.Write(writeCtx, websocket.MessageText, payload)
}

func (c *Client) event(t messages.Type, content string) []byte {
	payload, err := json.Marshal(messages.Outgoing{
		Type:     t,
		Room:     c.room.name,
		Nickname: c.nickname,
		Content:  content,
		SentAt:   time.Now(),
	})
	if err != nil {
		log.Printf("cannot JSON Marshal : %v", err)
		return nil
	}

	return payload
}
