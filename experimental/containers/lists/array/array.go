package array

import (
	"fmt"
	"strings"

	"github.com/L0rd1k/uprise-api/experimental/containers/lists"
	"github.com/L0rd1k/uprise-api/experimental/containers/utils"
)

type Array struct {
	elements []interface{}
	size     int
}

var _ lists.List = (*Array)(nil)

func New(elements ...interface{}) *Array {
	array := &Array{}
	if len(elements) > 0 {
		array.Add(elements...)
	}
	return array
}

// Add element to the list
func (array *Array) Add(elements ...interface{}) {
	array.growCapacity(len(elements))
	for _, elem := range elements {
		array.elements[array.size] = elem
		array.size++
	}
}

// Receive element from the list
func (array *Array) Get(index int) (interface{}, bool) {
	if !array.inRange(index) { // Check if index in range of list
		return nil, false
	}
	return array.elements[index], true
}

// Remove element by index
func (array *Array) Remove(index int) bool {
	if !array.inRange(index) {
		return false
	}
	array.elements[index] = nil                                      // Clean reference
	copy(array.elements[index:], array.elements[index+1:array.size]) // shift to the left by one
	array.size--
	array.shrink()
	return true
}

// Find elements in the list
func (array *Array) Contains(elements ...interface{}) bool {
	for _, searchValue := range elements {
		found := false
		for _, element := range array.elements {
			if element == searchValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Return all elements from the list
func (array *Array) List() []interface{} {
	newElements := make([]interface{}, array.size, array.size)
	copy(newElements, array.elements[:array.size])
	return newElements
}

// Return index of given element
func (array *Array) IndexOf(element interface{}) int {
	if array.size == 0 {
		return -1
	}
	for index, value := range array.elements {
		if value == element {
			return index
		}
	}
	return -1
}

// Check if list is empty
func (array *Array) Empty() bool {
	return array.size == 0
}

// Number of elements in the list.
func (array *Array) Size() int {
	return array.size
}

// Remove all elements from the list
func (array *Array) Clear() {
	array.size = 0
	array.elements = []interface{}{}
}

// Exchange position of two list's elements
func (array *Array) Swap(i, j int) {
	if array.inRange(i) && array.inRange(j) {
		array.elements[i], array.elements[j] = array.elements[j], array.elements[i]
	}
}

// Paste element to the given position
func (array *Array) Insert(index int, elements ...interface{}) bool {
	array.growCapacity(len(elements))
	if !array.inRange(index) {
		return false
	}
	copy(array.elements[index+len(elements):], array.elements[index:array.size])
	count := 0
	for _, elem := range elements {
		array.elements[index+count] = elem
		count++
		array.size++
	}
	return true
}

// Change value by index
func (array *Array) Set(index int, value interface{}) bool {
	if !array.inRange(index) {
		return false
	}
	array.elements[index] = value
	return true
}

func (array *Array) Sort(comparator utils.Comparator) {
	if len(array.elements) < 2 {
		return
	}
	utils.Sort(array.elements[:array.size], comparator)
}

func (array *Array) inRange(index int) bool {
	return index >= 0 && index <= array.size
}

// Expand the array if necessary
func (array *Array) growCapacity(n int) {
	currentCapacity := cap(array.elements)
	if array.size+n >= currentCapacity {
		newCapacity := int(2.0 * float32(currentCapacity+n))
		array.resize(newCapacity)
	}
}

// Shrink when size is 25% of capacity
func (array *Array) shrink() {
	currentCapacity := cap(array.elements)
	if array.size <= int(float32(currentCapacity)*0.25) {
		array.resize(array.size)
	}
}

func (array *Array) resize(capacity int) {
	newElements := make([]interface{}, capacity, capacity)
	copy(newElements, array.elements)
	array.elements = newElements
}

func (array *Array) ToString() string {
	var str = ""
	items := []string{}
	for _, value := range array.elements[:array.size] {
		items = append(items, fmt.Sprintf("%v", value))
	}
	str += strings.Join(items, "\n")
	return str
}
