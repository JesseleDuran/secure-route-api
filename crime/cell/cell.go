package cell

import (
  "math"

  "github.com/golang/geo/s2"
)

type Cell struct {
  ID    uint64
  Depth uint8
}

type Cells []Cell

// Make creates a new cell from an uint64.
// The uint 64 represents a point on the hilbert curve.
func Make(id uint64) Cell {
  return Cell{ID: id}
}

// Children retrieves the four children of a cell.
//
// For instance:
//
// |          |                |  a | b  |
// |   Cell   |.Children() ->  |----|----|
// |          |                |  c | d  |
//
func (c Cell) Children() []Cell {
  children := s2.CellID(c.ID).Children()
  return []Cell{
    {ID: uint64(children[0]), Depth: c.Depth + 1},
    {ID: uint64(children[1]), Depth: c.Depth + 1},
    {ID: uint64(children[2]), Depth: c.Depth + 1},
    {ID: uint64(children[3]), Depth: c.Depth + 1},
  }
}

// ContainsPoint Determine if the given cell contains a point.
func (c Cell) ContainsPoint(p s2.Point) bool {
  return s2.CellID(c.ID).Contains(s2.CellFromPoint(p).ID())
}

// Token retrieve the token of the given cell.
// The token is a hexadecimal string.
func (c Cell) Token() string {
  return s2.CellID(c.ID).ToToken()
}

// Level retrieve the level of the given cell.
func (c Cell) Level() int {
  return s2.CellID(c.ID).Level()
}

// LatLng Retrieve the lat lng centroid of the given cell
func (c Cell) LatLng() [2]float64 {
  ll := s2.CellID(c.ID).LatLng()
  return [2]float64{ll.Lat.Degrees(), ll.Lng.Degrees()}
}

// calculate the circle radius that wrap the given cell.
func (c Cell) CapRadius() float64 {
  circle := s2.CellFromCellID(s2.CellID(c.ID)).CapBound()
  return 6378000 * circle.Radius().Radians()
}

func (c Cell) IsEmpty() bool {
  return c.ID == math.MaxInt64
}

func (c Cell) IntersectsCap(cap s2.Cap) bool {
  return cap.IntersectsCell(s2.CellFromCellID(s2.CellID(c.ID)))
}
