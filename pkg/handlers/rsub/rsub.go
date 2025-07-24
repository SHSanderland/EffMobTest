// Пакет rsub для хендлера ReadSubscription.
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

// readSubscription Интерефейс с методами к базе данных,
// который использует хендлер.
type readSubscription interface {
	ReadSubscription(ctx context.Context, id int64) (*model.Subscription, error)
}

// helper Интерефейс с методами к Service,
// который использует хендлер.
type helper interface {
	GetSubID(r *http.Request) (int64, error)
	CheckSubscriptionID(ctx context.Context, subID int64) (bool, error)
}

// @Summary Получить подписку по ID
// @Description Возвращает информацию о подписке по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки" Example(123)
// @Success 200 {object} model.Subscription "Успешный запрос"
// @Failure 400 {string} string "Невалидный ID подписки"
// @Failure 404 {string} string "Подписка не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func Handler(
	l *slog.Logger, rs readSubscription, h helper,
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
