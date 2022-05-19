package schema

import "github.com/semi-technologies/weaviate/entities/models"

const (
	ClassName = "FoodTruck"

	PropName                = "name"
	PropFacilityType        = "facility_type"
	PropFoodItems           = "food_items"
	PropLocation            = "location"
	PropLocationLatitude    = "latitude"
	PropLocationLongitude   = "longitude"
	PropAdditional          = "_additional"
	PropAdditionalCertainty = "certainty"
)

// New returns a Foodtruck Class instance
func New() *models.Class {
	return &models.Class{
		Class:       ClassName,
		Description: "A mobile sustenance dispenser",
		Properties: []*models.Property{
			{
				DataType: []string{"string"},
				Name:     PropName,
			},
			{
				DataType: []string{"string"},
				Name:     PropFacilityType,
			},
			{
				DataType: []string{"text"},
				Name:     PropFoodItems,
			},
			{
				DataType: []string{"geoCoordinates"},
				Name:     PropLocation,
			},
		},
	}
}
