package csub

import (
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/go-chi/chi/v5/middleware"
)

type createSubscription interface {
	CreateSubscription(sub *model.Subscription) error
}

type checker interface {
	CheckBody(sub *model.Subscription) bool
	CheckSubActive(sub *model.Subscription) bool
}

func Handler(
	l *slog.Logger, cs createSubscription,
	c checker, w http.ResponseWriter,
	r *http.Request,
) {
	const fn = "handlers.csub.Handler"
	log := l.With(
		slog.String("fn", fn),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)

	sub, err := model.GetSubFromBody(r)
	if err != nil {
		log.Error("failed get user body", slog.String("err", err.Error()))
		http.Error(w, "Wrong body", http.StatusBadRequest)

		return
	}

	log.Debug("request body", slog.Any("body", sub))

	if !c.CheckBody(sub) {
		log.Error("bad body", slog.Any("body", sub))
		http.Error(w, "Wrong body", http.StatusBadRequest)

		return
	}

	if !c.CheckSubActive(sub) {
		log.Error("subscription already active")
		http.Error(w, "subscription already active", http.StatusConflict)

		return
	}

	err = cs.CreateSubscription(sub)
	if err != nil {
		log.Error("failed to create subscription", slog.String("err", err.Error()))
		http.Error(w, "something wrong", http.StatusInternalServerError)

		return
	}

	log.Info("Subscription created successfully!", slog.String("userID", sub.UserID.String()))
	w.WriteHeader(http.StatusCreated)
}
