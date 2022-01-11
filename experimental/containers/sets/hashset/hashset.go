package hashset

import (
	"fmt"
	"strings"
)

type Set struct {
	elements map[interface{}]struct{}
}

var itemExists = struct{}{} // New struct{} instance

func New(values ...interface{}) *Set {
	set := &Set{
		elements: make(map[interface{}]struct{}),
	}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

func (set *Set) Add(elements ...interface{}) {
	for _, elem := range elements {
		set.elements[elem] = itemExists
	}
}

func (set *Set) Remove(elements ...interface{}) {
	for _, elem := range elements {
		delete(set.elements, elem)
	}
}

func (set *Set) Contains(elements ...interface{}) bool {
	for _, elem := range elements {
		if _, contains := set.elements[elem]; !contains {
			return false
		}
	}
	return true
}

func (set *Set) Empty() bool {
	return set.Size() == 0
}

func (set *Set) Size() int {
	return len(set.elements)
}

func (set *Set) Clear() {
	set.elements = make(map[interface{}]struct{})
}

func (set *Set) List() []interface{} {
	values := make([]interface{}, set.Size())
	counter := 0
	for elem := range set.elements {
		values[counter] = elem
		counter++
	}
	return values
}

func (set *Set) ToString() string {
	var str = ""
	items := []string{}
	for itr := range set.elements {
		items = append(items, fmt.Sprintf("%v", itr))
	}
	str += strings.Join(items, "\n")
	return str
}
