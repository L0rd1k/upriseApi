package array

// +	Get(index int) (interface{}, bool)
// +	Remove(index int)
// +	Add(elements ...interface{})
// +	Contains(elements ...interface{}) bool
// +	Swap(index_1, index_2 int)
// Insert(index int, elements ...interface{})
// Set(index int, element interface{})
// +	Empty() bool
// +	Size() int
// + 	Clear()
// +	List() []interface{}
// +	ToString() string

import (
	"errors"
	"fmt"
	"strings"
)

type List struct {
	elements []interface{}
	size     int
}

func New(elements ...interface{}) *List {
	list := &List{}
	if len(elements) > 0 {
		list.Add(elements...)
	}
	return list
}

// ============================================================

// Add element to the list
func (list *List) Add(elements ...interface{}) {
	list.growCapacity(len(elements))
	for _, elem := range elements {
		list.elements[list.size] = elem
		list.size++
	}
}

// Receive element from the list
func (list *List) Get(index int) (interface{}, bool) {
	if !list.inRange(index) { // Check if index in range of list
		return nil, false
	}
	return list.elements[index], true
}

// Remove element by index
func (list *List) Remove(index int) bool {
	if !list.inRange(index) {
		return false
	}
	list.elements[index] = nil                                    // Clean reference
	copy(list.elements[index:], list.elements[index+1:list.size]) // shift to the left by one
	list.size--
	list.shrink()
	return true
}

// Find elements in the list
func (list *List) Contains(elements ...interface{}) bool {
	for _, searchValue := range elements {
		found := false
		for _, element := range list.elements {
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
func (list *List) List() []interface{} {
	newElements := make([]interface{}, list.size, list.size)
	copy(newElements, list.elements[:list.size])
	return newElements
}

// Return index of given element
func (list *List) IndexOf(element interface{}) int {
	if list.size == 0 {
		return -1
	}
	for index, value := range list.elements {
		if value == element {
			return index
		}
	}
	return -1
}

// Check if list is empty
func (list *List) Empty() bool {
	return list.size == 0
}

// Number of elements in the list.
func (list *List) Size() int {
	return list.size
}

// Remove all elements from the list
func (list *List) Clear() {
	list.size = 0
	list.elements = []interface{}{}
}

// Exchange position of two list's elements
func (list *List) Swap(i, j int) {
	if list.inRange(i) && list.inRange(j) {
		list.elements[i], list.elements[j] = list.elements[j], list.elements[i]
	}
}

// Paste element to the given position
func (list *List) Insert(index int, elements ...interface{}) bool {
	list.growCapacity(len(elements))
	if !list.inRange(index) {
		fmt.Println(errors.New("error : out of range"))
		return false
	}
	copy(list.elements[index+len(elements):], list.elements[index:list.size])
	count := 0
	for _, elem := range elements {
		list.elements[index+count] = elem
		count++
		list.size++
	}
	return true
}

// Change value by index
func (list *List) Set(index int, value interface{}) bool {
	if !list.inRange(index) {
		return false
	}
	list.elements[index] = value
	return true
}

// ============================================================

func (list *List) inRange(index int) bool {
	return index >= 0 && index <= list.size
}

// Expand the array if necessary
func (list *List) growCapacity(n int) {
	currentCapacity := cap(list.elements)
	if list.size+n >= currentCapacity {
		newCapacity := int(2.0 * float32(currentCapacity+n))
		list.resize(newCapacity)
	}
}

// Shrink when size is 25% of capacity
func (list *List) shrink() {
	currentCapacity := cap(list.elements)
	if list.size <= int(float32(currentCapacity)*0.25) {
		list.resize(list.size)
	}
}

func (list *List) resize(capacity int) {
	newElements := make([]interface{}, capacity, capacity)
	copy(newElements, list.elements)
	list.elements = newElements
}

func (list *List) ToString() string {
	var str = ""
	items := []string{}
	for _, value := range list.elements[:list.size] {
		items = append(items, fmt.Sprintf("%v", value))
	}
	str += strings.Join(items, "\n")
	return str
}
