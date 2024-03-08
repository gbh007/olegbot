package render

import (
	"app/internal/domain"
	"time"
)

type Quote struct {
	ID              int64     `json:"id"`
	Text            string    `json:"text"`
	CreatorID       int64     `json:"creator_id,omitempty"`
	CreatedInChatID int64     `json:"created_in_chat_id,omitempty"`
	CreateAt        time.Time `json:"create_at"`
}

func QuoteFromDomain(raw domain.Quote) Quote {
	return Quote(raw)
}
