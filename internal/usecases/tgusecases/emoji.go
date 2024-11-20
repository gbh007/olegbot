package tgusecases

import (
	"context"
	"math/rand"
)

func (u *UseCases) RandomEmoji(ctx context.Context) (string, bool, error) {
	if rand.Float32() > u.emojiChance || len(u.emojiList) == 0 {
		return "", false, nil
	}

	return (u.emojiList)[rand.Intn(len(u.emojiList))], true, nil
}
