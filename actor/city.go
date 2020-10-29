package actor

import (
  "log"
  "secure-route-api/crime"
  "time"

  "github.com/AsynkronIT/protoactor-go/actor"
  "github.com/JesseleDuran/osm-graph/graph"
  "github.com/JesseleDuran/osm-graph/graph/shortest_path/dijkstra"
  "github.com/golang/geo/s2"
)

type City struct {
  PID    *actor.PID
  Graph  graph.Graph
  Crimes *crime.Crimes
}

func (c *City) TellSync(msg interface{}) interface{} {
  context := actor.EmptyRootContext
  future := context.RequestFuture(c.PID, msg, time.Second)
  r, _ := future.Result()
  return r
}

func (c *City) Receive(context actor.Context) {
  switch msg := context.Message().(type) {
  case *actor.Started:
    c.Graph = graph.FromJSONGraphFileStream("downloads/osm-graph-1.json")
    log.Println("done")

  case Route:
    log.Println("ADA", msg.Origin)
    s := s2.CellIDFromToken("94ce595164")
    e := s2.CellIDFromToken("94ce50b26c")
    _, previous := dijkstra.DijkstraFromToken(s, e, c.Graph)
    context.Respond(previous)
  }
}
