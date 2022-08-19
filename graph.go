package graph

import (
	"context"
	"fmt"
	graph "github.com/JesseleDuran/gograph"
	osm "github.com/JesseleDuran/gograph/osm/pbf"
	geojson "github.com/paulmach/go.geojson"
)

//go:generate mockery --name Graph
type Graph interface {
	DijkstraPathCoord(source, target graph.Coordinate) (float32, geojson.FeatureCollection, []uint64)
}

//go:generate mockery --name GraphInMemoryProvider
type GraphInMemoryProvider interface {
	Fetch(content string, mode osm.Mode) (Graph, error)
}

//go:generate mockery --name GraphRepository
type GraphRepository interface {
	Path(ctx context.Context, origin, destination [2]float64) (float32, geojson.FeatureCollection, []uint64)
}

//go:generate mockery --name S3Client
type S3Client interface {
	Get(bucketName, objectName, fileName string) error
	GetAllObjectKeys(bucketName string) []string
}

// ModeFromString pass a mode string to a type Mode. Default is
// cycling.
func ModeFromString(modeS string) (osm.Mode, error) {
	switch modeS {
	case "driving":
		return osm.Driving, nil
	case "bicycling":
		return osm.Cycling, nil
	case "":
		return osm.Cycling, nil
	default:
		return 0, fmt.Errorf("not valid mode")
	}
}

type Source struct {
	data []GraphRepository
}

func NewSource(repoCar, repoBike, repoCrimeCar, repoCrimeBike GraphRepository) Source {
	return Source{
		data: []GraphRepository{repoCar, repoBike, repoCrimeCar, repoCrimeBike},
	}
}

func (ms Source) GraphByModeAndContent(mode osm.Mode, content string) (GraphRepository, error) {
	if content == "distance" {
		if mode.ToString() == "drive" {
			return ms.data[0], nil
		}
		return ms.data[1], nil
	}
	if mode.ToString() == "drive" {
		return ms.data[2], nil
	}
	return ms.data[3], nil
}
