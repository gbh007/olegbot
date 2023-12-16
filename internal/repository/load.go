package repository

import (
	"context"
	"fmt"
)

func (r *Repository) Load(ctx context.Context, source string) error {
	err := r.connect(ctx, source)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	err = r.migrate(ctx)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	err = r.reloadCache(ctx)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	return nil
}

func (r *Repository) reloadCache(ctx context.Context) error {
	rawQuotes, err := r.allQuotes(ctx)
	if err != nil {
		return fmt.Errorf("reload cache: %w", err)
	}

	rawModerators, err := r.allModerators(ctx)
	if err != nil {
		return fmt.Errorf("reload cache: %w", err)
	}

	data := make([]string, len(rawQuotes))
	for _, quote := range rawQuotes {
		data = append(data, quote.Text)
	}

	moderators := make(map[int64]struct{}, len(rawModerators))
	for _, moderator := range rawModerators {
		moderators[moderator.UserID] = struct{}{}
	}

	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	r.moderatorsMutex.Lock()
	defer r.moderatorsMutex.Unlock()

	r.data = data
	r.moderators = moderators

	return nil
}
