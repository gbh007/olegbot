package tgusecases

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"app/internal/domain"

	"github.com/samber/lo"
)

func (u *UseCases) randomQuote(ctx context.Context, botInfo *domain.Bot, msgHistory []string, useLLM bool) (string, error) {
	quotes, err := u.repo.Quotes(ctx, u.botID)
	if err != nil {
		return "", fmt.Errorf("use case: random quote: %w", err)
	}

	if len(quotes) == 0 {
		return "", fmt.Errorf("use case: random quote: no quotes")
	}

	if u.llm != nil && botInfo != nil && useLLM {
		tStart := time.Now()
		q, err := u.llmQuote(ctx, botInfo, quotes, msgHistory)
		if err != nil {
			u.logger.ErrorContext(ctx, "get llm quote", "err", err.Error())
		} else {
			u.logger.DebugContext(ctx, "use llm quote", "quote", q, "duration", time.Since(tStart).String())
			return q, nil
		}
	}

	return quotes[rand.Intn(len(quotes))].Text, nil
}

func (u UseCases) llmQuote(ctx context.Context, botInfo *domain.Bot, quotes []domain.Quote, msgHistory []string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	q, err := u.llm.GetQuote(ctx, botInfo.Tags, lo.Map(quotes, func(q domain.Quote, _ int) string { return q.Text }), msgHistory)
	if err != nil {
		return "", fmt.Errorf("get quote: %w", err)
	}

	return q, nil
}
