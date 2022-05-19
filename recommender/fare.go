package recommender

import (
	"context"
	"net/http"

	"github.com/parkerduckworth/foodtruck-recommender/app/config"
	"github.com/parkerduckworth/foodtruck-recommender/failure"
	"github.com/parkerduckworth/foodtruck-recommender/recommender/schema"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/graphql"
)

const (
	ErrFailedToRecommendByFare = "failed to recommend by fare"
)

func ByFare(ctx context.Context, question string, limit int) (*Response, *failure.Error) {
	client := weaviate.New(config.Conf.Weaviate)

	fields := []graphql.Field{
		{Name: schema.PropName},
		{Name: schema.PropFacilityType},
		{Name: schema.PropFoodItems},
		{Name: schema.PropLocation, Fields: []graphql.Field{
			{Name: schema.PropLocationLatitude},
			{Name: schema.PropLocationLongitude},
		}},
		{Name: schema.PropAdditional, Fields: []graphql.Field{
			{Name: schema.PropAdditionalCertainty},
		}},
	}

	ask := client.GraphQL().AskArgBuilder().
		WithQuestion(question).
		WithCertainty(0.6)

	result, err := client.GraphQL().Get().
		WithClassName(schema.ClassName).
		WithFields(fields...).
		WithAsk(ask).
		WithLimit(limit).
		Do(ctx)

	err = checkWeaviateResponse(result, err)
	if err != nil {
		return nil, failure.NewError(
			http.StatusInternalServerError, ErrFailedToRecommendByFare, err)
	}

	resp, err := buildResponse(result)
	if err != nil {
		return nil, failure.NewError(
			http.StatusInternalServerError, ErrFailedToRecommendByFare, err)
	}

	return resp, nil
}
