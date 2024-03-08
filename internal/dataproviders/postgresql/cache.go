package postgresql

import (
	"context"
	"fmt"
)

func (r *Repository) reloadCache(ctx context.Context) error {
	err := r.reloadQuoteCache(ctx)
	if err != nil {
		return err
	}

	err = r.reloadModeratorCache(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) reloadModeratorCache(ctx context.Context) error {
	rawModerators, err := r.allModerators(ctx)
	if err != nil {
		return fmt.Errorf("reload moderator cache: %w", err)
	}

	moderators := make(map[int64]struct{}, len(rawModerators))
	for _, moderator := range rawModerators {
		moderators[moderator.UserID] = struct{}{}
	}

	r.moderatorsMutex.Lock()
	defer r.moderatorsMutex.Unlock()

	r.moderators = moderators

	return nil
}

func (r *Repository) reloadQuoteCache(ctx context.Context) error {
	rawQuotes, err := r.allQuotes(ctx)
	if err != nil {
		return fmt.Errorf("reload quote cache: %w", err)
	}

	data := make([]string, 0, len(rawQuotes))
	for _, quote := range rawQuotes {
		data = append(data, quote.Text)
	}

	quoteCount.Set(float64(len(data)))

	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	r.data = data

	return nil
}
