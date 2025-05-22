package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"services/clm/internal/models"
	"services/clm/internal/storage"
)

// CreateLimitRequestHandler handles the creation of new limit requests.
func CreateLimitRequestHandler(c *gin.Context) {
	var req models.LimitRequestCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Basic validation
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}
	if req.Currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Currency is required"})
		return
	}
	if req.Justification == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Justification is required"})
		return
	}
	if req.DesiredDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Desired date is required"})
		return
	}

	// Parse DesiredDate
	desiredDate, err := time.Parse("2006-01-02", req.DesiredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid desired_date format, expected YYYY-MM-DD"})
		return
	}

	// Construct LimitRequest object
	limitRequest := &models.LimitRequest{
		UserID:        uuid.New(), // Placeholder UserID
		Amount:        req.Amount,
		Currency:      req.Currency,
		Justification: req.Justification,
		DesiredDate:   desiredDate,
		Status:        "PENDING_TEAM_LEAD",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Instantiate store and call Create
	// In a real app, the store would be injected via dependency injection.
	store := storage.NewDBStore()
	createdRequest, err := store.Create(c.Request.Context(), limitRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create limit request: " + err.Error()})
		return
	}

	// Transform to LimitRequestView
	response := &models.LimitRequestView{
		ID:            createdRequest.ID.String(),
		UserID:        createdRequest.UserID.String(),
		Amount:        createdRequest.Amount,
		Currency:      createdRequest.Currency,
		Justification: createdRequest.Justification,
		DesiredDate:   createdRequest.DesiredDate.Format("2006-01-02"),
		Status:        createdRequest.Status,
		CreatedAt:     createdRequest.CreatedAt,
		UpdatedAt:     createdRequest.UpdatedAt,
	}

	if createdRequest.CurrentApproverID.Valid {
		approverIDStr := createdRequest.CurrentApproverID.UUID.String()
		response.CurrentApproverID = &approverIDStr
	}

	c.JSON(http.StatusCreated, response)
}
