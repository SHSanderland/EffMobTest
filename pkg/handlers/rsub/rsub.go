package rsub

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/service"
	"github.com/go-chi/chi/v5/middleware"
)

type readSubscription interface {
	ReadSubscription(ctx context.Context, id int64) (*model.Subscription, error)
}

type urlParser interface {
	GetSubID(r *http.Request) (int64, error)
}

func Handler(
	l *slog.Logger, rs readSubscription, up urlParser,
	w http.ResponseWriter, r *http.Request,
) {
	const fn = "handlers.rsub.Handler"
	log := l.With(
		slog.String("fn", fn),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)

	intsubID, err := up.GetSubID(r)
	if err != nil {
		log.Error(
			service.ErrInvalidSubID.Error(),
			slog.Int64("ID", intsubID),
			slog.String("err", err.Error()),
		)
		http.Error(w, service.ErrInvalidSubID.Error(), http.StatusBadRequest)

		return
	}

	sub, err := rs.ReadSubscription(r.Context(), intsubID)
	if err != nil {
		log.Error("failed to read subscription from DB", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(sub); err != nil {
		log.Error("failed to send JSON", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	log.Info("Subscription send successfully!", slog.Int64("ID", intsubID))
}
