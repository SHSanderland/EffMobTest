package dsub

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type deleteSubscription interface {
	DeleteSubscription(ctx context.Context, id int64) error
}

func Handler(
	l *slog.Logger, ds deleteSubscription,
	w http.ResponseWriter, r *http.Request,
) {
	const fn = "handlers.rsub.Handler"
	log := l.With(
		slog.String("fn", fn),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)

	subID := chi.URLParam(r, "id")

	intsubID, err := strconv.ParseInt(subID, 10, 64)
	if err != nil || intsubID < 0 {
		log.Error(
			"invalid subscription ID",
			slog.Int64("ID", intsubID),
			slog.String("err", err.Error()),
		)
		http.Error(w, "Invalid ID", http.StatusBadRequest)

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
