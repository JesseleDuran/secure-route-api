package main

import (
	osm "github.com/JesseleDuran/gograph/osm/pbf"
	graph "github.com/JesseleDuran/secure-route-api"
	"github.com/JesseleDuran/secure-route-api/config"
	"github.com/JesseleDuran/secure-route-api/inmen"
	"github.com/JesseleDuran/secure-route-api/s3"
	"github.com/JesseleDuran/secure-route-api/service"
	"log"
	"net/http"
)

func main() {
	// Init configs
	config.Initialize()

	// Init repo
	res := s3.NewS3GraphInMemoryProvider(s3.GetClient())
	gDistanceCar, err := res.Fetch("distance", osm.Driving)
	if err != nil {
		log.Fatalf("Unable to load drive graph: %s", err.Error())
	}
	gDistanceBike, err := res.Fetch("distance", osm.Cycling)
	if err != nil {
		log.Fatalf("Unable to load bike graph: %s", err.Error())
	}
	gCrimesCar, err := res.Fetch("crimes", osm.Driving)
	if err != nil {
		log.Fatalf("Unable to load drive graph: %s", err.Error())
	}
	gCrimesBike, err := res.Fetch("crimes", osm.Cycling)
	if err != nil {
		log.Fatalf("Unable to load bike graph: %s", err.Error())
	}
	graphSource := graph.NewSource(inmen.NewGraphRepository(gDistanceCar), inmen.NewGraphRepository(gDistanceBike),
		inmen.NewGraphRepository(gCrimesCar), inmen.NewGraphRepository(gCrimesBike))

	// Raise http router.
	router := setupRoutes(graphSource, service.NewGoogleService())
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
}
