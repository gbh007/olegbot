package domain

import (
	"math/rand"
)

type Bot struct {
	ID int64

	Enabled bool

	EmojiList   []string
	EmojiChance float32
	Tags        []string
	Name        string // TODO: подумать над надобностью использовать
	Tag         string // TODO: подумать над надобностью использовать

	Token        string
	AllowedChats []int64
}

func (r Bot) RandomEmoji() (string, bool) {
	if rand.Float32() > r.EmojiChance || len(r.EmojiList) == 0 {
		return "", false
	}

	return (r.EmojiList)[rand.Intn(len(r.EmojiList))], true
}
