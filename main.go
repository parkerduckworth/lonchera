package main

import (
	"github.com/parkerduckworth/foodtruck-recommender/app"
	"github.com/parkerduckworth/foodtruck-recommender/app/config"
	"github.com/parkerduckworth/foodtruck-recommender/log"
)

func init() {
	config.Setup()
	log.Setup()
}

func main() {
	app.Run()
}
