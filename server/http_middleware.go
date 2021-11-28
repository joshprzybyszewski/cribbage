package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func addCORS(
	router *gin.Engine,
) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"https://hobbycribbage.com",
		"http://localhost:3000",
	}
	config.AllowMethods = []string{
		"GET",
		"POST",
	}

	router.Use(cors.New(config))
}
