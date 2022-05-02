package arraystack

import (
	"fmt"
	"strings"

	"github.com/L0rd1k/uprise-api/experimental/containers/lists/array"
	"github.com/L0rd1k/uprise-api/experimental/containers/stacks"
)

type ArrayStack struct {
	arr *array.Array
}

var _ stacks.Stack = (*ArrayStack)(nil)

func New() *ArrayStack {
	return &ArrayStack{arr: array.New()}
}

func (stack *ArrayStack) Push(element interface{}) {
	stack.arr.Add(element)
}

func (stack *ArrayStack) Pop() (element interface{}, status bool) {
	element, status = stack.arr.Get(stack.arr.Size() - 1)
	stack.arr.Remove(stack.arr.Size() - 1)
	return
}

func (stack *ArrayStack) Top() (elements interface{}, status bool) {
	elements, status = stack.arr.Get(stack.arr.Size() - 1)
	return elements, status
}

func (stack *ArrayStack) Empty() bool {
	return stack.arr.Empty()
}

func (stack *ArrayStack) Size() int {
	return stack.arr.Size()
}

func (stack *ArrayStack) Clear() {
	stack.arr.Clear()
}

func (stack *ArrayStack) List() []interface{} {
	size := stack.arr.Size()
	elements := make([]interface{}, size, size)
	for i := 1; i <= size; i++ {
		elements[size-1], _ = stack.arr.Get(i - 1)
	}
	return elements
}

func (stack *ArrayStack) ToString() string {
	var str = ""
	items := []string{}
	for _, value := range stack.arr.List() {
		items = append(items, fmt.Sprintf("%v", value))
	}
	str += strings.Join(items, "\n")
	return str
}
