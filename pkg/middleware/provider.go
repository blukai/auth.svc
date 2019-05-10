package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

func ProviderName() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
		c.Next()
	}
}
