package main

import (
	"github.com/parkerduckworth/lonchera/app"
	"github.com/parkerduckworth/lonchera/app/config"
	"github.com/parkerduckworth/lonchera/log"
)

func init() {
	config.Setup()
	log.Setup()
}

func main() {
	app.Run()
}
