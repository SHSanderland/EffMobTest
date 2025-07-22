package psql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/SHSanderland/EffMobTest/pkg/config"
	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func InitDB(ctx context.Context, log *slog.Logger, cfg *config.Config) (*Storage, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		log.Error("failed to add new connection to DB", slog.String("err", err.Error()))

		return nil, fmt.Errorf("failed to add new connection: %w", err)
	}

	err = pool.Ping(ctx)
	for i := 0; err != nil && i < 3; i++ {
		log.Error(
			"failed to connect to DB",
			slog.Int("Try", i+1),
			slog.String("err", err.Error()),
		)
		time.Sleep(time.Second)

		err = pool.Ping(ctx)
	}

	if err != nil {
		log.Error("failed to init DB", slog.String("err", err.Error()))

		return nil, fmt.Errorf("failed to init DB: %w", err)
	}

	m, err := migrate.New(cfg.SourceURL, cfg.DSN)
	if err != nil {
		log.Error("failed to run migrator", slog.String("err", err.Error()))

		return nil, fmt.Errorf("failed to run migrator: %w", err)
	}

	if err := m.Up(); err != nil {
		log.Error("failed to migrate DB", slog.String("err", err.Error()))

		return nil, fmt.Errorf("failed to migrate DB: %w", err)
	}

	return &Storage{log: log, db: pool}, nil
}

func (s *Storage) CreateSubscription(sub *model.Subscription) error {

	return nil
}

func (s *Storage) CloseConnection() {
	s.db.Close()
}
