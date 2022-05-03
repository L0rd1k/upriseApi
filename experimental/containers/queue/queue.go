package queue

type Queue interface {
	Push(value interface{})
	Pop() (element interface{}, status bool)
	Peek() (element interface{}, status bool)
	Empty() bool
	Size() int
	Clear()
	List() []interface{}
	ToString() string
}
