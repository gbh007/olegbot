package domain

import "time"

type Moderator struct {
	UserID      int64
	CreateAt    time.Time
	Description string
}
