package psql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"

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

func (s *Storage) ReadSubscription(ctx context.Context, subID int64) (*model.Subscription, error) {
	const fn = "psql.ReadSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.Int64("subID", subID),
	)

	var (
		sub       model.Subscription
		startTime time.Time
		endTime   *time.Time
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", slog.String("err", err.Error()))

		return nil, fmt.Errorf("%w: %w", storage.ErrBeginTrans, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("failed to rollback transaction", slog.String("err", err.Error()))
		}
	}()

	err = tx.QueryRow(ctx, storage.ReadSubscriptionSchema, subID).Scan(
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&startTime,
		&endTime,
	)
	if err != nil {
		log.Error("failed to exec schema", slog.String("err", err.Error()))

		return nil, fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return nil, fmt.Errorf("%w: %w", storage.ErrCommitTrans, err)
	}

	sub.StartDate = startTime.Format("01-2006")
	if endTime != nil {
		sub.EndDate = endTime.Format("01-2006")
	}

	log.Info("Subscription is readed!")

	return &sub, nil
}

func (s *Storage) UpdateSubscription(ctx context.Context, subID int64, sub *model.Subscription) error {
	const fn = "psql.UpdateSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.Int64("subID", subID),
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

	updates, args := prepareUpdate(sub)

	if len(updates) == 0 {
		return storage.ErrEmptySub
	}

	query := fmt.Sprintf(
		"UPDATE subscriptions SET %s WHERE id = $%d",
		strings.Join(updates, ", "),
		len(updates)+1,
	)

	args = append(args, subID)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to exec schema", slog.String("err", err.Error()))

		return fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return fmt.Errorf("%w: %w", storage.ErrCommitTrans, err)
	}

	log.Info("Subscription is update!")

	return nil
}

func (s *Storage) DeleteSubscription(ctx context.Context, subID int64) error {
	const fn = "psql.DeleteSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.Int64("subID", subID),
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

	_, err = tx.Exec(ctx, storage.DeleteSubscriptionSchema, subID)
	if err != nil {
		log.Error("failed to exec schema", slog.String("err", err.Error()))

		return fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return fmt.Errorf("%w: %w", storage.ErrCommitTrans, err)
	}

	return nil
}

func (s *Storage) GetListSubscription(
	ctx context.Context, userID uuid.UUID, serviceName string,
) ([]*model.Subscription, error) {
	const fn = "psql.GetListSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.String("userID", userID.String()),
		slog.String("serviceName", serviceName),
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", slog.String("err", err.Error()))

		return nil, fmt.Errorf("%w: %w", storage.ErrBeginTrans, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("failed to rollback transaction", slog.String("err", err.Error()))
		}
	}()

	rows, err := tx.Query(ctx, storage.ListSubscriptionSchema, userID, serviceName)
	if err != nil {
		log.Error("failed to exec schema", slog.String("err", err.Error()))

		return nil, fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	subs, err := s.getSubList(rows)
	if err != nil {
		log.Error("failed to scan rows", slog.String("err", err.Error()))

		return nil, fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return nil, fmt.Errorf("%w: %w", storage.ErrCommitTrans, err)
	}

	return subs, nil
}

func (s *Storage) CostSubscription(ctx context.Context, filter *model.CostParams) (int64, error) {
	const fn = "psql.CostSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.String("userID", filter.UserID.String()),
		slog.String("serviceName", filter.ServiceName),
	)

	var total int64

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", slog.String("err", err.Error()))

		return total, fmt.Errorf("%w: %w", storage.ErrBeginTrans, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("failed to rollback transaction", slog.String("err", err.Error()))
		}
	}()

	err = tx.QueryRow(
		ctx, storage.CountSubscriptionsSchema,
		filter.UserID,
		filter.ServiceName,
		filter.StartDate,
		filter.EndDate,
	).Scan(&total)
	if err != nil {
		log.Error("failed to scan rows", slog.String("err", err.Error()))

		return total, fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return total, fmt.Errorf("%w: %w", storage.ErrCommitTrans, err)
	}

	return total, nil
}

func (s *Storage) CloseConnection() {
	s.db.Close()
	s.log.Info("Connection to DB is closed!")
}

func prepareUpdate(sub *model.Subscription) ([]string, []any) {
	var (
		updates []string
		args    []any
	)

	if sub.ServiceName != "" {
		updates = append(updates, fmt.Sprintf("service_name = $%d", len(updates)+1))
		args = append(args, sub.ServiceName)
	}

	if sub.Price > 0 {
		updates = append(updates, fmt.Sprintf("price = $%d", len(updates)+1))
		args = append(args, sub.Price)
	}

	if sub.EndDate != "" {
		updates = append(updates, fmt.Sprintf("end_date = TO_DATE($%d, 'MM-YYYY')", len(updates)+1))
		args = append(args, sub.EndDate)
	}

	return updates, args
}

func (s *Storage) getSubscriptionStartDate(ctx context.Context, subID int64) *time.Time {
	const fn = "psql.UpdateSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.Int64("subID", subID),
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", slog.String("err", err.Error()))

		return nil
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("failed to rollback transaction", slog.String("err", err.Error()))
		}
	}()

	var startDateFromDB *time.Time

	err = tx.QueryRow(ctx, storage.ReadSubscriptionStartDateSchema, subID).Scan(
		&startDateFromDB,
	)
	if err != nil {
		log.Error("failed to exec schema", slog.String("err", err.Error()))

		return nil
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return nil
	}

	return startDateFromDB
}

func (s *Storage) getSubList(rows pgx.Rows) ([]*model.Subscription, error) {
	var subs []*model.Subscription

	for rows.Next() {
		var (
			sub       model.Subscription
			startDate time.Time
			endDate   *time.Time
		)

		err := rows.Scan(
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&startDate,
			&endDate,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
		}

		sub.StartDate = startDate.Format("01-2006")
		if endDate != nil {
			sub.EndDate = endDate.Format("01-2006")
		}

		subs = append(subs, &sub)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrExecSchema, rows.Err())
	}

	return subs, nil
}
