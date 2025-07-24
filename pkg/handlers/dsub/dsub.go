// Пакет dsub для хендлера DeleteSubscription.
package dsub

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/service"
	"github.com/go-chi/chi/v5/middleware"
)

// deleteSubscription Интерефейс с методами к базе данных,
// который использует хендлер.
type deleteSubscription interface {
	DeleteSubscription(ctx context.Context, id int64) error
}

// helper Интерефейс с методами к Service,
// который использует хендлер.
type helper interface {
	GetSubID(r *http.Request) (int64, error)
	CheckSubscriptionID(ctx context.Context, subID int64) (bool, error)
}

// @Summary Удалить подписку
// @Description Удаляет подписку по указанному ID
// @Tags subscriptions
// @Produce plain
// @Param id path int true "ID удаляемой подписки" Example(123)
// @Success 204 "Подписка успешно удалена"
// @Failure 400 {string} string "Невалидный ID подписки"
// @Failure 404 {string} string "Подписка не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func Handler(
	l *slog.Logger, ds deleteSubscription, h helper,
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

	err = ds.DeleteSubscription(r.Context(), intsubID)
	if err != nil {
		log.Error("failed to delete subscription from DB", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)

	log.Info("Subscription delete successfully!", slog.Int64("ID", intsubID))
}
