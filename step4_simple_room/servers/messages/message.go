package messages

import "time"

type Incoming struct {
	Content string `json:"content"`
}

type Outgoing struct {
	Nickname string    `json:"nickname,omitempty"`
	Content  string    `json:"content,omitempty"`
	SentAt   time.Time `json:"sentAt,omitempty"`
}
