package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshprzybyszewski/cribbage/auth"
)

var authService *auth.AuthService

func (cs *cribbageServer) ginGetAuthURL(c *gin.Context) {
	url := authService.GetAuthURL()
	c.String(http.StatusOK, url)
}

func (cs *cribbageServer) ginGetAuthToken(c *gin.Context) {
	state, code := c.Query(`state`), c.Query(`code`)
	fmt.Println(`hit callback!`)
	tok, err := authService.GetAccessToken(context.Background(), state, code)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(http.StatusOK, tok)
}
