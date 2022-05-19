package recommender

import (
	"context"
	"math"
	"net/http"

	"github.com/parkerduckworth/foodtruck-recommender/app/config"
	"github.com/parkerduckworth/foodtruck-recommender/failure"
	"github.com/parkerduckworth/foodtruck-recommender/recommender/schema"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/filters"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/graphql"
)

const (
	ErrFailedToRecommendByLocation = "failed to recommend by location"
)

type GeoCoordinates struct {
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
	MaxDistance float32 `json:"maxMetersAway"`
}

func ByLocation(ctx context.Context, coord *filters.GeoCoordinatesParameter, limit int) (interface{}, *failure.Error) {
	client := weaviate.New(config.Conf.Weaviate)

	fields := []graphql.Field{
		{Name: schema.PropName},
		{Name: schema.PropFacilityType},
		{Name: schema.PropFoodItems},
		{Name: schema.PropLocation, Fields: []graphql.Field{
			{Name: schema.PropLocationLatitude},
			{Name: schema.PropLocationLongitude},
		}},
	}

	where := filters.Where().
		WithOperator(filters.WithinGeoRange).
		WithPath([]string{schema.PropLocation}).
		WithValueGeoRange(coord)

	result, err := client.GraphQL().Get().
		WithClassName(schema.ClassName).
		WithFields(fields...).
		WithWhere(where).
		WithLimit(limit).
		Do(ctx)

	err = checkWeaviateResponse(result, err)
	if err != nil {
		return nil, failure.NewError(
			http.StatusInternalServerError, ErrFailedToRecommendByLocation, err)
	}

	resp, err := buildResponse(result)
	if err != nil {
		return nil, failure.NewError(
			http.StatusInternalServerError, ErrFailedToRecommendByLocation, err)
	}

	insertDistances(coord, resp)
	return resp, nil
}

type geoPoints struct {
	src geoPoint
	dst geoPoint
}

type geoPoint struct {
	lat float32
	lng float32
}

func insertDistances(coord *filters.GeoCoordinatesParameter, resp *Response) {
	for _, res := range *resp {
		if res.Location.Latitude != 0 && res.Location.Longitude != 0 {
			metersAway, milesAway := calculateGeoDistance(geoPoints{
				src: geoPoint{lat: coord.Latitude, lng: coord.Longitude},
				dst: geoPoint{res.Location.Latitude, res.Location.Longitude},
			})

			res.Location.MetersAway = metersAway
			res.Location.MilesAway = milesAway
		}
	}
}

// calculateGeoDistance uses the haversine distance formula
// to calculate the number of meters/miles each result is
// from the input location
func calculateGeoDistance(pts geoPoints) (metersAway, milesAway float32) {
	radiusMeters := 6371e3

	φ1 := float64(pts.src.lat * math.Pi / 180) // φ, λ in radians
	φ2 := float64(pts.dst.lat * math.Pi / 180)

	Δφ := float64((pts.dst.lat - pts.src.lat) * math.Pi / 180)
	Δλ := float64((pts.dst.lng - pts.src.lng) * math.Pi / 180)

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distMeters := radiusMeters * c
	distMiles := distMeters * 0.00062137

	return float32(distMeters), float32(distMiles)
}
