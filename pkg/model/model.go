package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date,omitempty"`
}

func GetSubFromBody(r *http.Request) (*Subscription, error) {
	sub := Subscription{}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		return nil, fmt.Errorf("bad user body: %w", err)
	}

	return &sub, nil
}

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
