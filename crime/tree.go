package crime

import (
  "log"
  "secure-route-api/crime/cell"
  "secure-route-api/crime/polygon"

  "github.com/golang/geo/s1"
  "github.com/golang/geo/s2"
)

type Tree struct {
  root *node
  Len  int
}

type node struct {
  cid      cell.Cell
  children []*node
  Crimes   []Crime
}

func MakeTree(crimes []Crime) Tree {
  ch := s2.NewConvexHullQuery()
  for _, c := range crimes {
    ch.AddPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(c.Lat, c.Lng)))
  }
  cells := polygon.MakeFromPoints(ch.ConvexHull().Vertices()).Tessellate(5)
  nodes := make([]*node, 0, cells.Size())
  for {
    c := cells.Pop()
    if c.IsEmpty() {
      break
    }
    nodes = append(nodes, &node{
      cid: c, children: make([]*node, 0), Crimes: make([]Crime, 0),
    })
  }
  t := Tree{root: &node{children: nodes}}
  for _, c := range crimes {
    t.Insert(c)
  }
  log.Printf("Tree of %d crimes built", len(crimes))
  return t
}

func (t *Tree) Insert(c Crime) bool {
  for _, n := range t.root.children {
    if n.insert(c) {
      t.Len += 1
      return true
    }
  }
  return false
}

func (n *node) insert(c Crime) bool {
  p := s2.PointFromLatLng(s2.LatLngFromDegrees(c.Lat, c.Lng))

  if !n.cid.ContainsPoint(p) {
    return false
  }

  if len(n.Crimes) < 1 || n.cid.Level() == 25 {
    n.Crimes = append(n.Crimes, c)
    return true
  } else if len(n.children) == 0 {
    for _, c := range n.cid.Children() {
      n.children = append(n.children, &node{
        cid: c, children: make([]*node, 0), Crimes: make([]Crime, 0),
      })
    }
  }

  a, b, e, f := n.children[0], n.children[1], n.children[2], n.children[3]
  return a.insert(c) || b.insert(c) || e.insert(c) || f.insert(c)
}

func (t *Tree) Search(lat, lng, radius float64) []Crime {
  p := s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng))
  circle := s2.CapFromCenterAngle(p, s1.Angle(radius/float64(6378000)))
  result := make([]Crime, 0)
  for _, n := range t.root.children {
    if n.cid.IntersectsCap(circle) {
      result = append(result, n.search(circle)...)
    }
  }
  return result
}

func (n *node) search(cap s2.Cap) []Crime {
  result := make([]Crime, 0)
  if !n.cid.IntersectsCap(cap) {
    return result
  }
  for _, u := range n.Crimes {
    p := s2.PointFromLatLng(s2.LatLngFromDegrees(u.Lat, u.Lng))
    if cap.ContainsPoint(p) {
      result = append(result, u)
    }
  }
  for _, child := range n.children {
    if child.cid.IntersectsCap(cap) {
      result = append(result, child.search(cap)...)
    }
  }
  return result
}

func (n *node) replace(cap s2.Cap, crime Crime) bool {
  if !n.cid.IntersectsCap(cap) {
    return false
  }
  for idx, u := range n.Crimes {
    if u.ID == crime.ID {
      n.Crimes[idx] = crime
      return true
    }
  }
  for _, child := range n.children {
    if child.cid.IntersectsCap(cap) && child.replace(cap, crime) {
      return true
    }
  }
  return false
}

func (t *Tree) Replace(lat, lng, radius float64, crime Crime) bool {
  p := s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng))
  circle := s2.CapFromCenterAngle(p, s1.Angle(radius/float64(6378000)))

  for _, n := range t.root.children {
    if n.cid.IntersectsCap(circle) && n.replace(circle, crime) {
      return true
    }
  }
  return false
}
