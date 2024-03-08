package postgresql

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
