package stacks

type Stack interface {
	Push(element interface{})
	Pop(elements interface{}, status bool)
	Top(elements interface{}, status bool)
	Empty() bool
	Size() int
	Clear()
	List() []interface{} // slice of interfaces
	ToString() string
	// containers.Containers // general methods to use
}
