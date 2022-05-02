package stacks

type Stack interface {
	Push(element interface{})
	Pop() (elements interface{}, status bool)
	Top() (elements interface{}, status bool)
	Empty() bool
	Size() int
	Clear()
	List() []interface{}
	ToString() string
}
