package storage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"services/clm/internal/models"
)

func TestDBStore_Create_Placeholder(t *testing.T) {
	store := NewDBStore()
	ctx := context.Background()

	userID := uuid.New()
	desiredDate := time.Now().Add(24 * time.Hour)

	request := &models.LimitRequest{
		UserID:        userID,
		Amount:        1000.00,
		Currency:      "USD",
		Justification: "Business travel expenses",
		DesiredDate:   desiredDate,
		Status:        "PENDING_TEAM_LEAD", // This will be overwritten by the Create method
	}

	createdRequest, err := store.Create(ctx, request)

	assert.NoError(t, err, "Expected no error from placeholder Create")
	assert.NotNil(t, createdRequest, "Expected a non-nil request from Create")

	// Assert that placeholder behavior is as expected
	assert.NotEqual(t, uuid.Nil, createdRequest.ID, "Expected a new ID to be generated")
	assert.Equal(t, request.UserID, createdRequest.UserID)
	assert.Equal(t, request.Amount, createdRequest.Amount)
	assert.Equal(t, request.Currency, createdRequest.Currency)
	assert.Equal(t, request.Justification, createdRequest.Justification)
	assert.Equal(t, request.DesiredDate, createdRequest.DesiredDate)
	// The placeholder Create method currently sets these:
	// request.Status is not set by the caller in the handler, but set before calling store.Create
	// For this test, we are testing the store's Create method directly.
	// The handler sets status to PENDING_TEAM_LEAD before calling Create.
	// The store's Create method does not currently modify status.
	assert.Equal(t, request.Status, createdRequest.Status)
	assert.False(t, createdRequest.CurrentApproverID.Valid, "Expected CurrentApproverID to be invalid for a new request")

	assert.WithinDuration(t, time.Now(), createdRequest.CreatedAt, 5*time.Second, "Expected CreatedAt to be recent")
	assert.WithinDuration(t, time.Now(), createdRequest.UpdatedAt, 5*time.Second, "Expected UpdatedAt to be recent")

	// Verify that the input request object was modified in place (which is current behavior)
	assert.Equal(t, createdRequest.ID, request.ID, "Expected input request ID to be updated")
	assert.Equal(t, createdRequest.CreatedAt, request.CreatedAt, "Expected input request CreatedAt to be updated")
	assert.Equal(t, createdRequest.UpdatedAt, request.UpdatedAt, "Expected input request UpdatedAt to be updated")
}
