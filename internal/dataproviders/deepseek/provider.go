package deepseek

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
	"github.com/samber/lo/mutable"
)

type Provider struct {
	client deepseek.Client
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	token string,
) (*Provider, error) {
	client, err := deepseek.NewClient(token)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return &Provider{
		client: client,
	}, nil
}

func (p Provider) GetQuote(ctx context.Context, names, quotes, messages []string) (string, error) {
	var temperature float32 = 1.5

	resp, err := p.client.CallChatCompletionsReasoner(ctx, &request.ChatCompletionsRequest{
		Messages:    makeMessages(names, quotes, messages),
		Model:       deepseek.DEEPSEEK_REASONER_MODEL,
		Temperature: &temperature,
	})
	if err != nil {
		return "", fmt.Errorf("make request: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("empty choices")
	}

	return resp.Choices[0].Message.Content, nil
}

func makeMessages(names, quotes, messages []string) []*request.Message {
	quotes = slices.Clone(quotes)
	mutable.Shuffle(quotes)
	if len(quotes) > 50 {
		quotes = quotes[:50]
	}

	result := make([]*request.Message, 0, len(messages)+1)

	result = append(result, &request.Message{
		Role: "system",
		Content: "Тебя зовут именами " +
			strings.Join(names, ", ") +
			". Ты пишешь смешные фразы, на основании следующих фраз придумай новый смешной ответ для пользователя из одной фразы:\n" +
			strings.Join(quotes, "\n"),
	})

	for _, msg := range messages {
		result = append(result, &request.Message{
			Role:    "user",
			Content: msg,
		})
	}

	return result
}
