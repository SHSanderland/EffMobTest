// Пакет model описывает основные модели, использующиеся
// в сервисе, а также функции работы с ними.
package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Subscription Структура подписки.
//
// Примечание: StartDate/EndDate можно было использовать с типом
// time.Time, но так как в запросе они передаются строкой, то
// решено было оставить их также строкой и преобразовывать в нужный
// вид по мере необходимости.
type Subscription struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date,omitempty"`
}

// GetSubFromBody Получения тела запроса и маршал в Subscription.
func GetSubFromBody(r *http.Request) (*Subscription, error) {
	sub := Subscription{}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		return nil, fmt.Errorf("bad user body: %w", err)
	}

	return &sub, nil
}

// IsValidSubscription Валидация структуры Subscription.
func IsValidSubscription(sub *Subscription) bool {
	if sub.ServiceName == "" || sub.Price <= 0 {
		return false
	}

	if uuid.Validate(sub.UserID.String()) != nil {
		return false
	}

	if sub.EndDate != "" {
		layout := "01-2006"

		startDate, err := time.Parse(layout, sub.StartDate)
		if err != nil {
			return false
		}

		endDate, err := time.Parse(layout, sub.EndDate)
		if err != nil {
			return false
		}

		if !endDate.After(startDate) {
			return false
		}
	}

	return true
}

// CostParams Структура для хендлера CostSubscription
// с фильтрующими данными.
type CostParams struct {
	ServiceName string     `json:"service_name"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}
