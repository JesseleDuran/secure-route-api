package api

import (
  "secure-route-api/actor"

  "github.com/JesseleDuran/osm-graph/graph"
  "github.com/JesseleDuran/osm-graph/graph/shortest_path/dijkstra"
  "github.com/gin-gonic/gin"
  "github.com/golang/geo/s2"
)

func New(system actor.System) *gin.Engine {
  router := gin.Default()
  group := router.Group("secure-route")
  g := graph.FromJSONGraphFileStream("downloads/osm-graph-1.json")

  group.GET("directions", func(c *gin.Context) {
    //origin := strings.Split(c.Query("origin"), ",")
    //destination := strings.Split(c.Query("destination"), ",")
    //if len(origin) < 2 || len(destination) < 2 {
    //  c.JSON(400,
    //    gin.H{"error": "origin and destination should have 2 coordinates"})
    //  return
    //}
    //o, err := coordinates.FromStrings(origin[0], origin[1])
    //d, err1 := coordinates.FromStrings(destination[0], destination[1])
    //if err != nil || err1 != nil {
    //  c.JSON(400, gin.H{"error": "invalid origin or destination values"})
    //  return
    //}
    //r, _ := system.TellSync(actor.Route{
    //  CityID:      1,
    //  Origin:      o,
    //  Destination: d,
    //  Mode:        0,
    //})
    s := s2.CellIDFromToken("94ce595164")
    e := s2.CellIDFromToken("94ce50b26c")
    _, previous := dijkstra.DijkstraFromToken(s, e, g)
    c.JSON(200, gin.H{"path": previous})
  })

  group.GET("health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
  })
  return router
}
