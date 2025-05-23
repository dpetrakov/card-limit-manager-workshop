package handlers

import (
	"net/http"
	"time"

	"services/clm/internal/models"
	"services/clm/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateLimitRequestHandler godoc
// @Summary Create a new limit request (with R-2 auto team lead assignment)
// @Description Create a new limit request with the input payload. Team lead will be auto-assigned based on user's department.
// @Tags requests
// @Accept  json
// @Produce  json
// @Param   limitRequest body models.LimitRequestCreate true "Create Limit Request Payload"
// @Success 201 {object} models.LimitRequestView "Successfully created limit request"
// @Failure 400 {object} models.ErrorResponse "Invalid request payload or data"
// @Failure 500 {object} models.ErrorResponse "Internal server error or specific business logic error (e.g., team lead not configured)"
// @Router /requests [post]
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
	// For demo purposes, use a fixed test user. In real app, get from auth context.
	testUserID, _ := uuid.Parse("a005d32d-6190-477c-b23e-38c44eaaaae0")
	limitRequest := &models.LimitRequest{
		UserID:        testUserID, // Use existing test user
		Amount:        req.Amount,
		Currency:      req.Currency,
		Justification: req.Justification,
		DesiredDate:   desiredDate,
	}

	// Instantiate store and call Create (implements R-2 auto team lead assignment)
	store := storage.NewDBStore()
	createdRequest, err := store.Create(c.Request.Context(), limitRequest)
	if err != nil {
		// Handle R-2 specific error case (no team lead configured)
		if err.Error() == "team lead not configured for department" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Team lead not configured for department"})
			return
		}
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

	// Set current_assignee_id from R-2 logic
	if createdRequest.CurrentAssigneeID.Valid {
		assigneeIDStr := createdRequest.CurrentAssigneeID.String
		response.CurrentAssigneeID = &assigneeIDStr
	}

	c.JSON(http.StatusCreated, response)
}
