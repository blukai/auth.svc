package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func Provider() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		ses := sessions.Default(c)
		if ses.Get(provider) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "already logged in",
			})
			return
		}

		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func ProviderCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("could not complete auth: %v", err),
			})
			return
		}

		v, err := json.Marshal(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("could not marshal session: %v", err),
			})
			return
		}

		provider := c.Param("provider")
		ses := sessions.Default(c)
		ses.Set(provider, string(v))
		if err := ses.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("could not save session: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func ProviderLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := gothic.Logout(c.Writer, c.Request); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("could not logout: %v", err),
			})
			return
		}

		provider := c.Param("provider")
		ses := sessions.Default(c)
		ses.Delete(provider)
		if err := ses.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("could not delete session id: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func ProviderUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		ses := sessions.Default(c)
		v, _ := ses.Get(provider).(string)
		if v == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": strings.ToLower(http.StatusText(http.StatusUnauthorized)),
			})
			return
		}

		var user goth.User
		if err := json.Unmarshal([]byte(v), &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("could not unmarshal session: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, &user)
	}
}
