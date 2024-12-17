package tgusecases

import (
	"context"
	"strings"
)

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
	tags        []string
}

func New(
	repo repository,
	emojiList []string,
	emojiChance float32,
	rawTags []string,
	botName, botTag string,
) *UseCases {
	tags := make([]string, 0, len(rawTags)+2) // FIXME: в репозиторий

	if botName != "" {
		tags = append(tags, strings.ToLower(botName))
	}

	if botTag != "" {
		tags = append(tags, strings.ToLower(botTag))
	}

	for _, tag := range rawTags {
		if tag != "" {
			tags = append(tags, strings.ToLower(tag))
		}
	}

	return &UseCases{
		repo:        repo,
		emojiList:   emojiList,
		emojiChance: emojiChance,
		tags:        tags,
	}
}
