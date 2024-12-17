package postgresql

import (
	"app/internal/dataproviders/postgresql/migration"
	"app/internal/domain"
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // драйвер для PostgreSQL

	migrator "gitlab.com/gbh007/go-sql-migrator"
)

type Repository struct {
	data atomic.Pointer[[]string]

	moderators atomic.Pointer[map[int64]struct{}]

	botInfo domain.Bot // FIXME: получать данные из БД

	db *sqlx.DB
}

func New(botInfo domain.Bot) *Repository {
	return &Repository{
		data:       atomic.Pointer[[]string]{},
		moderators: atomic.Pointer[map[int64]struct{}]{},
		botInfo:    botInfo,
	}
}

func (r *Repository) connect(_ context.Context, source string) error {
	db, err := sqlx.Open("postgres", source)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	db.SetMaxOpenConns(10)

	r.db = db

	return nil
}

func (r *Repository) migrate(ctx context.Context) error {
	err := migrator.New().
		WithFS(migration.Migrations).
		WithProvider(migrator.PostgreSQLProvider).
		MigrateAll(ctx, r.db, true)

	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
