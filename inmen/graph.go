package inmen

import (
	"context"
	graph "github.com/JesseleDuran/gograph"

	api "github.com/JesseleDuran/secure-route-api"

	geojson "github.com/paulmach/go.geojson"
)

type GraphRepository struct {
	Repo api.Graph
}

func NewGraphRepository(repo api.Graph) api.GraphRepository {
	return GraphRepository{
		Repo: repo,
	}
}

func (r GraphRepository) Path(ctx context.Context, origin, destination [2]float64) (float32, geojson.FeatureCollection, []uint64) {
	return r.Repo.DijkstraPathCoord(graph.Coordinate{
		Lat: origin[0],
		Lng: origin[1],
	}, graph.Coordinate{
		Lat: destination[0],
		Lng: destination[1],
	})
}
