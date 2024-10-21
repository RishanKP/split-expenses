package api

import (
	"github.com/gin-gonic/gin"
)

func NewClientError(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, gin.H{
		"status":  "failed",
		"message": msg,
	})
}

func NewInternalError(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, gin.H{
		"status":  "failed",
		"message": msg,
	})
}

func Result(c *gin.Context, statusCode int, msg string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status": "success",
		"data":   data,
	})
}
