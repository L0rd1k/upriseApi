package singlell

import (
	"fmt"
	"strings"

	"github.com/L0rd1k/uprise-api/experimental/containers/utils"
)

type Node struct {
	value interface{}
	next  *Node
}

type SingleList struct {
	head *Node
	tail *Node
	size int
}

//=============================================================

func New(elements ...interface{}) *SingleList {
	list := &SingleList{}
	if len(elements) > 0 {
		list.Add(elements...)
	}
	return list
}

//=============================================================

func (list *SingleList) Add(elements ...interface{}) {
	for _, value := range elements {
		newNode := &Node{value: value}
		if list.size == 0 {
			list.head = newNode
		} else {
			list.tail.next = newNode
		}
		list.tail = newNode
		list.size++
	}
}

func (list *SingleList) Remove(index int) bool {
	if !list.inRange(index) {
		return false
	}

	if list.size == 1 {
		list.Clear()
		return true
	}

	var prevNode *Node // Create prev element
	currNode := list.head
	for i := 0; i != index; i, currNode = i+1, currNode.next {
		prevNode = currNode
	}

	if currNode == list.head {
		list.head = currNode.next
	} else if currNode == list.tail {
		list.tail = prevNode
	}

	if prevNode != nil {
		prevNode.next = currNode.next
	}
	currNode = nil
	list.size--
	return true
}

func (list *SingleList) ToString() string {
	var str = ""
	items := []string{}
	for node := list.head; node != nil; node = node.next {
		items = append(items, fmt.Sprintf("%v", node.value))
	}
	str += strings.Join(items, "\n")
	return str
}

func (list *SingleList) Clear() {
	list.size = 0
	list.head = nil
	list.tail = nil
}

func (list *SingleList) Empty() bool {
	return list.size == 0
}

func (list *SingleList) Size() int {
	return list.size
}

func (list *SingleList) List() []interface{} {
	elements := make([]interface{}, list.size, list.size)
	for i, element := 0, list.head; element != nil; i, element = i+1, element.next {
		elements[i] = element.value
	}
	return elements
}

func (list *SingleList) Get(index int) (interface{}, bool) {
	if !list.inRange(index) {
		return nil, false
	}
	element := list.head
	for i := 0; i != index; i, element = i+1, element.next {
	}
	return element.value, true
}

func (list *SingleList) Contains(elements ...interface{}) bool {
	if len(elements) == 0 {
		return true
	}
	if list.size == 0 {
		return false
	}

	for _, value := range elements {
		found := false
		for element := list.head; element != nil; element = element.next {
			if element.value == value {
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

func (list *SingleList) IndexOf(value interface{}) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.List() {
		if element == value {
			return index
		}
	}
	return -1
}

func (list *SingleList) Sort(comparator utils.Comparator) bool {
	if list.size < 2 {
		return false
	}
	elements := list.List()
	utils.Sort(elements, comparator)
	list.Clear()
	list.Add(elements...)
	return true
}

func (list *SingleList) Swap(val_1, val_2 int) {
	if list.inRange(val_1) && list.inRange(val_2) && val_1 != val_2 {
		var elem_1, elem_2 *Node
		for i, currentElement := 0, list.head; elem_1 == nil || elem_2 == nil; i, currentElement = i+1, currentElement.next {
			switch i {
			case val_1:
				elem_1 = currentElement
			case val_2:
				elem_2 = currentElement
			}
		}
		elem_1.value, elem_2.value = elem_2.value, elem_1.value
	}
}

func (list *SingleList) Set(index int, value interface{}) {
	if !list.inRange(index) {
		if index == list.size {
			list.Add(value)
		}
		return
	}
	foundElement := list.head
	for i := 0; i != index; {
		i, foundElement = i+1, foundElement.next
	}
	foundElement.value = value
}

func (list *SingleList) Insert(index int, elements ...interface{}) {
	if !list.inRange(index) {
		if index == list.size {
			list.Add(elements...)
		}
		return
	}
	list.size += len(elements)
	var prevElement *Node
	foundElement := list.head
	for i := 0; i != index; i, foundElement = i+1, foundElement.next {
		prevElement = foundElement
	}
	if foundElement == list.head {
		oldNextElement := list.head
		for i, value := range elements {
			newElement := &Node{value: value}
			if i == 0 {
				list.head = newElement
			} else {
				prevElement.next = newElement
			}
			prevElement = newElement
		}
		prevElement.next = oldNextElement
	} else {
		oldNextElement := prevElement.next
		for _, value := range elements {
			newElement := &Node{value: value}
			prevElement.next = newElement
			prevElement = newElement
		}
		prevElement.next = oldNextElement
	}
}

//=============================================================

func (list *SingleList) inRange(index int) bool {
	return index >= 0 && index < list.size
}
