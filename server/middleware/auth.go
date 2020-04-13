package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshprzybyszewski/cribbage/auth"
)

func Auth() gin.HandlerFunc {
	jwtSvc := auth.NewJWTService(`somethingSecret`, time.Second)
	return func(c *gin.Context) {
		tok := c.GetHeader(`x-auth-token`)
		isValid, user, err := jwtSvc.ValidateAndParseToken(tok)
		if err != nil {
			c.String(http.StatusInternalServerError, `Error: %s`, err)
			return
		}
		if !isValid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Println(`setting user to: ` + user)
		c.Set(`user`, user)
	}
}
