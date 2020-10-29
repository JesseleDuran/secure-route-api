package actor

import (
  "github.com/JesseleDuran/osm-graph/coordinates"
  "github.com/JesseleDuran/osm-graph/transport"
)

// Direction
type Route struct {
  CityID      int
  Origin      coordinates.Coordinates
  Destination coordinates.Coordinates
  Mode        transport.Mode
}
