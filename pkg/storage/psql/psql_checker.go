package psql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CheckSubscriptionID(ctx context.Context, subID int64) (bool, error) {
	const fn = "psql.CheckSubscriptionID"
	log := s.log.With(
		slog.String("fn", fn),
		slog.Int64("subID", subID),
	)

	var hasID bool

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", slog.String("err", err.Error()))

		return hasID, fmt.Errorf("%w: %w", storage.ErrBeginTrans, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("failed to rollback transaction", slog.String("err", err.Error()))
		}
	}()

	err = tx.QueryRow(ctx, storage.SubscriptionExistsSchema, subID).Scan(&hasID)
	if err != nil {
		log.Error("failed to exec schema", slog.String("err", err.Error()))

		return hasID, fmt.Errorf("%w: %w", storage.ErrExecSchema, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("err", err.Error()))

		return hasID, fmt.Errorf("%w: %w", storage.ErrCommitTrans, err)
	}

	log.Info("Subscription is checked!")

	return hasID, nil
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

func (s *Storage) CheckSubscriptionForUpdate(ctx context.Context, subID int64, sub *model.Subscription) (bool, error) {
	const fn = "psql.CheckSubscription"
	log := s.log.With(
		slog.String("fn", fn),
		slog.String("userID", sub.UserID.String()),
	)

	if sub.Price < 0 {
		log.Error("update price is negetive", slog.Int("price", sub.Price))

		return false, nil
	}

	if sub.EndDate == "" {
		return true, nil
	}

	layout := "01-2006"

	endDateFromSub, err := time.Parse(layout, sub.EndDate)
	if err != nil {
		log.Error("failed to parse date", slog.String("err", err.Error()))

		return false, fmt.Errorf("failed to parse date: %w", err)
	}

	startDateFromDB := s.getSubscriptionStartDate(ctx, subID)

	if !endDateFromSub.After(*startDateFromDB) {
		return false, nil
	}

	return true, nil
}
