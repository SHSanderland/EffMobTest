package handlers

import (
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/handlers/csub"
	"github.com/SHSanderland/EffMobTest/pkg/service"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
)

type SubscriptionHandlers struct {
	log      *slog.Logger
	service  service.SubscriptionService
	database storage.Storage
}

func InitHandlers(log *slog.Logger, db storage.Storage) SubscriptionHandlers {
	service := service.InitService(db)

	return SubscriptionHandlers{log: log, database: db, service: service}
}

func (sh *SubscriptionHandlers) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	csub.Handler(sh.log, sh.database, sh.service, w, r)
}
