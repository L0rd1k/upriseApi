package utils

import (
	"fmt"
	"sort"
)

/***

// Use the generic sort.Sort and sort.Stable functions.
// They sort any collection that implements the sort.Interface interface.

type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}
***/

func Sort(values []interface{}, comparator Comparator) {
	sort.Sort(sortable{values, comparator})
}

type sortable struct {
	elements   []interface{}
	comparator Comparator
}

func (s sortable) Len() int {
	fmt.Println("Len")
	return len(s.elements)
}
func (s sortable) Swap(val_1, val_2 int) {
	fmt.Println("Swap", val_1, val_2)
	s.elements[val_1], s.elements[val_2] = s.elements[val_2], s.elements[val_1]
}
func (s sortable) Less(val_1, val_2 int) bool {
	fmt.Println("Less", val_1, val_2)
	return s.comparator(s.elements[val_1], s.elements[val_2]) < 0
}
