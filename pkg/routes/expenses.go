package routes

import (
	"split-expenses/pkg/handlers/expenses"
	"split-expenses/pkg/middleware"
	"split-expenses/pkg/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ExpensesRoutes(r *gin.Engine, db *mongo.Database) {
	repo := repository.NewExpenseRepository(db, "expenses")
	userRepo := repository.NewUserRepository(db, "users")
	handler := expenses.Newhandler(repo, userRepo)

	userGroup := r.Group("/expense")
	userGroup.Use(middleware.AuthMiddleware("admin"))
	{
		userGroup.POST("/", handler.Create)
		userGroup.GET("/", handler.Get)
	}
}
