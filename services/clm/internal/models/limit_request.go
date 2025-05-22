package models

import (
	"time"

	"github.com/google/uuid"
)

// LimitRequestCreate represents the data required to create a new limit request.
type LimitRequestCreate struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Justification string  `json:"justification"`
	DesiredDate string  `json:"desired_date"` // Expected format: YYYY-MM-DD
}

// LimitRequestView represents a limit request as returned by the API.
type LimitRequestView struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Justification string    `json:"justification"`
	DesiredDate   string    `json:"desired_date"`
	Status            string    `json:"status"`
	CurrentApproverID *string   `json:"current_approver_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// LimitRequest represents a limit request in the database.
type LimitRequest struct {
	ID                uuid.UUID    `db:"id"`
	UserID            uuid.UUID    `db:"user_id"`
	Amount            float64      `db:"amount"`
	Currency          string       `db:"currency"`
	Justification     string       `db:"justification"`
	DesiredDate       time.Time    `db:"desired_date"`
	Status            string       `db:"status"`
	CurrentApproverID uuid.NullUUID `db:"current_approver_id"` // Nullable UUID
	CreatedAt         time.Time    `db:"created_at"`
	UpdatedAt         time.Time    `db:"updated_at"`
}
