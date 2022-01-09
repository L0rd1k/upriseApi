package singlelinkedlist

import (
	"fmt"
	"strings"
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

//=============================================================

func (list *SingleList) inRange(index int) bool {
	return index >= 0 && index < list.size
}
