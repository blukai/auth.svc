package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

func ProviderName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("provider")
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", name))
		c.Next()
	}
}
