package recommender

import (
	"encoding/json"
	"fmt"

	"github.com/parkerduckworth/foodtruck-recommender/failure"
	"github.com/semi-technologies/weaviate/entities/models"
)

type Response []Result

type Result struct {
	Name         string          `json:"name"`
	FacilityType string          `json:"facilityType"`
	Fare         string          `json:"fare"`
	Location     *ResultLocation `json:"location"`
}

type ResultLocation struct {
	Latitude   float32 `json:"latitude,omitempty"`
	Longitude  float32 `json:"longitude,omitempty"`
	MetersAway float32 `json:"metersAway,omitempty"`
	MilesAway  float32 `json:"milesAway,omitempty"`
}

type weaviateResponse struct {
	Get struct {
		FoodTruck []struct {
			FacilityType string `json:"facility_type"`
			FoodItems    string `json:"food_items"`
			Location     struct {
				Latitude  float32
				Longitude float32
			}
			Name string
		}
	}
}

func checkWeaviateResponse(resp *models.GraphQLResponse, err error) error {
	if err != nil {
		return err
	}

	if len(resp.Errors) != 0 {
		return failure.CombineGraphQLErrors(resp.Errors)
	}

	return nil
}

// because weaviate returns responses as `interface{}`, we have
// to unmarshal+marshal to access the response fields. this
// wouldn't be necessary if returning the weaviate payload
// directly to the user, but we are going to clean it up a bit
// first
func buildResponse(gql *models.GraphQLResponse) (*Response, error) {
	b, err := json.Marshal(gql.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal weaviate response")
	}

	var wResp weaviateResponse
	err = json.Unmarshal(b, &wResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal weaviate response")
	}

	resp := make(Response, len(wResp.Get.FoodTruck))
	for i := range wResp.Get.FoodTruck {
		resp[i] = Result{
			Name:         wResp.Get.FoodTruck[i].Name,
			FacilityType: wResp.Get.FoodTruck[i].FacilityType,
			Fare:         wResp.Get.FoodTruck[i].FoodItems,
			Location: &ResultLocation{
				Latitude:  wResp.Get.FoodTruck[i].Location.Latitude,
				Longitude: wResp.Get.FoodTruck[i].Location.Longitude,
			},
		}
	}

	return &resp, nil
}
