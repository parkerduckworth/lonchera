package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/parkerduckworth/foodtruck-recommender/app/config"
	"github.com/parkerduckworth/foodtruck-recommender/failure"
	"github.com/parkerduckworth/foodtruck-recommender/log"
	"github.com/parkerduckworth/foodtruck-recommender/recommender/schema"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/batch"
	"github.com/semi-technologies/weaviate/entities/models"
)

func init() {
	config.Setup()
	log.Setup()
}

// Column numbers for each target field
const (
	ApplicantCol    = 1
	FacilityTypeCol = 2
	FoodItemsCol    = 11
	LatitudeCol     = 14
	LongitudeCol    = 15
)

const (
	batchSize = 10
	csvPath   = "cmd/import/Mobile_Food_Facility_Permit.csv"
)

func main() {
	recs, err := readVectorFile()
	if err != nil {
		log.Fatal(err)
	}

	client := weaviate.New(config.Conf.Weaviate)

	err = createSchema(client)
	if err != nil {
		log.Fatalf("failed to create schema: %s", failure.WeaviateError(err))
	}

	batcher := client.Batch().ObjectsBatcher()

	log.Infof("importing %d objects...\n", len(recs)-1)
	var importCount int

	for i := 1; i < len(recs); i += batchSize {
		for j := i; j < i+batchSize && j < len(recs); j++ {
			addObjectToBatch(batcher, recs[j])
			importCount++
		}

		checkBatchInsertResult(batcher.Do(context.Background()))
		log.Infof("objects imported: %d\n", importCount)
	}
}

func readVectorFile() (records [][]string, err error) {
	f, err := os.Open(csvPath)
	if err != nil {
		err = fmt.Errorf("failed to open csv: %s", err)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err = r.ReadAll()
	if err != nil {
		err = fmt.Errorf("failed to read csv: %s", err)
		return
	}

	return
}

func createSchema(client *weaviate.Client) (err error) {
	return client.Schema().
		ClassCreator().
		WithClass(schema.New()).
		Do(context.Background())
}

func addObjectToBatch(batcher *batch.ObjectsBatcher, rec []string) {
	parsedLat, err := parseFloat32(rec[LatitudeCol])
	if err != nil {
		log.Fatalf("failed to parse latitude for applicant %s: %s", rec[ApplicantCol], err)
	}

	parsedLong, err := parseFloat32(rec[LongitudeCol])
	if err != nil {
		log.Fatalf("failed to parse longitude for applicant %s: %s", rec[ApplicantCol], err)
	}

	batcher.WithObject(&models.Object{
		Class: schema.ClassName,
		Properties: map[string]interface{}{
			schema.PropName:         rec[ApplicantCol],
			schema.PropFacilityType: rec[FacilityTypeCol],
			schema.PropFoodItems:    rec[FoodItemsCol],
			schema.PropLocation: &models.GeoCoordinates{
				Latitude:  &parsedLat,
				Longitude: &parsedLong,
			},
		},
	})
}

func parseFloat32(in string) (parsed float32, err error) {
	p, err := strconv.ParseFloat(in, 32)
	if err != nil {
		err = fmt.Errorf("failed to parse float32: %s", in)
		return
	}

	parsed = float32(p)
	return
}

func checkBatchInsertResult(created []models.ObjectsGetResponse, err error) {
	if err != nil {
		log.Fatal(failure.WeaviateError(err))
	}

	// each created object can contain its own error
	// as well. iterate through and check each one
	for _, c := range created {
		if c.Result != nil {
			if c.Result.Errors != nil && c.Result.Errors.Error != nil {
				log.Fatalf("failed to create obj: %+v, with status: %v",
					c.Result.Errors.Error[0], c.Result.Status)
			}
		}
	}
}
