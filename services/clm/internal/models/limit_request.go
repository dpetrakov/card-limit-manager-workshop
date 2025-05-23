package models

import (
	"time"

	"database/sql"

	"github.com/google/uuid"
)

// LimitRequestCreate represents the data required to create a new limit request.
type LimitRequestCreate struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Justification string  `json:"justification"`
	DesiredDate   string  `json:"desired_date"` // Expected format: YYYY-MM-DD
}

// LimitRequestView represents a limit request as returned by the API.
type LimitRequestView struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	Justification     string    `json:"justification"`
	DesiredDate       string    `json:"desired_date"`
	Status            string    `json:"status"`
	CurrentAssigneeID *string   `json:"current_assignee_id,omitempty"` // Changed to match spec
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// LimitRequest represents a limit request in the database.
type LimitRequest struct {
	ID                uuid.UUID      `db:"id"`
	UserID            uuid.UUID      `db:"user_id"`
	Amount            float64        `db:"amount"`
	Currency          string         `db:"currency"`
	Justification     string         `db:"justification"`
	DesiredDate       time.Time      `db:"desired_date"`
	Status            string         `db:"status"`
	CurrentApproverID sql.NullString `db:"current_approver_id"` // Fixed: matches DB schema
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
}

// User represents a user in the system.
type User struct {
	ID           uuid.UUID `db:"id"`
	ExternalID   string    `db:"external_id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	Role         string    `db:"role"`
	DepartmentID uuid.UUID `db:"department_id"` // Will be added in future migration for R-2
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ApprovalStep represents an approval step in the workflow.
type ApprovalStep struct {
	ID         uuid.UUID `db:"id"`
	RequestID  uuid.UUID `db:"request_id"`
	AssigneeID uuid.UUID `db:"assignee_id"`
	StepType   string    `db:"step_type"`
	Status     string    `db:"status"`
	AssignedAt time.Time `db:"assigned_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// AuditLog represents an audit log entry.
type AuditLog struct {
	ID          uuid.UUID `db:"id"`
	RequestID   uuid.UUID `db:"request_id"`
	Action      string    `db:"action"`
	Actor       uuid.UUID `db:"actor"`
	Timestamp   time.Time `db:"timestamp"`
	PayloadJSON string    `db:"payload_json"`
}

// ErrorResponse represents a generic error response for API failures.
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"Detailed error message here"`
}
