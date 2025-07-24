// Пакет usub для хендлера UpdateSubscription.
package usub

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/SHSanderland/EffMobTest/pkg/service"
	"github.com/go-chi/chi/v5/middleware"
)

// updateSubscription Интерефейс с методами к базе данных,
// который использует хендлер.
type updateSubscription interface {
	UpdateSubscription(ctx context.Context, id int64, sub *model.Subscription) error
}

// helper Интерефейс с методами к Service,
// который использует хендлер.
type helper interface {
	CheckSubscriptionID(ctx context.Context, subID int64) (bool, error)
	CheckSubscriptionForUpdate(ctx context.Context, subID int64, sub *model.Subscription) (bool, error)
	GetSubID(r *http.Request) (int64, error)
}

// @Summary Обновить подписку
// @Description Обновляет информацию о существующей подписке по её ID
// @Tags subscriptions
// @Accept json
// @Produce plain
// @Param id path int true "ID обновляемой подписки" Example(123)
// @Param input body model.Subscription true "Новые данные подписки"
// @Success 200 "Подписка успешно обновлена"
// @Failure 400 {string} string "Невалидные входные данные (ID или тело запроса)"
// @Failure 404 {string} string "Подписка с указанным ID не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
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
