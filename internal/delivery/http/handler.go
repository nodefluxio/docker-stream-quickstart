package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// New returns http handler
func New() *gin.Engine {
	router := gin.Default()
	// will use for futher development
	// router.Use(cors.Default())
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"*"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	// 	return origin == "https://github.com"
	// 	// },
	// 	MaxAge: 12 * time.Hour,
	// }))
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
