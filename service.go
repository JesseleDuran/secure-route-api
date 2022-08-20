package graph

import (
	geojson "github.com/paulmach/go.geojson"
)

type RoutesService interface {
	Route(origin, destination [2]float64, mode string) (float32, geojson.FeatureCollection)
}
