package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidSubID       = errors.New("invalid subscription ID")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidServiceName = errors.New("invalid service name")
	ErrInvalidDate        = errors.New("invalid date")
)

type SubscriptionService interface {
	CheckBody(sub *model.Subscription) bool
	CheckSubscription(ctx context.Context, sub *model.Subscription) (bool, error)
	CheckSubscriptionID(ctx context.Context, subID int64) (bool, error)
	CheckSubscriptionForUpdate(ctx context.Context, subID int64, sub *model.Subscription) (bool, error)
	GetSubID(r *http.Request) (int64, error)
	GetUserIDAndServiceName(r *http.Request) (uuid.UUID, string, error)
	GetCostParams(r *http.Request) (*model.CostParams, error)
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

func (s *Service) GetSubID(r *http.Request) (int64, error) {
	subID := chi.URLParam(r, "id")

	intsubID, err := strconv.ParseInt(subID, 10, 64)
	if err != nil || intsubID < 0 {
		return intsubID, fmt.Errorf("%w: %w", ErrInvalidSubID, err)
	}

	return intsubID, nil
}

func (s *Service) GetUserIDAndServiceName(r *http.Request) (uuid.UUID, string, error) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		return uuid.Nil, "", ErrInvalidUserID
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return userUUID, "", fmt.Errorf("%w: %w", ErrInvalidUserID, err)
	}

	serviceName := r.URL.Query().Get("service_name")
	if serviceName == "" {
		return userUUID, serviceName, ErrInvalidServiceName
	}

	return userUUID, serviceName, nil
}

func (s *Service) GetCostParams(r *http.Request) (*model.CostParams, error) {
	serviceName := r.URL.Query().Get("service_name")
	userID := r.URL.Query().Get("user_id")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	layout := "01-2006"

	if serviceName == "" {
		return nil, ErrInvalidServiceName
	}

	if userID == "" {
		return nil, ErrInvalidUserID
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidUserID, err)
	}

	if startDateStr == "" {
		return nil, ErrInvalidDate
	}

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidDate, err)
	}

	if endDateStr == "" {
		return nil, ErrInvalidDate
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidDate, err)
	}

	cost := model.CostParams{
		ServiceName: serviceName,
		UserID:      userUUID,
		StartDate:   startDate,
		EndDate:     &endDate,
	}

	return &cost, nil
}
