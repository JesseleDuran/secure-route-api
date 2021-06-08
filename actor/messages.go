package actor

import (
  "github.com/JesseleDuran/osm-graph/coordinates"
  "github.com/JesseleDuran/osm-graph/transport"
)

// Direction
type RouteCrimes struct {
  CityID      int
  Origin      coordinates.Coordinates
  Destination coordinates.Coordinates
  Mode        transport.Mode
}

// Direction
type Route struct {
  CityID      int
  Origin      coordinates.Coordinates
  Destination coordinates.Coordinates
  Mode        transport.Mode
}

type GradientMapGeoJSON struct {
  CityID int
}
