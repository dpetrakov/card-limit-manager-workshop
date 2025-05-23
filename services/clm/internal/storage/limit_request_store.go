package storage

import (
	"context"
	"database/sql"
	"log"
	"os"

	"services/clm/internal/models"

	_ "github.com/lib/pq"
)

// LimitRequestStore defines the interface for database operations on limit requests.
type LimitRequestStore interface {
	Create(ctx context.Context, request *models.LimitRequest) (*models.LimitRequest, error)
}

// DBStore is a basic implementation of LimitRequestStore.
type DBStore struct {
	db *sql.DB
}

// NewDBStore creates a new DBStore with database connection.
func NewDBStore() *DBStore {
	// Get database connection string from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return &DBStore{db: db}
}

// Create inserts a new limit request into the database.
func (s *DBStore) Create(ctx context.Context, request *models.LimitRequest) (*models.LimitRequest, error) {
	log.Printf("DBStore: Create called with request: %+v\n", request)

	query := `
		INSERT INTO limit_requests (user_id, amount, currency, justification, desired_date, status, current_approver_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	var currentApproverID interface{}
	if request.CurrentApproverID.Valid {
		currentApproverID = request.CurrentApproverID.UUID
	} else {
		currentApproverID = nil
	}

	err := s.db.QueryRowContext(ctx, query,
		request.UserID,
		request.Amount,
		request.Currency,
		request.Justification,
		request.DesiredDate,
		request.Status,
		currentApproverID,
	).Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt)

	if err != nil {
		log.Printf("Failed to insert limit request: %v", err)
		return nil, err
	}

	log.Printf("Successfully created limit request with ID: %s", request.ID)
	return request, nil
}
