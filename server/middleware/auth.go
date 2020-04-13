package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshprzybyszewski/cribbage/auth"
)

func Auth() gin.HandlerFunc {
	jwtSvc := auth.NewJWTService(`somethingSecret`, time.Second)
	return func(c *gin.Context) {
		tok, err := c.Cookie(`token`)
		if err != nil {
			c.String(http.StatusUnauthorized, `Unauthorized!`)
			c.Abort()
			return
		}
		isValid, user, err := jwtSvc.ValidateAndParseToken(tok)
		if err != nil {
			c.String(http.StatusInternalServerError, `Error: %s`, err)
			c.Abort()
			return
		}
		if !isValid {
			c.String(http.StatusUnauthorized, `Unauthorized!`)
			c.Abort()
			return
		}
		c.Set(`user`, user)
	}
}
