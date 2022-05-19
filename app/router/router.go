// Package router provides the server routes, delegates all the
// server routes to their appropriate handler funcs, and contains
// code to ensure dependencies are available at server starup.
package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parkerduckworth/foodtruck-recommender/app/router/foodtruck"
	"github.com/pkg/errors"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
)

// PingWeaviate checks for a running Weaviate instance at the host
// address provided in the config file. If an instance is running,
// nil is returned. Otherwise, the executing thread is panicked.
func PingWeaviate(config weaviate.Config) error {
	hostURL := fmt.Sprintf("%s://%s", config.Scheme, config.Host)

	_, err := http.Get(hostURL)
	if err != nil {
		err = errors.Wrapf(err, "failed to connect to weaviate host %s", hostURL)
		log.Fatal(err)
	}

	return nil
}

// SetupRoutes takes sets of routes, handler funcs, and middleware, and
// attaches them to the server router to be served during app lifetime
func SetupRoutes(r *gin.Engine) {
	apiRoutes := r.Group("/api")
	setupV1Routes(apiRoutes)
}

func setupV1Routes(r *gin.RouterGroup) {
	v1Routes := r.Group("/v1")
	{
		foodtruckRoutes := v1Routes.Group("/foodtrucks")
		{
			foodtruckRoutes.POST("/by-fare", foodtruck.ByFare)
			foodtruckRoutes.POST("/by-location", foodtruck.ByLocation)
		}
	}
}
