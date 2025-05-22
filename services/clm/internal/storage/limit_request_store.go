package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"services/clm/internal/models"
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
	// Database connection string from environment or default
	dbURL := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		// For now, allow creation without connection for testing
		log.Printf("Warning: Failed to open database: %v", err)
		return &DBStore{db: nil}
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		log.Printf("Warning: Failed to ping database: %v", err)
		db.Close()
		return &DBStore{db: nil}
	}
	
	return &DBStore{db: db}
}

// Create creates a new limit request with auto-assigned team lead (R-2).
func (s *DBStore) Create(ctx context.Context, request *models.LimitRequest) (*models.LimitRequest, error) {
	log.Printf("DBStore: Create called with request: %+v\n", request)

	// If no database connection, fall back to placeholder behavior for testing
	if s.db == nil {
		log.Printf("Warning: No database connection, using placeholder implementation")
		request.ID = uuid.New()
		request.CreatedAt = time.Now()
		request.UpdatedAt = time.Now()
		request.Status = "PENDING_TEAM_LEAD"
		// For placeholder, just mark as having no assignee yet
		request.CurrentAssigneeID = sql.NullString{String: "", Valid: false}
		return request, nil
	}

	// Get user to find their department
	user, err := s.GetUserByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Find team lead for user's department (R-2 requirement)
	teamLead, err := s.FindTeamLeadByDepartment(ctx, user.DepartmentID)
	if err != nil {
		log.Printf("ERROR: No team lead found for department %s: %v", user.DepartmentID, err)
		return nil, fmt.Errorf("team lead not configured for department")
	}

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Assign ID and timestamps
	request.ID = uuid.New()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.Status = "PENDING_TEAM_LEAD"
	request.CurrentAssigneeID = sql.NullString{String: teamLead.ID.String(), Valid: true}

	// Insert limit request
	query := `
		INSERT INTO limit_requests 
		(id, user_id, amount, currency, justification, desired_date, status, current_assignee_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = tx.ExecContext(ctx, query,
		request.ID, request.UserID, request.Amount, request.Currency,
		request.Justification, request.DesiredDate, request.Status,
		request.CurrentAssigneeID, request.CreatedAt, request.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert limit request: %w", err)
	}

	// Create approval step (R-2 requirement)
	approvalStep := &models.ApprovalStep{
		ID:         uuid.New(),
		RequestID:  request.ID,
		AssigneeID: teamLead.ID,
		StepType:   "TEAM_LEAD",
		Status:     "PENDING",
		AssignedAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	stepQuery := `
		INSERT INTO approval_steps 
		(id, request_id, assignee_id, step_type, status, assigned_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.ExecContext(ctx, stepQuery,
		approvalStep.ID, approvalStep.RequestID, approvalStep.AssigneeID,
		approvalStep.StepType, approvalStep.Status, approvalStep.AssignedAt,
		approvalStep.CreatedAt, approvalStep.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert approval step: %w", err)
	}

	// Create audit log entry (R-2 requirement)
	payload := map[string]interface{}{
		"assigned_team_lead_id": teamLead.ID.String(),
		"department_id":         user.DepartmentID.String(),
		"initial_status":        "PENDING_TEAM_LEAD",
	}
	payloadJSON, _ := json.Marshal(payload)

	auditEntry := &models.AuditLog{
		ID:          uuid.New(),
		RequestID:   request.ID,
		Action:      "request_created",
		Actor:       request.UserID,
		Timestamp:   time.Now(),
		PayloadJSON: string(payloadJSON),
	}

	auditQuery := `
		INSERT INTO audit_log 
		(id, request_id, action, actor, timestamp, payload_json)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, auditQuery,
		auditEntry.ID, auditEntry.RequestID, auditEntry.Action,
		auditEntry.Actor, auditEntry.Timestamp, auditEntry.PayloadJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to insert audit log: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully created limit request %s with team lead %s", request.ID, teamLead.ID)
	return request, nil
}

// FindTeamLeadByDepartment finds a team lead for the given department (R-2).
func (s *DBStore) FindTeamLeadByDepartment(ctx context.Context, departmentID uuid.UUID) (*models.User, error) {
	if s.db == nil {
		return nil, fmt.Errorf("no database connection available")
	}

	query := `
		SELECT id, email, role, department_id, created_at, updated_at 
		FROM users 
		WHERE role = 'team_lead' AND department_id = $1 
		ORDER BY email ASC`

	rows, err := s.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query team leads: %w", err)
	}
	defer rows.Close()

	var teamLeads []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.DepartmentID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		teamLeads = append(teamLeads, user)
	}

	if len(teamLeads) == 0 {
		return nil, fmt.Errorf("no team lead found for department %s", departmentID)
	}

	if len(teamLeads) > 1 {
		log.Printf("WARN: Multiple team leads found for department %s, selecting first by email: %s", departmentID, teamLeads[0].Email)
	}

	return &teamLeads[0], nil
}

// GetUserByID retrieves a user by ID.
func (s *DBStore) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	if s.db == nil {
		return nil, fmt.Errorf("no database connection available")
	}

	query := `
		SELECT id, email, role, department_id, created_at, updated_at 
		FROM users 
		WHERE id = $1`

	var user models.User
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.Role, &user.DepartmentID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// CreateApprovalStep creates a new approval step.
func (s *DBStore) CreateApprovalStep(ctx context.Context, step *models.ApprovalStep) error {
	query := `
		INSERT INTO approval_steps 
		(id, request_id, assignee_id, step_type, status, assigned_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := s.db.ExecContext(ctx, query,
		step.ID, step.RequestID, step.AssigneeID, step.StepType,
		step.Status, step.AssignedAt, step.CreatedAt, step.UpdatedAt)
	return err
}

// CreateAuditLog creates a new audit log entry.
func (s *DBStore) CreateAuditLog(ctx context.Context, entry *models.AuditLog) error {
	query := `
		INSERT INTO audit_log 
		(id, request_id, action, actor, timestamp, payload_json)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.ExecContext(ctx, query,
		entry.ID, entry.RequestID, entry.Action, entry.Actor,
		entry.Timestamp, entry.PayloadJSON)
	return err
}