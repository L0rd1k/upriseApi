package sets

type Set interface {
	Add(elements ...interface{})
	Remove(elements ...interface{})
	Contains(elements ...interface{}) bool
	Empty() bool
	Size() int
	Clear()
	List() []interface{} // slice of interfaces
	ToString() string
	// containers.Containers // general methods to use
}
