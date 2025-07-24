package costsub

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/go-chi/chi/v5/middleware"
)

type costSubscription interface {
	CostSubscription(ctx context.Context, filter *model.CostParams) (int64, error)
}

type helper interface {
	GetCostParams(r *http.Request) (*model.CostParams, error)
}

type userResponse struct {
	ServiceName string     `json:"service_name"`
	StartDate   time.Time  `json:"start_period"`
	EndDate     *time.Time `json:"end_period"`
	TotalCost   int64      `json:"total_cost"`
}

func Handler(
	l *slog.Logger, cs costSubscription, h helper,
	w http.ResponseWriter, r *http.Request,
) {
	const fn = "handlers.costsub.Handler"
	log := l.With(
		slog.String("fn", fn),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)

	filters, err := h.GetCostParams(r)
	if err != nil {
		log.Error("invalid params", slog.String("err", err.Error()))
		http.Error(w, "Invalid params", http.StatusBadRequest)

		return
	}

	total, err := cs.CostSubscription(r.Context(), filters)
	if err != nil {
		log.Error("failed to count total cost", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	ur := userResponse{
		ServiceName: filters.ServiceName,
		StartDate:   filters.StartDate,
		EndDate:     filters.EndDate,
		TotalCost:   total,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&ur); err != nil {
		log.Error("failed to send response", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	log.Info(
		"Total cost counted successfully!",
		slog.String("userID", filters.UserID.String()),
		slog.String("serviceName", filters.ServiceName),
	)
}
