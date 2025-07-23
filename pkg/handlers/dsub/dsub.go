package dsub

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/service"
	"github.com/go-chi/chi/v5/middleware"
)

type deleteSubscription interface {
	DeleteSubscription(ctx context.Context, id int64) error
}

type urlParser interface {
	GetSubID(r *http.Request) (int64, error)
}

func Handler(
	l *slog.Logger, ds deleteSubscription, up urlParser,
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

	err = ds.DeleteSubscription(r.Context(), intsubID)
	if err != nil {
		log.Error("failed to delete subscription from DB", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	log.Info("Subscription delete successfully!", slog.Int64("ID", intsubID))
}
