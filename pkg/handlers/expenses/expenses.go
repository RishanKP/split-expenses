package expenses

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"time"

	"split-expenses/library/api"
	"split-expenses/library/jwt"
	"split-expenses/pkg/interfaces"
	"split-expenses/pkg/repository"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Download(c *gin.Context)
}

type handler struct {
	repo     repository.ExpenseRepository
	userRepo repository.UserRepository
}

func (h handler) Create(c *gin.Context) {

	var req interfaces.CreateExpenseInput

	err := c.BindJSON(&req)
	if err != nil {
		api.NewClientError(c, http.StatusBadRequest, "invalid request")
		return
	}

	input, err := req.AsExpense()
	if err != nil {
		api.NewClientError(c, http.StatusBadRequest, err.Error())
		return
	}

	for _, u := range input.Participants {
		_, err := h.userRepo.GetById(context.TODO(), u.UserID)
		if err != nil {
			api.NewClientError(c, http.StatusBadRequest, "invalid participant id(s)")
			return
		}
	}

	err = h.repo.CreateExpense(input)
	if err != nil {
		api.NewInternalError(c, http.StatusInternalServerError, "failed to create expense")
		return
	}

	api.Result(c, http.StatusOK, "success", struct{}{})
}

func (h handler) Get(c *gin.Context) {
	claims, exists := c.Get("user")
	if !exists {
		api.NewClientError(c, http.StatusUnauthorized, "claims not found")
		return
	}

	userClaims, ok := claims.(jwt.Claims)
	if !ok {
		api.NewInternalError(c, http.StatusInternalServerError, "failed to fetch claims")
		return
	}

	expenses, err := h.repo.GetAll(userClaims.UserId)
	if err != nil {
		api.NewInternalError(c, http.StatusInternalServerError, "failed to fetch documents")
		return
	}

	api.Result(c, http.StatusOK, "success", expenses)
}

func (h handler) Download(c *gin.Context) {
	claims, exists := c.Get("user")
	if !exists {
		api.NewClientError(c, http.StatusUnauthorized, "claims not found")
		return
	}

	userClaims, ok := claims.(jwt.Claims)
	if !ok {
		api.NewInternalError(c, http.StatusInternalServerError, "failed to fetch claims")
		return
	}

	expenses, err := h.repo.GetAll(userClaims.UserId)
	if err != nil {
		api.NewInternalError(c, http.StatusInternalServerError, "failed to fetch documents")
		return
	}

	var csvData [][]string
	csvData = append(csvData, []string{"Expense ID", "Description", "Total Amount", "Split Type", "Participant User ID", "Participant Amount", "Participant Percentage", "Created At"})

	for _, expense := range expenses {
		createdAt := expense.CreatedAt.Format(time.RFC3339)
		for _, participant := range expense.Participants {
			csvData = append(csvData, []string{
				expense.ID.Hex(),
				expense.Description,
				fmt.Sprintf("%.2f", expense.Amount),
				expense.SplitType,
				participant.UserID,
				fmt.Sprintf("%.2f", participant.Amount),
				fmt.Sprintf("%.2f", participant.Percentage),
				createdAt,
			})
		}
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.WriteAll(csvData); err != nil {
		log.Println("Error writing CSV data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV"})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=balance_sheet.csv")

	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}

func Newhandler(repo repository.ExpenseRepository, userRepo repository.UserRepository) ExpenseHandler {
	return handler{
		repo:     repo,
		userRepo: userRepo,
	}
}
