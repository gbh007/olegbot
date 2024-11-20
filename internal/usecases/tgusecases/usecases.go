package tgusecases

import "context"

type repository interface {
	RandomQuote(ctx context.Context) (string, error)
	AddQuote(ctx context.Context, text string, userID, chatID int64) error
	IsModerator(ctx context.Context, userID int64) (bool, error)
	QuoteExists(ctx context.Context, text string) (bool, error)
}

type UseCases struct {
	repo        repository
	emojiList   []string
	emojiChance float32
}

func New(
	repo repository,
	emojiList []string,
	emojiChance float32,
) *UseCases {
	return &UseCases{
		repo:        repo,
		emojiList:   emojiList,
		emojiChance: emojiChance,
	}
}
