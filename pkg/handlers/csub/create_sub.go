// Пакет csub для хендлера CreateSubscription.
package csub

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/go-chi/chi/v5/middleware"
)

// createSubscription Интерефейс с методами к базе данных,
// который использует хендлер.
type createSubscription interface {
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
}

// checker Интерефейс с методами к Service,
// который использует хендлер.
type checker interface {
	CheckBody(sub *model.Subscription) bool
	CheckSubscription(ctx context.Context, sub *model.Subscription) (bool, error)
}

// @Summary		Создать новую подписку
// @Description	Создает новую подписку после проверки валидности данных и отсутствия активной подписки
// @Tags			subscriptions
// @Accept			json
// @Produce		plain
// @Param			input	body	model.Subscription	true	"Данные для создания подписки"
// @Success		201		"Подписка успешно создана"
// @Failure		400		{string}	string	"Невалидные входные данные"
// @Failure		409		{string}	string	"Подписка уже активна"
// @Failure		500		{string}	string	"Внутренняя ошибка сервера"
// @Router			/subscriptions [post]
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

	isActive, err := c.CheckSubscription(r.Context(), sub)
	if err != nil {
		log.Error("failed to check subscription", slog.String("err", err.Error()))
		http.Error(w, "something wrong", http.StatusInternalServerError)

		return
	}

	if isActive {
		log.Error("subscription already active")
		http.Error(w, "subscription already active", http.StatusConflict)

		return
	}

	err = cs.CreateSubscription(r.Context(), sub)
	if err != nil {
		log.Error("failed to create subscription", slog.String("err", err.Error()))
		http.Error(w, "something wrong", http.StatusInternalServerError)

		return
	}

	log.Info("Subscription created successfully!", slog.String("userID", sub.UserID.String()))
	w.WriteHeader(http.StatusCreated)
}
