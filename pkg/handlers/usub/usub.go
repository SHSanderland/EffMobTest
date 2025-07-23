package usub

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/service"
	"github.com/go-chi/chi/v5/middleware"
)

type updateSubscription interface {
	UpdateSubscription(ctx context.Context, id int64, sub *model.Subscription) error
}

type helper interface {
	CheckSubscriptionID(ctx context.Context, subID int64) (bool, error)
	CheckSubscriptionForUpdate(ctx context.Context, subID int64, sub *model.Subscription) (bool, error)
	GetSubID(r *http.Request) (int64, error)
}

func Handler(
	l *slog.Logger, us updateSubscription, h helper,
	w http.ResponseWriter, r *http.Request,
) {
	const fn = "handlers.rsub.Handler"
	log := l.With(
		slog.String("fn", fn),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)

	intsubID, err := h.GetSubID(r)
	if err != nil {
		log.Error(
			service.ErrInvalidSubID.Error(),
			slog.Int64("ID", intsubID),
			slog.String("err", err.Error()),
		)
		http.Error(w, service.ErrInvalidSubID.Error(), http.StatusBadRequest)

		return
	}

	hasID, err := h.CheckSubscriptionID(r.Context(), intsubID)
	if err != nil {
		log.Error("failed check subscription ID", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	if !hasID {
		log.Error("ID not exists", slog.Int64("ID", intsubID))
		http.Error(w, "Subscription ID not ex", http.StatusNotFound)

		return
	}

	sub, err := model.GetSubFromBody(r)
	if err != nil {
		log.Error("failed to get body", slog.String("err", err.Error()))
		http.Error(w, "Bad body", http.StatusBadRequest)

		return
	}

	isValidUpdate, err := h.CheckSubscriptionForUpdate(r.Context(), intsubID, sub)
	if err != nil {
		log.Error("failed check subscription for update", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	if !isValidUpdate {
		log.Error("update not valid", slog.Any("sub", sub))
		http.Error(w, "Bad body", http.StatusBadRequest)

		return
	}

	err = us.UpdateSubscription(r.Context(), intsubID, sub)
	if err != nil {
		log.Error("failed update subscription", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	log.Info("Subscription update successfully!", slog.Int64("ID", intsubID))
}
