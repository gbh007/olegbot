package multillm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"app/internal/dataproviders/deepseek"
	"app/internal/dataproviders/ollama"
	"app/internal/dataproviders/openai"
	"app/internal/domain"
)

type Config struct {
	Addr  string `toml:"addr"`
	Model string `toml:"model"`
	Token string `toml:"token"`
	Type  string `toml:"type"`
}

type llmProvider struct {
	provider    domain.Llm
	limitBefore time.Time
	errCount    int
}

type Provider struct {
	providers      []llmProvider
	activeProvider int
	mutex          sync.Mutex
	logger         *slog.Logger
}

func New(ctx context.Context, logger *slog.Logger, cfgs ...Config) (*Provider, error) {
	providers := make([]llmProvider, 0, len(cfgs))

	for _, cfg := range cfgs {
		var (
			provider domain.Llm
			err      error
		)
		switch {
		case cfg.Type == "openai" && cfg.Token != "" && cfg.Addr != "" && cfg.Model != "":
			provider, err = openai.New(ctx, logger, cfg.Addr, cfg.Token, cfg.Model)
			if err != nil {
				return nil, fmt.Errorf("app: init: openai: %w", err)
			}
		case cfg.Type == "deepseek" && cfg.Token != "":
			provider, err = deepseek.New(ctx, logger, cfg.Token)
			if err != nil {
				return nil, fmt.Errorf("app: init: deepseek: %w", err)
			}
		case cfg.Type == "ollama" && cfg.Addr != "" && cfg.Model != "":
			provider, err = ollama.New(ctx, logger, cfg.Addr, cfg.Model)
			if err != nil {
				return nil, fmt.Errorf("app: init: ollama: %w", err)
			}
		}

		providers = append(providers, llmProvider{
			provider: provider,
		})
	}

	return &Provider{
		providers:      providers,
		logger:         logger,
		activeProvider: 0,
		mutex:          sync.Mutex{},
	}, nil
}

func (p *Provider) GetQuote(ctx context.Context, prompt string, messages []string) (string, error) {
	p.mutex.Lock()

	now := time.Now()

	for i, pr := range p.providers {
		if p.activeProvider <= i {
			break
		}

		if pr.limitBefore.Before(now) {
			if i < p.activeProvider {
				p.logger.WarnContext(ctx, "swap provider due to limits end", "old", p.activeProvider, "new", i)

				p.activeProvider = i
				p.providers[i].errCount = 0

				break
			}
		}
	}

	cur := p.activeProvider
	provider := p.providers[cur]

	p.mutex.Unlock()

	q, err := provider.provider.GetQuote(ctx, prompt, messages)
	if err == nil {
		return q, nil
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.activeProvider != cur {
		return "", err
	}

	if errors.Is(err, domain.ErrLimitExceeded) {
		p.providers[cur].limitBefore = time.Now().Add(time.Hour * 4)
		p.activeProvider = (cur + 1) % len(p.providers)

		p.logger.WarnContext(ctx, "swap provider due to limits", "old", cur, "new", p.activeProvider)

		return "", err
	}

	p.providers[cur].errCount++

	if p.providers[cur].errCount > 3 {
		p.providers[cur].limitBefore = time.Now().Add(time.Minute * 10)
		p.activeProvider = (cur + 1) % len(p.providers)

		p.logger.WarnContext(ctx, "swap provider due to errors", "old", cur, "new", p.activeProvider)
	}

	return "", err
}
