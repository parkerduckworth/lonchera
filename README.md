# Lonchera - SF Food Truck Recommendation API 🍔🚚

A simple REST API designed to match you with the perfect mobile refreshment for any time and place (within the City of San Francisco) :)

## About

Lonchera is a Golang service built around the [Weaviate Vector Search Engine](https://weaviate.io), using machine learning models to:

1. Recommend vendors based on free-form questions related to which food item is desired
2. Fetch a list of vendors within a specified radius

https://user-images.githubusercontent.com/31421773/169399883-d9e8c720-2500-4927-80cb-70f5a3e9a1df.mov

## Usage

### Starting/Stopping the Service

To start the service:
```
docker-compose up
```

To stop the service:
```
docker-compose stop
```

### Importing Data

> Note: The service must be started with docker-compose prior to importing!

Simply run:
```
go run cmd/import/import.go
```

### Recommend By Fare

> Note: to search, there must be data! See the [Importing Data](#importing-data) section above.

```
POST /api/v1/foodtrucks/by-fare

{
  "question": "where can i get a burger?",
  "limit": 1
}
```

Example Response:

```
[
	{
		"name": "Natan's Catering",
		"facilityType": "Truck",
		"fare": "Burgers: melts: hot dogs: burritos:sandwiches: fries: onion rings: drinks",
		"location": {
			"latitude": 37.747772,
			"longitude": -122.39703
		}
	}
]
```

### Recommend By Location

> Note: to search, there must be data! See the [Importing Data](#importing-data) section above.

```
POST /api/v1/foodtrucks/by-location

{
	"latitude": 37.798207610167076,
	"longitude": -122.43364918356474,
	"maxMilesAway": 2,
	"limit": 1
}
```

Example Response:

```
[
	{
		"name": "BOWL'D ACAI, LLC.",
		"facilityType": "Truck",
		"fare": "Acai Bowls: Poke Bowls: Smoothies: Juices",
		"location": {
			"latitude": 37.804577,
			"longitude": -122.433014,
			"metersAway": 710.5537,
			"milesAway": 0.4415168
		}
	}
]
```

### Errors

All errors are returned with the following format:
```
{
  "message": "<the problem>"
  "statusCode": <related HTTP code>
}
```

## Project Architecture

### App

Top-level application abstraction which encapsulates service, routing, and configuration.

### Env

Contains application configuration for all required environments, encoded in YAML. Absolutely no secrets should be stored in these files, or in source control in general. All secrets should be loaded into an environment, via remote manager, securely stored in a hosted platform, or simply managed locally outside of source control. Each env may have its variables injected differently, based on the environment's specific needs/usage.

### Failure

Centralized error management library which provides a common interface for handling errors of different types. This package aims to enable easy translation between language runtime errors (such as JSON failures), HTTP errors, Weaviate server/client errors, etc. 

### Log

Wrapper package for a [logrus](https://github.com/sirupsen/logrus) instance, providing structured, leveled logging throughout the application. The log level is set in the application configuration YAML files.

### Recommender

The business logic components of the application. Includes the functions responsible for recommending mobile food vendors using the various recommendation methods. Also contains the schema which is used inside the Weaviate instance, which models the FoodTruck resource.

## Roadmap

- Geocoding/reverse geocoding so that the user can determine and use the location unit that best fits their usecase.
- Configurable distance units (currently just miles for requests, miles+meters in response).
- Include data for more cities! (this is already possible, you just need the data, and to slightly refactor the import script).
