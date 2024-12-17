package postgresql

import (
	"context"
	"math/rand"
)

func (r *Repository) RandomEmoji(ctx context.Context) (string, bool, error) {
	if rand.Float32() > r.emojiChance || len(r.emojiList) == 0 {
		return "", false, nil
	}

	return (r.emojiList)[rand.Intn(len(r.emojiList))], true, nil
}
