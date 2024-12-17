package postgresql

import "context"

func (r *Repository) Tags(ctx context.Context) ([]string, error) {
	return r.tags, nil
}
