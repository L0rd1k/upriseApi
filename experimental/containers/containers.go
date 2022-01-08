package containers

type Containers interface {
	Empty() bool
	Size() int
	Clear()
	List() []interface{} // slice of interfaces
	ToString() string
}
