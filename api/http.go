package api

import (
  "log"
  "secure-route-api/actor"
  "strings"

  "github.com/JesseleDuran/osm-graph/coordinates"
  "github.com/gin-gonic/gin"
)

func New(system actor.System) *gin.Engine {
  router := gin.Default()
  group := router.Group("secure-route")

  group.GET("directions", func(c *gin.Context) {
    origin := strings.Split(c.Query("origin"), ",")
    destination := strings.Split(c.Query("destination"), ",")
    if len(origin) < 2 || len(destination) < 2 {
      c.JSON(400, gin.H{"error": "origin and destination should have 2 coordinates"})
      return
    }
    o, err := coordinates.FromStrings(origin[0], origin[1])
    d, err1 := coordinates.FromStrings(destination[0], destination[1])
    if err != nil || err1 != nil {
      log.Println(err, err1)
      c.JSON(400, gin.H{"error": "invalid origin or destination values"})
      return
    }
    r, _ := system.TellSync(actor.Route{
      CityID:      1,
      Origin:      o,
      Destination: d,
      Mode:        0,
    })
    c.JSON(200, gin.H{"routes": []interface{}{r}})
  })

  group.GET("health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
  })
  return router
}
