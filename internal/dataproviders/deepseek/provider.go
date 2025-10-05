package deepseek

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
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

func (p Provider) GetQuote(ctx context.Context, prompt string, messages []string) (string, error) {
	var temperature float32 = 1.5

	resp, err := p.client.CallChatCompletionsReasoner(ctx, &request.ChatCompletionsRequest{
		Messages:    makeMessages(prompt, messages),
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

func makeMessages(prompt string, messages []string) []*request.Message {
	result := make([]*request.Message, 0, len(messages)+1)

	result = append(result, &request.Message{
		Role:    "system",
		Content: prompt,
	})

	for _, msg := range messages {
		result = append(result, &request.Message{
			Role:    "user",
			Content: msg,
		})
	}

	return result
}
