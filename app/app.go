// Package app provides top-level application abstraction
// which encapsulates service, routing, and configuration
package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/parkerduckworth/lonchera/app/config"
	"github.com/parkerduckworth/lonchera/app/router"
)

// Run sets up all application dependencies
// and starts the server
func Run() {
	r := gin.New()
	gin.SetMode(toGinMode(config.Conf.Env))

	r.Use(gin.Logger())
	r.SetTrustedProxies(nil)
	router.SetupRoutes(r)

	err := router.PingWeaviate(config.Conf.Weaviate)
	if err != nil {
		log.Fatal(err)
	}

	r.Run(":" + config.Conf.Server.HTTPPort)
}

func toGinMode(env string) string {
	switch env {
	case "dev":
		return gin.DebugMode
	case "staging":
		return gin.ReleaseMode
	case "prod":
		return gin.ReleaseMode
	default:
		return gin.DebugMode
	}
}
