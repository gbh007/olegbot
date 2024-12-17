package postgresql

import (
	"app/internal/dataproviders/postgresql/migration"
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // драйвер для PostgreSQL

	migrator "gitlab.com/gbh007/go-sql-migrator"
)

type Repository struct {
	data atomic.Pointer[[]string]

	moderators atomic.Pointer[map[int64]struct{}]

	emojiList   []string // FIXME: брать из БД
	emojiChance float32  // FIXME: брать из БД
	tags        []string // FIXME: брать из БД

	db *sqlx.DB
}

func New(
	emojiList []string,
	emojiChance float32,
	rawTags []string,
	botName, botTag string,
) *Repository {
	tags := make([]string, 0, len(rawTags)+2)

	if botName != "" {
		tags = append(tags, strings.ToLower(botName))
	}

	if botTag != "" {
		tags = append(tags, strings.ToLower(botTag))
	}

	for _, tag := range rawTags {
		if tag != "" {
			tags = append(tags, strings.ToLower(tag))
		}
	}

	return &Repository{
		data:        atomic.Pointer[[]string]{},
		moderators:  atomic.Pointer[map[int64]struct{}]{},
		emojiList:   emojiList,
		emojiChance: emojiChance,
		tags:        tags,
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
