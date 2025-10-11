package mistral

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"app/internal/domain"

	"github.com/gage-technologies/mistral-go"
)

type Provider struct {
	client   *mistral.MistralClient
	llmModel string
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	token string,
	llmModel string,
) (*Provider, error) {
	client := mistral.NewMistralClientDefault(token)

	return &Provider{
		client:   client,
		llmModel: llmModel,
	}, nil
}

func (p Provider) GetQuote(ctx context.Context, prompt string, messages []string) (string, error) {
	params := mistral.DefaultChatRequestParams

	params.Temperature = 1.5

	resp, err := p.client.Chat(
		p.llmModel,
		makeMessages(prompt, messages),
		&params,
	)
	if err != nil {
		s := err.Error()

		if strings.Contains(s, "429") || strings.Contains(s, "Rate limit exceeded") {
			return "", domain.ErrLimitExceeded
		}

		return "", fmt.Errorf("make request: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("empty choices")
	}

	return resp.Choices[0].Message.Content, nil
}

func makeMessages(prompt string, messages []string) []mistral.ChatMessage {
	result := make([]mistral.ChatMessage, 0, len(messages)+1)

	result = append(result, mistral.ChatMessage{
		Role:    mistral.RoleSystem,
		Content: prompt,
	})

	for _, msg := range messages {
		result = append(result, mistral.ChatMessage{
			Role:    mistral.RoleUser,
			Content: msg,
		})
	}

	return result
}
