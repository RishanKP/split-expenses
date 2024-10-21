package routes

import (
	"split-expenses/pkg/handlers/user"
	"split-expenses/pkg/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(r *gin.Engine, db *mongo.Database) {
	repo := repository.NewUserRepository(db, "user")
	handler := user.Newhandler(repo)

	userGroup := r.Group("/user")
	{
		userGroup.POST("/signup", handler.SignUp)
		userGroup.POST("/login", handler.Login)
		userGroup.GET("/:id", handler.GetUserById)
	}
}
