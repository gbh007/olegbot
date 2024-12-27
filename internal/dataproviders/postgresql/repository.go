package postgresql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	pool *pgxpool.Pool
	db   *sqlx.DB

	logger *slog.Logger
	debug  bool
}

func New(ctx context.Context, dataSourceName string, maxConn int32, logger *slog.Logger, debug bool) (*Repository, error) {
	pgxConfig, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if maxConn > 0 {
		pgxConfig.MaxConns = maxConn
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	db := sqlx.NewDb(stdlib.OpenDBFromPool(dbpool), "pgx")

	err = migrate(ctx, logger, db.DB)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &Repository{
		pool:   dbpool,
		db:     db,
		logger: logger,
		debug:  debug,
	}, nil
}

func (r *Repository) squirrelDebugLog(ctx context.Context, query string, args []any) {
	if !r.debug {
		return
	}

	r.logger.DebugContext(
		ctx, "squirrel build request",
		slog.String("query", query),
		slog.Any("args", args),
	)
}
