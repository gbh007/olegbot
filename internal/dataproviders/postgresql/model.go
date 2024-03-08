package postgresql

import (
	"app/internal/domain"
	"database/sql"
	"time"
)

type quoteModel struct {
	ID       int64         `db:"id"`
	Text     string        `db:"text"`
	CreateAt time.Time     `db:"create_at"`
	UserID   sql.NullInt64 `db:"user_id"`
	ChatID   sql.NullInt64 `db:"chat_id"`
}

func (v quoteModel) toDomain() domain.Quote {
	return domain.Quote{
		ID:              v.ID,
		Text:            v.Text,
		CreatorID:       v.UserID.Int64,
		CreatedInChatID: v.ChatID.Int64,
		CreateAt:        v.CreateAt,
	}
}

type moderatorModel struct {
	UserID      int64          `db:"user_id"`
	CreateAt    time.Time      `db:"create_at"`
	Description sql.NullString `db:"description"`
}

func (v moderatorModel) toDomain() domain.Moderator {
	return domain.Moderator{
		UserID:      v.UserID,
		CreateAt:    v.CreateAt,
		Description: v.Description.String,
	}
}
