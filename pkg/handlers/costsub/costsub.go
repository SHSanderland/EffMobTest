// Пакет costsub для хендлера CostSubscription.
package costsub

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/go-chi/chi/v5/middleware"
)

// costSubscription Интерефейс с методами к базе данных,
// который использует хендлер.
type costSubscription interface {
	CostSubscription(ctx context.Context, filter *model.CostParams) (int64, error)
}

// helper Интерефейс с методами к Service,
// который использует хендлер.
type helper interface {
	GetCostParams(r *http.Request) (*model.CostParams, error)
}

// userResponse Структура для ответа пользователю.
type userResponse struct {
	ServiceName string     `json:"service_name"`
	StartDate   time.Time  `json:"start_period"`
	EndDate     *time.Time `json:"end_period"`
	TotalCost   int64      `json:"total_cost"`
}

// @Summary		Рассчитать стоимость подписок
// @Description	Возвращает суммарную стоимость подписок за указанный период с возможностью фильтрации
// @Tags			subscriptions
// @Produce		json
// @Param			user_id			query		string			true	"UUID пользователя"					Example(550e8400-e29b-41d4-a716-446655440000)
// @Param			service_name	query		string			true	"Название сервиса"					Example(netflix)
// @Param			start_date		query		string			true	"Начало периода (формат MM-YYYY)"	Example(01-2023)
// @Param			end_date		query		string			true	"Конец периода (формат MM-YYYY)"	Example(12-2023)
// @Success		200				{object}	userResponse	"Успешный расчет стоимости"
// @Failure		400				{string}	string			"Невалидные параметры запроса"
// @Failure		500				{string}	string			"Внутренняя ошибка сервера"
// @Router			/subscriptions/cost [get]
func Handler(
	l *slog.Logger, cs costSubscription, h helper,
	w http.ResponseWriter, r *http.Request,
) {
	const fn = "handlers.costsub.Handler"
	log := l.With(
		slog.String("fn", fn),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)

	filters, err := h.GetCostParams(r)
	if err != nil {
		log.Error("invalid params", slog.String("err", err.Error()))
		http.Error(w, "Invalid params", http.StatusBadRequest)

		return
	}

	total, err := cs.CostSubscription(r.Context(), filters)
	if err != nil {
		log.Error("failed to count total cost", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	ur := userResponse{
		ServiceName: filters.ServiceName,
		StartDate:   filters.StartDate,
		EndDate:     filters.EndDate,
		TotalCost:   total,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&ur); err != nil {
		log.Error("failed to send response", slog.String("err", err.Error()))
		http.Error(w, "Something wrong", http.StatusInternalServerError)

		return
	}

	log.Info(
		"Total cost counted successfully!",
		slog.String("userID", filters.UserID.String()),
		slog.String("serviceName", filters.ServiceName),
	)
}
