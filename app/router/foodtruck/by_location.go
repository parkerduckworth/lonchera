package foodtruck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parkerduckworth/lonchera/failure"
	"github.com/parkerduckworth/lonchera/recommender"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/filters"
)

type locationRequest struct {
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
	MaxMilesAway float32 `json:"maxMilesAway"`
	Limit        int     `json:"limit"`
}

func (r *locationRequest) validate() *failure.Error {
	if r.MaxMilesAway == 0 {
		return failure.NewError(http.StatusBadRequest, "must provide maxMilesAway", nil)
	}

	if r.Limit < 1 {
		r.Limit = defaultQueryLimit
	}

	return nil
}

// ByLocation is a handler func for fetching food trucks
// near a given set of geo coordinates
func ByLocation(c *gin.Context) {
	var request locationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ferr := failure.NewError(http.StatusBadRequest, "invalid request body", err)

		c.AbortWithStatusJSON(ferr.StatusCode, ferr)
		return
	}

	if ferr := (&request).validate(); ferr != nil {
		c.AbortWithStatusJSON(ferr.StatusCode, ferr)
		return
	}

	data, ferr := recommender.ByLocation(c, &filters.GeoCoordinatesParameter{
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		MaxDistance: milesToMeters(request.MaxMilesAway),
	}, request.Limit)

	if ferr != nil {
		c.AbortWithStatusJSON(ferr.StatusCode, ferr)
		return
	}

	c.JSON(http.StatusOK, data)
}

func milesToMeters(m float32) float32 {
	return m * 1609.344
}
