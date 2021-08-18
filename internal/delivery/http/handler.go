package httphandler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// New returns http handler
func New() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(configureMiddleware()))
	router.Use(gin.Logger())

	return router
}

func configureMiddleware() cors.Config {
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origins", "Content-Type", "Authorization"}
	config.AllowAllOrigins = true

	return config
}
