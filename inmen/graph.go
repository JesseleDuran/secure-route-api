package inmen

import (
	"context"
	graph "github.com/JesseleDuran/gograph"
	"github.com/JesseleDuran/gograph/nearest_edge"

	api "github.com/JesseleDuran/secure-route-api"

	geojson "github.com/paulmach/go.geojson"
)

type GraphRepository struct {
	Repo api.Graph
	Node nearest_edge.Node
}

func NewGraphRepository(repo api.Graph) api.GraphRepository {
	return GraphRepository{
		Repo: repo,
		Node: repo.BuildEdgeIndex(),
	}
}

func (r GraphRepository) Index() nearest_edge.Node {
	return r.Node
}

func (r GraphRepository) Path(ctx context.Context, origin, destination [2]float64, node nearest_edge.Node) (float32, geojson.FeatureCollection, []uint64) {
	return r.Repo.DijkstraPathCoord(graph.Coordinate{
		Lat: origin[0],
		Lng: origin[1],
	}, graph.Coordinate{
		Lat: destination[0],
		Lng: destination[1],
	}, node)
}
