package render

import (
	"app/internal/domain"
	"time"
)

type Bot struct {
	ID            int64     `json:"id"`
	Enabled       bool      `json:"enabled"`
	EmojiList     []string  `json:"emoji_list,omitempty"`
	EmojiChance   float32   `json:"emoji_chance,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	Name          string    `json:"name"`
	Tag           string    `json:"tag"`
	Description   string    `json:"description,omitempty"`
	Token         string    `json:"token"` // TODO: обезопасить токен при работе с апи.
	AllowedChats  []int64   `json:"allowed_chats,omitempty"`
	StickerChance float32   `json:"sticker_chance,omitempty"`
	GifChance     float32   `json:"gif_chance,omitempty"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at,omitempty"`
}

func BotFromDomain(raw domain.Bot) Bot {
	return Bot{
		ID:            raw.ID,
		Enabled:       raw.Enabled,
		EmojiList:     raw.EmojiList,
		EmojiChance:   raw.EmojiChance,
		Tags:          raw.Tags,
		Name:          raw.Name,
		Tag:           raw.Tag,
		Description:   raw.Description,
		Token:         raw.Token,
		AllowedChats:  raw.AllowedChats,
		StickerChance: raw.StickerChance,
		GifChance:     raw.GifChance,
		CreateAt:      raw.CreateAt,
		UpdateAt:      raw.UpdateAt,
	}
}

func BotToDomain(raw Bot) domain.Bot {
	return domain.Bot{
		ID:            raw.ID,
		Enabled:       raw.Enabled,
		EmojiList:     raw.EmojiList,
		EmojiChance:   raw.EmojiChance,
		Tags:          raw.Tags,
		Name:          raw.Name,
		Tag:           raw.Tag,
		Description:   raw.Description,
		Token:         raw.Token,
		AllowedChats:  raw.AllowedChats,
		StickerChance: raw.StickerChance,
		GifChance:     raw.GifChance,
		CreateAt:      raw.CreateAt,
		UpdateAt:      raw.UpdateAt,
	}
}
