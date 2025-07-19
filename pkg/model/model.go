package model

import (
	"encoding/json"
	"fmt"
	"net/http"

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
