package postgresql

import (
	"app/internal/domain"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ pgx.RowScanner = (*botModel)(nil)

type botModel struct {
	ID           int64                    `db:"id"`
	Name         string                   `db:"name"`
	BotTag       string                   `db:"bot_tag"`
	Token        string                   `db:"token"`
	Enabled      bool                     `db:"enabled"`
	Description  sql.NullString           `db:"description"`
	EmojiList    pgtype.FlatArray[string] `db:"emoji_list"`
	EmojiChance  sql.NullFloat64          `db:"emoji_chance"`
	Tags         pgtype.FlatArray[string] `db:"tags"`
	AllowedChats pgtype.FlatArray[int64]  `db:"allowed_chats"`
	CreateAt     time.Time                `db:"create_at"`
	UpdateAt     sql.NullTime             `db:"update_at"`
}

func (v botModel) toDomain() domain.Bot {
	return domain.Bot{
		ID:           v.ID,
		Enabled:      v.Enabled,
		Name:         v.Name,
		Tag:          v.BotTag,
		Description:  v.Description.String,
		Token:        v.Token,
		EmojiList:    v.EmojiList,
		EmojiChance:  float32(v.EmojiChance.Float64),
		Tags:         v.Tags,
		AllowedChats: v.AllowedChats,
		CreateAt:     v.CreateAt,
		UpdateAt:     v.UpdateAt.Time,
	}
}

func (v *botModel) fromDomain(raw domain.Bot) {
	v.ID = raw.ID
	v.Name = raw.Name
	v.BotTag = raw.Tag
	v.Token = raw.Token
	v.Enabled = raw.Enabled
	v.Description = StringToDB(raw.Description)
	v.EmojiList = ArrayToDB(raw.EmojiList)
	v.EmojiChance = sql.NullFloat64{
		Float64: float64(raw.EmojiChance),
		Valid:   raw.EmojiChance > 0,
	}
	v.Tags = ArrayToDB(raw.Tags)
	v.AllowedChats = ArrayToDB(raw.AllowedChats)
	v.CreateAt = raw.CreateAt
	v.UpdateAt = TimeToDB(raw.UpdateAt)
}

func (v botModel) columns() []string {
	return []string{
		"id",
		"name",
		"bot_tag",
		"token",
		"enabled",
		"description",
		"emoji_list",
		"emoji_chance",
		"tags",
		"allowed_chats",
		"create_at",
		"update_at",
	}
}

func (v *botModel) ScanRow(rows pgx.Rows) error {
	return rows.Scan(
		&v.ID,
		&v.Name,
		&v.BotTag,
		&v.Token,
		&v.Enabled,
		&v.Description,
		&v.EmojiList,
		&v.EmojiChance,
		&v.Tags,
		&v.AllowedChats,
		&v.CreateAt,
		&v.UpdateAt,
	)
}

type quoteModel struct {
	ID       int64         `db:"id"`
	BotID    int64         `db:"bot_id"`
	Text     string        `db:"text"`
	CreateAt time.Time     `db:"create_at"`
	UpdateAt sql.NullTime  `db:"update_at"`
	UserID   sql.NullInt64 `db:"user_id"`
	ChatID   sql.NullInt64 `db:"chat_id"`
}

func (v quoteModel) toDomain() domain.Quote {
	return domain.Quote{
		ID:              v.ID,
		BotID:           v.BotID,
		Text:            v.Text,
		CreatorID:       v.UserID.Int64,
		CreatedInChatID: v.ChatID.Int64,
		CreateAt:        v.CreateAt,
	}
}

type moderatorModel struct {
	UserID      int64          `db:"user_id"`
	BotID       int64          `db:"bot_id"`
	CreateAt    time.Time      `db:"create_at"`
	Description sql.NullString `db:"description"`
}

func (v moderatorModel) toDomain() domain.Moderator {
	return domain.Moderator{
		UserID:      v.UserID,
		BotID:       v.BotID,
		CreateAt:    v.CreateAt,
		Description: v.Description.String,
	}
}

func StringToDB(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func TimeToDB(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t.UTC(),
		Valid: !t.IsZero(),
	}
}

func Int64ToDB(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: i != 0,
	}
}

func ArrayToDB[T any](raw []T) pgtype.FlatArray[T] {
	if len(raw) == 0 {
		return nil
	}

	return raw
}
