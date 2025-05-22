package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"services/clm/internal/models"
	"services/clm/internal/storage"
)

// MockLimitRequestStore is a mock implementation of LimitRequestStore for testing.
type MockLimitRequestStore struct {
	CreateFunc func(ctx context.Context, request *models.LimitRequest) (*models.LimitRequest, error)
}

func (m *MockLimitRequestStore) Create(ctx context.Context, request *models.LimitRequest) (*models.LimitRequest, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, request)
	}
	// Default behavior if CreateFunc is not set
	request.ID = uuid.New()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	return request, nil
}

// SetupRouter sets up a Gin router for testing and injects the mock store.
// This function is not directly used in the CreateLimitRequestHandler tests as the handler
// currently instantiates its own store. When dependency injection for the store is implemented
// in the handler, this setup will be more relevant.
// For now, we will test CreateLimitRequestHandler directly.
func SetupRouterWithStore(store storage.LimitRequestStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	// This is how we would inject the store if CreateLimitRequestHandler accepted it.
	// For now, CreateLimitRequestHandler creates its own store, so this mock isn't fully utilized
	// in the way a typical DI setup would work. The tests will work around this.
	router.POST("/requests", CreateLimitRequestHandler) // CreateLimitRequestHandler will use its own store
	return router
}

func TestCreateLimitRequestHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Note: CreateLimitRequestHandler currently creates its own DBStore instance.
	// To properly test with a mock, the handler would need to accept a store via DI.
	// The current test will therefore exercise the actual DBStore's placeholder Create.

	router := gin.Default()
	router.POST("/requests", CreateLimitRequestHandler)

	payload := models.LimitRequestCreate{
		Amount:      1500.50,
		Currency:    "EUR",
		Justification: "New equipment purchase",
		DesiredDate: "2025-12-01",
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected HTTP status 201 Created")

	var responseView models.LimitRequestView
	err := json.Unmarshal(rr.Body.Bytes(), &responseView)
	assert.NoError(t, err, "Error unmarshalling response body")

	assert.NotEmpty(t, responseView.ID, "Expected ID to be populated")
	assert.NotEmpty(t, responseView.UserID, "Expected UserID to be populated (placeholder)")
	assert.Equal(t, payload.Amount, responseView.Amount)
	assert.Equal(t, payload.Currency, responseView.Currency)
	assert.Equal(t, payload.Justification, responseView.Justification)
	assert.Equal(t, payload.DesiredDate, responseView.DesiredDate)
	assert.Equal(t, "PENDING_TEAM_LEAD", responseView.Status)
	assert.Nil(t, responseView.CurrentApproverID, "Expected CurrentApproverID to be nil for a new request")
	assert.WithinDuration(t, time.Now(), responseView.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now(), responseView.UpdatedAt, 5*time.Second)

	// Check that current_approver_id is not present in the raw JSON response
	var rawResponse map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &rawResponse)
	assert.NoError(t, err, "Error unmarshalling raw response body")
	_, currentApproverIDExists := rawResponse["current_approver_id"]
	assert.False(t, currentApproverIDExists, "Expected 'current_approver_id' key to be absent in JSON response")
}

func TestCreateLimitRequestHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/requests", CreateLimitRequestHandler)

	// Malformed JSON
	req, _ := http.NewRequest(http.MethodPost, "/requests", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request payload", "Expected error message for malformed JSON")

	// Missing required fields (e.g., Amount)
	payload := map[string]string{
		"currency":    "USD",
		"justification": "Test",
		"desired_date": "2025-10-10",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ = http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Amount must be positive", "Expected error message for missing amount")
}


func TestCreateLimitRequestHandler_InvalidPayloadData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/requests", CreateLimitRequestHandler)

	testCases := []struct {
		name          string
		payload       interface{}
		expectedError string
	}{
		{
			name: "Negative Amount",
			payload: models.LimitRequestCreate{Amount: -100, Currency: "USD", Justification: "Test", DesiredDate: "2025-01-01"},
			expectedError: "Amount must be positive",
		},
		{
			name: "Missing Currency",
			payload: models.LimitRequestCreate{Amount: 100, Justification: "Test", DesiredDate: "2025-01-01"},
			expectedError: "Currency is required",
		},
		{
			name: "Missing Justification",
			payload: models.LimitRequestCreate{Amount: 100, Currency: "USD", DesiredDate: "2025-01-01"},
			expectedError: "Justification is required",
		},
		{
			name: "Missing DesiredDate",
			payload: models.LimitRequestCreate{Amount: 100, Currency: "USD", Justification: "Test"},
			expectedError: "Desired date is required",
		},
		{
			name: "Invalid DesiredDate Format",
			payload: models.LimitRequestCreate{Amount: 100, Currency: "USD", Justification: "Test", DesiredDate: "01-01-2025"},
			expectedError: "Invalid desired_date format, expected YYYY-MM-DD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedError)
		})
	}
}


// TestCreateLimitRequestHandler_StorageError simulates a storage error.
// This test requires modification of CreateLimitRequestHandler to allow injection of a mock store.
// Since CreateLimitRequestHandler currently creates its own DBStore, this test, as written,
// cannot directly inject a mock that returns an error.
// For a true unit test of this scenario, the handler needs to be refactored for DI.
// The following is a conceptual test.
func TestCreateLimitRequestHandler_StorageError(t *testing.T) {
	// This test is more of an integration test with the current handler implementation,
	// or it needs the handler to be refactored.
	// If store.Create were to actually fail (e.g. DB down), this would be the expected path.
	// For now, the placeholder DBStore.Create always succeeds.
	// To make this test fail the store.Create, we'd need to modify DBStore.Create
	// to simulate an error, or (preferably) use DI.

	// As a placeholder, we acknowledge this limitation.
	// If CreateLimitRequestHandler is refactored to take `storage.LimitRequestStore` as a parameter,
	// this test would look like:
	/*
	   gin.SetMode(gin.TestMode)
	   router := gin.New() // Use gin.New() for more control if needed

	   mockStore := &MockLimitRequestStore{
	       CreateFunc: func(ctx context.Context, request *models.LimitRequest) (*models.LimitRequest, error) {
	           return nil, errors.New("simulated database error")
	       },
	   }

	   // Assumes CreateLimitRequestHandler is refactored to accept a store
	   // e.g., router.POST("/requests", handlers.NewCreateLimitRequestHandler(mockStore))
	   // For now, this test will use the existing handler which creates its own store.
	   // We can't easily make the internal store fail without changing its code for testing purposes.

	   // So, this test case cannot be effectively implemented without handler refactoring.
	   // We'll skip trying to make the real store fail for this unit test.
	   t.Skip("Skipping storage error test: Requires handler refactoring for store DI to mock storage errors.")


	   // If DI was in place, the rest of the test would be:
	   payload := models.LimitRequestCreate{
	       Amount:      2000.00,
	       Currency:    "GBP",
	       Justification: "Emergency fund",
	       DesiredDate: "2025-11-15",
	   }
	   payloadBytes, _ := json.Marshal(payload)

	   req, _ := http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(payloadBytes))
	   req.Header.Set("Content-Type", "application/json")
	   rr := httptest.NewRecorder()

	   // This assumes the router is set up with a handler that uses the mockStore
	   // router.ServeHTTP(rr, req) // This line would use the DI'd handler

	   // assert.Equal(t, http.StatusInternalServerError, rr.Code)
	   // assert.Contains(t, rr.Body.String(), "Failed to create limit request")
	   // assert.Contains(t, rr.Body.String(), "simulated database error")
	*/
	t.Log("TestCreateLimitRequestHandler_StorageError: This test is currently a placeholder as the handler does not support DI for the store. A real storage error cannot be easily simulated for this specific handler unit test without refactoring the handler or the store itself.")
}
