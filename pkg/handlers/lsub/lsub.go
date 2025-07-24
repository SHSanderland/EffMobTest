// Пакет lsub для хендлера ListSubscription.
package lsub

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// listSubscription Интерефейс с методами к базе данных,
// который использует хендлер.
type listSubscription interface {
	GetListSubscription(ctx context.Context, userID uuid.UUID, serviceName string) ([]*model.Subscription, error)
}

// urlParser Интерефейс с методами к Service,
// который использует хендлер.
type urlParser interface {
	GetUserIDAndServiceName(r *http.Request) (uuid.UUID, string, error)
}

// userResponse Структура для ответа пользователю.
type userResponse struct {
	Subscriptions []*model.Subscription `json:"subscriptions"`
	Total         int                   `json:"total"`
}

// @Summary		Получить список подписок
// @Description	Возвращает список подписок с возможностью фильтрации по user_id и service_name
// @Tags			subscriptions
// @Produce		json
// @Param			user_id			query		string			true	"UUID пользователя для фильтрации"	Example(550e8400-e29b-41d4-a716-446655440000)
// @Param			service_name	query		string			true	"Название сервиса для фильтрации"	Example(netflix)
// @Success		200				{object}	userResponse	"Успешный запрос"
// @Failure		400				{string}	string			"Невалидные параметры запроса"
// @Failure		500				{string}	string			"Внутренняя ошибка сервера"
// @Router			/subscriptions [get]
func Handler(
	l *slog.Logger, ls listSubscription, up urlParser,
	w http.ResponseWriter, r *http.Request,
) {
	const fn = "handlers.lsub.Handler"
	log := l.With(
		slog.String("fn", fn),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)

	userID, serviceName, err := up.GetUserIDAndServiceName(r)
	if err != nil {
		log.Error("failed to parse url", slog.String("err", err.Error()))
		http.Error(w, "Invalid URL params", http.StatusBadRequest)

		return
	}

	subs, err := ls.GetListSubscription(r.Context(), userID, serviceName)
	if err != nil {
		log.Error("failed to get list subscriptions", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	userResp := userResponse{subs, len(subs)}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(userResp); err != nil {
		log.Error("failed to encode json", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	log.Info(
		"List of subscriptions sended successfully!",
		slog.String("userID", userID.String()),
		slog.String("serviceName", serviceName),
	)
}
