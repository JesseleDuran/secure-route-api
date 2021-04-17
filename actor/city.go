package actor

import (
  "log"
  "secure-route-api/crime"
  "time"

  "github.com/AsynkronIT/protoactor-go/actor"
  "github.com/JesseleDuran/osm-graph/graph"
  "github.com/JesseleDuran/osm-graph/graph/shortest_path/dijkstra"
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
    log.Println("building graph")
    c.Graph = graph.BuildFromJsonFile("downloads/osm-graph-medellin-name-17." +
      "json", c.Crimes.SetWeight)
    log.Println("done graph")

  case Route:
    d := dijkstra.Dijkstra{c.Graph}
    result := d.FromCoordinates(msg.Origin, msg.Destination)
    context.Respond(result)
  }
}
