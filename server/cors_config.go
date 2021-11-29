package server

import (
	"github.com/gin-contrib/cors"
)

func getCORSConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		`https://hobbycribbage.com`,
		`http://localhost:3000`,
	}
	config.AllowMethods = []string{
		`GET`,
		`POST`,
	}
	return config
}
