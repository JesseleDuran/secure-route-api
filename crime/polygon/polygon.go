package polygon

import (
  "secure-route-api/crime/cell"

  "github.com/golang/geo/s2"
)

// Polygon Represents a projection of coordinates to a set of points on a sphere.
// it should be noted that a point on the sphere is a vector in the
// three-dimensional plane.
type Polygon struct {
  Points []s2.Point
}

// MakeFromFlatCoords creates a polygon from a set of coordinates.
func MakeFromFlatCoords(coords [][]float64) Polygon {
  points := make([]s2.Point, 0, len(coords))
  // c -> represents a  coordinate (lat, lng)
  for _, c := range coords {
    ll := s2.LatLngFromDegrees(c[1], c[0]) // ll means latLng.
    points = append(points, s2.PointFromLatLng(ll))
  }
  return Polygon{Points: points}
}

func MakeFromPoints(points []s2.Point) Polygon {
  return Polygon{Points: points}
}

// Tessellate retrieves the cell representation of a polygon.
func (p Polygon) Tessellate(min int) cell.LinkedList {
  rc := &s2.RegionCoverer{MaxLevel: 15, MinLevel: min}
  return cell.MakeCellListFromCellUnion(rc.Covering(s2.LoopFromPoints(p.Points)))
}
