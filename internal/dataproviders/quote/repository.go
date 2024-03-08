package quote

import (
	"app/internal/dataproviders/quote/migration"
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // драйвер для PostgreSQL

	migrator "gitlab.com/gbh007/go-sql-migrator"
)

type Repository struct {
	data      []string // FIXME: переписать на атомики
	dataMutex *sync.RWMutex

	moderators      map[int64]struct{} // FIXME: переписать на атомики
	moderatorsMutex *sync.RWMutex

	db *sqlx.DB
}

func New() *Repository {
	return &Repository{
		data:            []string{},
		dataMutex:       &sync.RWMutex{},
		moderators:      make(map[int64]struct{}),
		moderatorsMutex: &sync.RWMutex{},
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
