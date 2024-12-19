package domain

import "time"

type Quote struct {
	ID              int64
	BotID           int64
	Text            string
	CreatorID       int64
	CreatedInChatID int64
	CreateAt        time.Time
}
