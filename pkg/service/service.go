package service

import (
	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
)

type SubscriptionService interface {
	CheckBody(sub *model.Subscription) bool
	CheckSubActive(sub *model.Subscription) bool
}

type Service struct {
	database storage.Storage
}

func InitService(db storage.Storage) *Service {
	return &Service{database: db}
}

func (s *Service) CheckBody(sub *model.Subscription) bool {
	return true
}

func (s *Service) CheckSubActive(sub *model.Subscription) bool {
	return true
}
