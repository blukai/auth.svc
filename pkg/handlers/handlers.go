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

func Root(c *gin.Context) {
	name := c.Param("provider")
	if name == "user" {
		User(c)
		return
	}
	Provider(c)
}

func Provider(c *gin.Context) {
	ses := sessions.Default(c)
	if ses.Get("id") == nil {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func ProviderCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("could not complete auth: %v", err),
		})
		return
	}

	v, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("could not marshal: %v", err),
		})
		return
	}

	ses := sessions.Default(c)
	ses.Set("id", string(v))
	if err := ses.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("could not save: %v", err),
		})
		return
	}

	// TODO: redirect
}

func User(c *gin.Context) {
	ses := sessions.Default(c)
	v, _ := ses.Get("id").(string)
	if v == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": strings.ToLower(http.StatusText(http.StatusUnauthorized)),
		})
		return
	}

	var user goth.User
	if err := json.Unmarshal([]byte(v), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("could not unmarshal: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func Logout(c *gin.Context) {
	if err := gothic.Logout(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("could not logout: %v", err),
		})
		return
	}

	ses := sessions.Default(c)
	ses.Delete("id")
	if err := ses.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("could not delete: %v", err),
		})
		return
	}

	// TODO: redirect
}
