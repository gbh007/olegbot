package render

import (
	"app/internal/domain"
	"time"
)

type Moderator struct {
	UserID      int64     `json:"user_id"`
	CreateAt    time.Time `json:"create_at"`
	Description string    `json:"description,omitempty"`
}

func ModeratorFromDomain(raw domain.Moderator) Moderator {
	return Moderator(raw)
}
