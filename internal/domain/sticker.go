package domain

import "time"

type Sticker struct {
	ID              int64
	BotID           int64
	FileID          string
	CreatorID       int64
	CreatedInChatID int64
	CreateAt        time.Time
}
