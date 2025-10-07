package openai

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"app/internal/domain"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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

func (p Provider) GetQuote(ctx context.Context, prompt string, messages []string) (string, error) {
	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       p.llmModel,
		Messages:    makeMessages(prompt, messages),
		Temperature: openai.Float(1.5),
	})
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

func makeMessages(prompt string, messages []string) []openai.ChatCompletionMessageParamUnion {
	result := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)

	result = append(result, openai.SystemMessage(prompt))

	for _, msg := range messages {
		result = append(result, openai.UserMessage(msg))
	}

	return result
}
