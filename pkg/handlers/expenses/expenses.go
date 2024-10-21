package expenses

import (
	"context"
	"net/http"

	"split-expenses/library/api"
	"split-expenses/library/jwt"
	"split-expenses/pkg/interfaces"
	"split-expenses/pkg/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpenseHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
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
		if err == mongo.ErrNoDocuments {
			api.NewClientError(c, http.StatusBadRequest, "invalid participant id")
			return
		}
	}

	err = h.repo.CreateExpense(input)
	if err != nil {
		api.NewInternalError(c, http.StatusInternalServerError, "failed to create expense")
		return
	}

	api.Result(c, http.StatusOK, "success", input)
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

func Newhandler(repo repository.ExpenseRepository, userRepo repository.UserRepository) ExpenseHandler {
	return handler{
		repo:     repo,
		userRepo: userRepo,
	}
}
