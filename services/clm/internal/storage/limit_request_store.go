package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"services/clm/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// LimitRequestStore defines the interface for database operations on limit requests.
type LimitRequestStore interface {
	Create(ctx context.Context, request *models.LimitRequest) (*models.LimitRequest, error)
	FindTeamLeadByDepartment(ctx context.Context, departmentID uuid.UUID) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateApprovalStep(ctx context.Context, step *models.ApprovalStep) error
	CreateAuditLog(ctx context.Context, entry *models.AuditLog) error
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

// FindTeamLeadByDepartment finds a team lead for the given department (R-2).
// Note: This is a placeholder implementation as our current schema doesn't have departments.
func (s *DBStore) FindTeamLeadByDepartment(ctx context.Context, departmentID uuid.UUID) (*models.User, error) {
	// For now, return a placeholder error since department functionality is not yet implemented
	return nil, fmt.Errorf("department functionality not yet implemented")
}

// GetUserByID retrieves a user by ID.
func (s *DBStore) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, external_id, email, name, role, created_at, updated_at 
		FROM users 
		WHERE id = $1`

	var user models.User
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.ExternalID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// CreateApprovalStep creates a new approval step.
// Note: This is a placeholder implementation.
func (s *DBStore) CreateApprovalStep(ctx context.Context, step *models.ApprovalStep) error {
	// Placeholder implementation - not yet connected to actual schema
	return fmt.Errorf("approval steps functionality not yet implemented")
}

// CreateAuditLog creates a new audit log entry.
// Note: This is a placeholder implementation.
func (s *DBStore) CreateAuditLog(ctx context.Context, entry *models.AuditLog) error {
	// Placeholder implementation - not yet connected to actual schema
	return fmt.Errorf("audit log functionality not yet implemented")
}
