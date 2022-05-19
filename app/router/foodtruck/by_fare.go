package foodtruck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parkerduckworth/foodtruck-recommender/failure"
	"github.com/parkerduckworth/foodtruck-recommender/recommender"
)

type fareRequest struct {
	Question string `json:"question"`
	Limit    int    `json:"limit"`
}

func (r *fareRequest) validate() *failure.Error {
	if len(r.Question) == 0 {
		return failure.NewError(http.StatusBadRequest, "must provide question", nil)
	}

	if r.Limit < 1 {
		r.Limit = defaultQueryLimit
	}

	return nil
}

// ByFare is a handler func for fetching food trucks
// based on a provided question or statement indicating
// which type of food items are desired
func ByFare(c *gin.Context) {
	var request fareRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ferr := failure.NewError(http.StatusBadRequest, "invalid request body", err)

		c.AbortWithStatusJSON(ferr.StatusCode, ferr)
		return
	}

	if ferr := (&request).validate(); ferr != nil {
		c.AbortWithStatusJSON(ferr.StatusCode, ferr)
		return
	}

	data, ferr := recommender.ByFare(c, request.Question, request.Limit)
	if ferr != nil {
		c.AbortWithStatusJSON(ferr.StatusCode, ferr)
		return
	}

	c.JSON(http.StatusOK, data)
}
