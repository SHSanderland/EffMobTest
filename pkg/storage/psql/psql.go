package psql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/SHSanderland/EffMobTest/pkg/config"
	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
	"github.com/jackc/pgx/v5"
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

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error("failed to migrate DB", slog.String("err", err.Error()))

		return nil, fmt.Errorf("failed to migrate DB: %w", err)
	}

	return &Storage{log: log, db: pool}, nil
}

func (s *Storage) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
	const fn = "psql.CreateSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.String("userID", sub.UserID.String()),
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", slog.String("err", err.Error()))

		return fmt.Errorf("%w: %w", storage.ErrBeginTrans, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("failed to rollback transaction", slog.String("err", err.Error()))
		}
	}()

	if sub.EndDate == "" {
		_, err = tx.Exec(
			ctx,
			storage.CreateSubscriptionSchema,
			sub.ServiceName,
			sub.Price,
			sub.UserID,
			sub.StartDate,
			nil,
		)
	} else {
		_, err = tx.Exec(
			ctx,
			storage.CreateSubscriptionSchema,
			sub.ServiceName,
			sub.Price,
			sub.UserID,
			sub.StartDate,
			sub.EndDate,
		)
	}

	if err != nil {
		log.Error("failed to exec schema", slog.String("err", err.Error()))

		return fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return fmt.Errorf("%w: %w", storage.ErrCommitTrans, err)
	}

	log.Info("Subscription is created!")

	return nil
}

func (s *Storage) CheckSubscription(ctx context.Context, sub *model.Subscription) (bool, error) {
	const fn = "psql.CheckSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.String("userID", sub.UserID.String()),
	)

	var isActive bool

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", slog.String("err", err.Error()))

		return isActive, fmt.Errorf("%w: %w", storage.ErrBeginTrans, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("failed to rollback transaction", slog.String("err", err.Error()))
		}
	}()

	err = tx.QueryRow(
		ctx,
		storage.SubscriptionActiveSchema,
		sub.UserID,
		sub.ServiceName,
	).Scan(&isActive)
	if err != nil {
		log.Error("failed to exec schema", slog.String("err", err.Error()))

		return isActive, fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return isActive, fmt.Errorf("%w: %w", storage.ErrCommitTrans, err)
	}

	log.Info("Subscription is checked!")

	return isActive, nil
}

func (s *Storage) CloseConnection() {
	s.db.Close()
	s.log.Info("Connection to DB is closed!")
}
