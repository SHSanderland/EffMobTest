package service

import (
	"context"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
)

type SubscriptionService interface {
	CheckBody(sub *model.Subscription) bool
	CheckSubscription(ctx context.Context, sub *model.Subscription) (bool, error)
	CheckSubscriptionID(ctx context.Context, subID int64) (bool, error)
	CheckSubscriptionForUpdate(ctx context.Context, subID int64, sub *model.Subscription) (bool, error)
}

type Service struct {
	database storage.CheckStorage
}

func InitService(db storage.Storage) *Service {
	return &Service{database: db}
}

func (s *Service) CheckBody(sub *model.Subscription) bool {
	return model.IsValidSubscription(sub)
}

func (s *Service) CheckSubscription(ctx context.Context, sub *model.Subscription) (bool, error) {
	return s.database.CheckSubscription(ctx, sub)
}

func (s *Service) CheckSubscriptionID(ctx context.Context, subID int64) (bool, error) {
	return s.database.CheckSubscriptionID(ctx, subID)
}

func (s *Service) CheckSubscriptionForUpdate(ctx context.Context, subID int64, sub *model.Subscription) (bool, error) {
	return s.database.CheckSubscriptionForUpdate(ctx, subID, sub)
}
