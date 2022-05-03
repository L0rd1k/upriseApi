package circlequeue

import (
	"fmt"
	"strings"

	"github.com/L0rd1k/uprise-api/experimental/containers/queue"
)

type CircleQueue struct {
	elements []interface{}
	front    int
	peak     int
	isFull   bool
	size     int
	maxSize  int
}

var _ queue.Queue = (*CircleQueue)(nil)

func New(size int) *CircleQueue {
	if size < 1 {
		panic("Undefined queue size")
	}
	queue := &CircleQueue{maxSize: size}
	queue.Clear()
	return queue
}

func (queue *CircleQueue) Push(element interface{}) {
	if queue.Full() {
		queue.Pop()
	}
	queue.elements[queue.peak] = element
	queue.peak = queue.peak + 1
	if queue.peak >= queue.maxSize {
		queue.peak = 0
	}
	if queue.peak == queue.front {
		queue.isFull = true
	}

	// if queue.peak < queue.front {
	// 	queue.size = queue.maxSize - queue.front + queue.peak
	// } else if queue.peak == queue.front {
	// 	if queue.isFull {
	// 		queue.size = queue.maxSize
	// 	}
	// 	queue.size = 0
	// } else {
	// 	queue.size = queue.peak - queue.front
	// }
	queue.size = queue.getNewSize()
}

func (queue *CircleQueue) getNewSize() int {
	if queue.peak < queue.front {
		return queue.maxSize - queue.front + queue.peak
	} else if queue.peak == queue.front {
		if queue.isFull {
			return queue.maxSize
		}
		return 0
	}
	return queue.peak - queue.front
}

func (queue *CircleQueue) List() []interface{} {
	elements := make([]interface{}, queue.Size(), queue.Size())
	for i := 0; i < queue.Size(); i++ {
		elements[i] = queue.elements[(queue.front+i)%queue.maxSize]
	}
	return elements
}

func (queue *CircleQueue) Peek() (elements interface{}, status bool) {
	if queue.Empty() {
		return nil, false
	}
	return queue.elements[queue.front], true
}

func (queue *CircleQueue) ToString() string {
	str := ""
	var elements []string
	for _, element := range queue.List() {
		elements = append(elements, fmt.Sprintf("%v", element))
	}
	str += strings.Join(elements, " ")
	return str
}

func (queue *CircleQueue) Pop() (element interface{}, status bool) {
	if queue.Empty() {
		return nil, false
	}
	element, status = queue.elements[queue.front], true
	if element != nil {
		queue.elements[queue.front] = nil
		queue.front = queue.front + 1
		if queue.front >= queue.maxSize {
			queue.front = 0
		}
		queue.isFull = false
	}
	queue.size = queue.size - 1
	return
}

func (queue *CircleQueue) Empty() bool {
	return queue.Size() == 0
}

func (queue *CircleQueue) Size() int {
	return queue.size
}

func (queue *CircleQueue) Full() bool {
	return queue.Size() == queue.maxSize
}

func (queue *CircleQueue) Clear() {
	queue.elements = make([]interface{}, queue.maxSize, queue.maxSize)
	queue.front = 0
	queue.peak = 0
	queue.isFull = false
	queue.size = 0
}
