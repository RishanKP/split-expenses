package main

import (
	"split-expenses/library/config"
	"split-expenses/library/db"
	"split-expenses/pkg/middleware"
	"split-expenses/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	defer db.Disconnect()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	routes.RegisterRoutes(r, db.Client.Database(config.DB_NAME))

	r.Run(":8080")
}
