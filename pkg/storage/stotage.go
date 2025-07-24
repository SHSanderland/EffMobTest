package storage

import (
	"context"
	"errors"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/google/uuid"
)

var (
	ErrBeginTrans  = errors.New("failed to begin transaction")
	ErrCommitTrans = errors.New("failed to commit transaction")
	ErrExecSchema  = errors.New("failed to exec schema")
	ErrEmptySub    = errors.New("nothing to update")
)

type Storage interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	ReadSubscription(ctx context.Context, subID int64) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, subID int64, sub *model.Subscription) error
	DeleteSubscription(ctx context.Context, subID int64) error
	GetListSubscription(ctx context.Context, userID uuid.UUID, serviceName string) ([]*model.Subscription, error)
	CostSubscription(ctx context.Context, filter *model.CostParams) (int64, error)
	CloseConnection()
	CheckStorage
}

type CheckStorage interface {
	CheckSubscription(ctx context.Context, sub *model.Subscription) (bool, error)
	CheckSubscriptionID(ctx context.Context, subID int64) (bool, error)
	CheckSubscriptionForUpdate(ctx context.Context, subID int64, sub *model.Subscription) (bool, error)
}

const (
	CreateSubscriptionSchema = `
		INSERT INTO subscriptions (
			service_name, price, user_id, start_date, end_date
		)
		VALUES (
			$1, $2, $3, TO_DATE($4, 'MM-YYYY'), TO_DATE($5, 'MM-YYYY')
		);
	`
	SubscriptionActiveSchema = `
		SELECT EXISTS (
			SELECT 1
			FROM subscriptions
			WHERE user_id = $1
				AND service_name = $2
				AND (
					end_date IS NULL 
					OR end_date > CURRENT_DATE
				)
				AND CURRENT_DATE >= start_date
		);
	`
	ReadSubscriptionSchema = `
		SELECT service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1;
	`
	ReadSubscriptionStartDateSchema = `
		SELECT start_date
		FROM subscriptions
		WHERE id = $1;
	`
	SubscriptionExistsSchema = `
		SELECT EXISTS (
			SELECT 1
			FROM subscriptions
			WHERE id = $1
		);	
	`
	DeleteSubscriptionSchema = `
		DELETE FROM subscriptions
		WHERE id = $1;
	`
	ListSubscriptionSchema = `
		SELECT service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1 AND service_name = $2;
	`
	CountSubscriptionsSchema = `
		SELECT sum(price)
		FROM subscriptions
		WHERE user_id = $1 
			AND service_name = $2
			AND start_date >= $3
			AND (end_date <= $4 OR end_date IS NULL);
	`
)
