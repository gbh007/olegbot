package openai

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/samber/lo/mutable"
)

type Provider struct {
	client   openai.Client
	llmModel string
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	u, token string,
	llmModel string,
) (*Provider, error) {
	client := openai.NewClient(
		option.WithAPIKey(token),
		option.WithBaseURL(u),
	)

	return &Provider{
		client:   client,
		llmModel: llmModel,
	}, nil
}

func (p Provider) GetQuote(ctx context.Context, names, quotes, messages []string) (string, error) {
	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       p.llmModel,
		Messages:    makeMessages(names, quotes, messages),
		Temperature: openai.Float(1.5),
	})
	if err != nil {
		return "", fmt.Errorf("make request: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("empty choices")
	}

	return resp.Choices[0].Message.Content, nil
}

func makeMessages(names, quotes, messages []string) []openai.ChatCompletionMessageParamUnion {
	quotes = slices.Clone(quotes)
	mutable.Shuffle(quotes)
	if len(quotes) > 50 {
		quotes = quotes[:50]
	}

	result := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)

	result = append(result, openai.SystemMessage("Тебя зовут именами "+
		strings.Join(names, ", ")+
		". Ты пишешь смешные фразы, на основании следующих фраз придумай новый смешной ответ для пользователя из одной фразы:\n"+
		strings.Join(quotes, "\n"),
	))

	for _, msg := range messages {
		result = append(result, openai.UserMessage(msg))
	}

	return result
}
