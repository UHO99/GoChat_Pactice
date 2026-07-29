package messages

import "time"

type Type string

const (
	TypeMessage Type = "message"
	TypeJoin    Type = "join"
	TypeLeave   Type = "leave"
)

type Incoming struct {
	Content string `json:"content"`
}

type Outgoing struct {
	Type     Type      `json:"type"`
	Room     string    `json:"room"`
	Nickname string    `json:"nickname,omitempty"`
	Content  string    `json:"content,omitempty"`
	SentAt   time.Time `json:"sentAt,omitempty"`
}
