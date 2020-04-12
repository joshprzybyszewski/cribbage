package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshprzybyszewski/cribbage/auth"
)

var authService *auth.Provider

func (cs *cribbageServer) ginDirectToAuthProvider(c *gin.Context) {
	service := getAuthServiceFromContext(c)
	url, err := authService.GetAuthURL(service)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (cs *cribbageServer) ginGetAuthToken(c *gin.Context) {
	service := getAuthServiceFromContext(c)
	state, code := c.Query(`state`), c.Query(`code`)
	tok, err := authService.GetAccessToken(context.Background(), service, state, code)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(http.StatusOK, tok)
}

func getAuthServiceFromContext(c *gin.Context) auth.Service {
	return auth.Service(c.Param(`service`))
}
