package test

import (
	"fmt"

	"github.com/L0rd1k/uprise-api/experimental/containers/lists/array"
	"github.com/L0rd1k/uprise-api/experimental/containers/lists/singlell"
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
