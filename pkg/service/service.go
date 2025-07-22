package service

import (
	"context"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
)

type SubscriptionService interface {
	CheckBody(sub *model.Subscription) bool
	CheckSubscription(ctx context.Context, sub *model.Subscription) (bool, error)
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
