package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"split-expenses/pkg/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/add-expense", Authorize(), AddExpense) // Assuming you have an Authorize middleware
	return router
}

func TestCreateExpense_Success(t *testing.T) {
	router := setupRouter()

	input := CreateExpenseInput{
		Amount:      150.00,
		SplitType:   "exact",
		Description: "Dinner with friends",
		Participants: []models.Participant{
			{UserID: "user123", Amount: 50.00, Percentage: 33.33},
			{UserID: "user456", Amount: 50.00, Percentage: 33.33},
			{UserID: "user789", Amount: 50.00, Percentage: 33.34},
		},
	}

	jsonData, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/add-expense", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer your_token_here") // Add your authorization token

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Expense
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Dinner with friends", response.Description)
	assert.Equal(t, 150.00, response.Amount)
}

func TestCreateExpense_BadRequest(t *testing.T) {
	router := setupRouter()

	// Missing required fields
	invalidInput := `{"splitType":"exact","description":"Dinner with friends"}`
	req, _ := http.NewRequest("POST", "/add-expense", bytes.NewBuffer([]byte(invalidInput)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer your_token_here") // Add your authorization token

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
