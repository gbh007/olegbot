package domain

import (
	"math/rand"
	"time"
)

type Bot struct {
	ID int64

	Enabled bool

	EmojiList   []string
	EmojiChance float32
	Tags        []string
	Name        string
	Tag         string
	Description string

	Token        string
	AllowedChats []int64

	CreateAt time.Time
	UpdateAt time.Time
}

func (r Bot) RandomEmoji() (string, bool) {
	if rand.Float32() > r.EmojiChance || len(r.EmojiList) == 0 {
		return "", false
	}

	return (r.EmojiList)[rand.Intn(len(r.EmojiList))], true
}