package s3

import (
	"fmt"
	graph "github.com/JesseleDuran/gograph"
	osm "github.com/JesseleDuran/gograph/osm/pbf"

	api "github.com/JesseleDuran/secure-route-api"
	"github.com/JesseleDuran/secure-route-api/config"
)

// GraphInMemoryProvider implementation in S3.
type GraphInMemoryProvider struct {
	client api.S3Client
}

func NewS3GraphInMemoryProvider(client api.S3Client) api.GraphInMemoryProvider {
	return GraphInMemoryProvider{
		client: client,
	}
}

// Fetch builds a graph from a serialized file.
func (gp GraphInMemoryProvider) Fetch(content string, mode osm.Mode) (api.Graph, error) {
	path := config.Config.S3DownloadPath
	item := fmt.Sprintf("%s-%s-%s.gob", content, mode.ToString(), config.Config.Country)
	err := gp.client.Get(
		config.Config.S3BucketName,
		item,
		path+item,
	)
	if err != nil {
		return nil, fmt.Errorf("[graph_memory_provider:fetch][s3 downloading][err: %w]", err)
	}
	return graph.Deserialize(path + item), nil
}
