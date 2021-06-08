package actor

import (
  "log"
  "secure-route-api/crime"
  "time"

  "github.com/AsynkronIT/protoactor-go/actor"
  "github.com/JesseleDuran/osm-graph/coordinates"
  "github.com/JesseleDuran/osm-graph/graph"
  "github.com/JesseleDuran/osm-graph/graph/shortest_path"
  "github.com/JesseleDuran/osm-graph/graph/shortest_path/dijkstra"
  "github.com/JesseleDuran/osm-graph/json"
  "github.com/golang/geo/s2"
  geojson "github.com/paulmach/go.geojson"
)

type City struct {
  PID           *actor.PID
  GraphCrimes   graph.Graph
  GraphDistance graph.Graph
  Crimes        *crime.Crimes
}

type DirectionResponse struct {
  DurationWalking float64             `json:"duration_walking"`
  DurationDriving float64             `json:"duration_driving"`
  TotalCrimes     float64             `json:"total_crimes"`
  Steps           shortest_path.Steps `json:"steps"`
  Polyline        [][2]float64        `json:"polyline"`
  Distance        float64             `json:"distance"`
}

const AvgSpeedCar = 833 //meters/min
const AvgSpeedHuman = 100 //meters/min

func (c *City) TellSync(msg interface{}) interface{} {
  context := actor.EmptyRootContext
  future := context.RequestFuture(c.PID, msg, time.Second)
  r, _ := future.Result()
  return r
}

func (c *City) Tell(msg interface{}) {
  context := actor.EmptyRootContext
  context.Send(c.PID, msg)
}

func (c *City) Receive(context actor.Context) {
  switch msg := context.Message().(type) {
  case *actor.Started:
    log.Println("building graph")
    c.GraphCrimes = graph.BuildFromJsonFile(
      "downloads/osm-graph-medellin-name-17.json",
      c.Crimes.SetWeight,
    )
    c.GraphDistance = graph.BuildFromJsonFile(
      "downloads/osm-graph-medellin-name-17.json",
      nil,
    )
    log.Println("done graph")

  case RouteCrimes:
    d := dijkstra.Dijkstra{Graph: c.GraphCrimes}
    result := d.FromCoordinates(msg.Origin, msg.Destination)
    distance := DistanceFromSteps(result.Steps)
    context.Respond(DirectionResponse{
      DurationWalking: timeFormula(AvgSpeedHuman, result.TotalWeight),
      DurationDriving: timeFormula(AvgSpeedCar, result.TotalWeight),
      TotalCrimes:     result.TotalWeight,
      Steps:           result.Steps,
      Polyline:        result.Polyline,
      Distance:        distance,
    })

  case Route:
    d := dijkstra.Dijkstra{Graph: c.GraphDistance}
    result := d.FromCoordinates(msg.Origin, msg.Destination)
    context.Respond(DirectionResponse{
      DurationWalking: timeFormula(AvgSpeedHuman, result.TotalWeight),
      DurationDriving: timeFormula(AvgSpeedCar, result.TotalWeight),
      TotalCrimes:     CrimesFromSteps(result.Steps, *c.Crimes),
      Steps:           result.Steps,
      Polyline:        result.Polyline,
      Distance:        result.TotalWeight,
    })

  case GradientMapGeoJSON:
    result := geojson.NewFeatureCollection()
    for id, _ := range c.GraphCrimes.Nodes {
      crimes := c.Crimes.FindByRadius(coordinates.Coordinates{
        Lat: id.LatLng().Lat.Degrees(),
        Lng: id.LatLng().Lng.Degrees(),
      }, 100)
      p := geojson.NewPolygonFeature(SimplePolygonFromCell(id))
      p.SetProperty("crimes", len(crimes))
      p.SetProperty("color", ColorByNCrimes(len(crimes)))
      p.SetProperty("ID", id.String())
      result.AddFeature(p)
    }
    json.Write("nodes.json", result)
    log.Println("done geojson")
  }
}

func ColorByNCrimes(n int) string {
  if n >= 1000 {
    return "#f75257"
  } else if n > 900 {
    return "#fb7054"
  } else if n > 800 {
    return "#ff7d66"
  } else if n > 700 {
    return "#ff9365"
  } else if n > 600 {
    return "#ffa969"
  } else if n > 500 {
    return "#ffbd72"
  } else if n > 400 {
    return "#f5c670"
  } else if n > 300 {
    return "#e9ce71"
  } else if n > 200 {
    return "#dcd675"
  } else if n > 100 {
    return "#c0d66e"
  } else if n > 0 {
    return "#a1d56d"
  } else {
    return "#43d27d"
  }
}

func SimplePolygonFromCell(c s2.CellID) [][][]float64 {
  helperCell := s2.CellFromCellID(c)
  lowerLeft := s2.LatLngFromPoint(helperCell.Vertex(0))
  lowerRight := s2.LatLngFromPoint(helperCell.Vertex(1))
  upperRight := s2.LatLngFromPoint(helperCell.Vertex(2))
  upperLeft := s2.LatLngFromPoint(helperCell.Vertex(3))
  return [][][]float64{
    {
      {upperLeft.Lng.Degrees(), upperLeft.Lat.Degrees()},
      {upperRight.Lng.Degrees(), upperRight.Lat.Degrees()},
      {lowerRight.Lng.Degrees(), lowerRight.Lat.Degrees()},
      {lowerLeft.Lng.Degrees(), lowerLeft.Lat.Degrees()},
      {upperLeft.Lng.Degrees(), upperLeft.Lat.Degrees()},
    },
  }
}

func timeFormula(speed, distance float64) float64 {
  return distance/speed
}

func DistanceFromSteps(steps shortest_path.Steps) float64 {
  result := 0.0
  for i := 0; i < len(steps)-1; i++ {
    start := steps[i].StartLocation
    end := steps[i+1].EndLocation
    result += coordinates.Distance(start, end)
  }
  return result
}

func CrimesFromSteps(steps shortest_path.Steps, crimes crime.Crimes) float64 {
  result := 0.0
  for i := 0; i < len(steps)-1; i++ {
    start := steps[i].StartLocation
    end := steps[i+1].EndLocation
    result += crimes.SetWeight(start, end)
  }
  return result
}
