package cell

import (
	"math"

	"github.com/golang/geo/s2"
)

type LinkedList struct {
	head, tail *Node
	size       int
}

type Node struct {
	Cell
	next *Node
}

// MakeCellList creates a new linked list of cells.
func MakeCellList(cells ...Cell) LinkedList {
	var list LinkedList
	list.Add(cells...)
	return list
}

// MakeCellListFromCellUnion creates a new linked list from a set of cell ids.
func MakeCellListFromCellUnion(union s2.CellUnion) LinkedList {
	var list LinkedList
	for _, c := range union {
		list.Add(Make(uint64(c)))
	}
	return list
}

// Add push new elements to the cell list.
func (l *LinkedList) Add(cells ...Cell) {
	for _, c := range cells {
		if l.head == nil {
			l.head = &Node{Cell: c}
			l.tail = l.head
		} else {
			n := &Node{Cell: c}
			l.tail.next, l.tail = n, n
		}
		l.size++
	}
}

// Pop "Remove" the first element of the list.
func (l *LinkedList) Pop() Cell {
	head := l.head
	if head != nil {
		l.head = head.next
		l.size--
		return head.Cell
	}
	return Cell{ID: math.MaxInt64}
}

// retrieves the number of nodes present on the list
func (l LinkedList) Size() int {
	return l.size
}
