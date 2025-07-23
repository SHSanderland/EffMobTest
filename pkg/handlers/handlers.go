package handlers

import (
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/handlers/csub"
	"github.com/SHSanderland/EffMobTest/pkg/handlers/dsub"
	"github.com/SHSanderland/EffMobTest/pkg/handlers/lsub"
	"github.com/SHSanderland/EffMobTest/pkg/handlers/rsub"
	"github.com/SHSanderland/EffMobTest/pkg/handlers/usub"
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

func (sh *SubscriptionHandlers) ReadSubscription(w http.ResponseWriter, r *http.Request) {
	rsub.Handler(sh.log, sh.database, sh.service, w, r)
}

func (sh *SubscriptionHandlers) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	usub.Handler(sh.log, sh.database, sh.service, w, r)
}

func (sh *SubscriptionHandlers) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	dsub.Handler(sh.log, sh.database, sh.service, w, r)
}

func (sh *SubscriptionHandlers) ListSubscription(w http.ResponseWriter, r *http.Request) {
	lsub.Handler(sh.log, sh.database, sh.service, w, r)
}

func (sh *SubscriptionHandlers) CostSubscription(w http.ResponseWriter, r *http.Request) {
	// costsub.Handler(sh.log, )
}
