package middleware

import (
	"net/http"
	"strings"

	"split-expenses/library/api"
	"split-expenses/library/jwt"
	"split-expenses/library/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header["Authorization"][0]
		if utils.IsEmpty(authHeader) {
			api.NewClientError(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.VerifyToken(tokenString, role)
		if err != nil {
			api.NewClientError(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}
