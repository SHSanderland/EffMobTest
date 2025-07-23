package lsub

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type listSubscription interface {
	GetListSubscription(ctx context.Context, userID uuid.UUID, serviceName string) ([]*model.Subscription, error)
}

type urlParser interface {
	GetUserIDAndServiceName(r *http.Request) (uuid.UUID, string, error)
}

type userResponse struct {
	Subscriptions []*model.Subscription `json:"subscriptions"`
	Total         int                   `json:"total"`
}

func Handler(
	l *slog.Logger, ls listSubscription, up urlParser,
	w http.ResponseWriter, r *http.Request,
) {
	const fn = "handlers.lsub.Handler"
	log := l.With(
		slog.String("fn", fn),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)

	userID, serviceName, err := up.GetUserIDAndServiceName(r)
	if err != nil {
		log.Error("failed to parse url", slog.String("err", err.Error()))
		http.Error(w, "Invalid URL params", http.StatusBadRequest)

		return
	}

	subs, err := ls.GetListSubscription(r.Context(), userID, serviceName)
	if err != nil {
		log.Error("failed to get list subscriptions", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	userResp := userResponse{subs, len(subs)}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(userResp); err != nil {
		log.Error("failed to encode json", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	log.Info(
		"List of subscriptions sended successfully!",
		slog.String("userID", userID.String()),
		slog.String("serviceName", serviceName),
	)
}
