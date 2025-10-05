package ollama

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/ollama/ollama/api"
)

type Provider struct {
	client   *api.Client
	llmModel string
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	u string,
	llmModel string,
) (*Provider, error) {
	uu, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	client := api.NewClient(uu, http.DefaultClient)

	list, err := client.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list models: %w", err)
	}

	found := false

	for _, m := range list.Models {
		logger.DebugContext(ctx, "exists model", "name", m.Name, "model", m.Model)

		if m.Model == llmModel {
			found = true
		}
	}

	if !found {
		go func() {
			err = client.Pull(ctx, &api.PullRequest{
				Model: llmModel,
			}, func(_ api.ProgressResponse) error {
				return nil
			})
			if err != nil {
				logger.ErrorContext(ctx, "load model", "err", err.Error())
			}
		}()
	}

	return &Provider{
		client:   client,
		llmModel: llmModel,
	}, nil
}

func (p Provider) GetQuote(ctx context.Context, prompt string, messages []string) (string, error) {
	buff := &bytes.Buffer{}

	err := p.client.Chat(ctx, &api.ChatRequest{
		Model:    p.llmModel,
		Messages: makeMessages(prompt, messages),
	}, func(resp api.ChatResponse) error {
		buff.WriteString(resp.Message.Content)

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("make request: %w", err)
	}

	s := buff.String()

	if strings.Contains(s, "</think>") {
		s = strings.TrimSpace(strings.Split(s, "</think>")[1])
	}

	return s, nil
}

func makeMessages(prompt string, messages []string) []api.Message {
	result := make([]api.Message, 0, len(messages)+1)

	result = append(result, api.Message{
		Role:    "system",
		Content: prompt,
	})

	for _, msg := range messages {
		result = append(result, api.Message{
			Role:    "user",
			Content: msg,
		})
	}

	return result
}
