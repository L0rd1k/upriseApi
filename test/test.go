package test

import (
	"fmt"

	"github.com/L0rd1k/uprise-api/experimental/containers/lists/array"
	"github.com/L0rd1k/uprise-api/experimental/containers/lists/singlell"
	"github.com/L0rd1k/uprise-api/experimental/containers/queue/circlequeue"
	"github.com/L0rd1k/uprise-api/experimental/containers/stacks/arraystack"
)

func TestArray() {
	tstArray := array.New("Flone", "Apple", "Diving", "Bucket")
	tstArray.Add("Make")
	fmt.Println(tstArray.ToString())
	// tstArray.Sort(utils.Comparator_String)
	// fmt.Println("\n", tstArray.ToString())
}

func TestSingleLinkedList() {
	tstList := singlell.New("One", "Two", "Three", "Four")
	tstList.Remove(3)
	fmt.Println(tstList.ToString())
}

func TestArrayStack() {
	stack := arraystack.New()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	fmt.Println(stack.ToString())
}

func TestCircleQueue() {
	cqueue := circlequeue.New(3)
	cqueue.Push(1)
	cqueue.Push(2)
	cqueue.Push(3)
	fmt.Println(cqueue.ToString())
	fmt.Println("-------")
	cqueue.Push(4)
	cqueue.Push(5)
	fmt.Println(cqueue.ToString())
}
