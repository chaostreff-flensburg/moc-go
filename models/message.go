package models

import (
	"time"
)

type Message struct {
	ID string `json:"id"`

	Text string `json:"message"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewMessage create new message
func NewMessage(text string) *Message {
	return &Message{
		Text: text,
	}
}
