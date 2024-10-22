package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"split-expenses/library/config"
	"split-expenses/library/db"
	"split-expenses/pkg/handlers/expenses"
	"split-expenses/pkg/handlers/user"
	"split-expenses/pkg/interfaces"
	"split-expenses/pkg/middleware"
	"split-expenses/pkg/models"
	"split-expenses/pkg/repository"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	db.Connect()
	defer db.Disconnect()
	router := gin.Default()

	dataBase := db.Client.Database(config.DB_NAME)
	repo := repository.NewExpenseRepository(dataBase, "expenses")
	userRepo := repository.NewUserRepository(dataBase, "user")
	handler := expenses.Newhandler(repo, userRepo)
	userHandler := user.Newhandler(userRepo)

	router.POST("/login", userHandler.Login)
	router.POST("/add-expense", middleware.AuthMiddleware("user"), handler.Create)
	return router
}

func getUserToken(router *gin.Engine) string {
	input := interfaces.LoginCredentials{
		Email:    "rishan@gmail.com",
		Password: "rish123",
	}

	jsonData, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/add-expense", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getUserToken(router)) // Add your authorization token

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	type response struct {
		Data struct {
			Email     string `json:"email"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Token     string `json:"token"`
			Contact   string `json:"contact"`
			ID        string `json:"id"`
		} `json:"data"`
		Status string `json:"status"`
	}

	var res response
	json.Unmarshal(w.Body.Bytes(), &res)

	return res.Data.Token
}

func TestCreateExpense_Success(t *testing.T) {
	router := setupRouter()

	input := interfaces.CreateExpenseInput{
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
	req.Header.Set("Authorization", "Bearer "+getUserToken(router)) // Add your authorization token

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	type response struct {
		Data   interface{} `json:"data"`
		Status string      `json:"status"`
	}

	var res response
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "success", res.Status)
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
