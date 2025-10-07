package domain

import (
	"context"
	"errors"
)

var ErrLimitExceeded = errors.New("limit exceeded")

type Llm interface {
	GetQuote(ctx context.Context, prompt string, messages []string) (string, error)
}
