package lists

type List interface {
	Get(index int) (interface{}, bool)
	Remove(index int) bool
	Add(elements ...interface{})
	Contains(elements ...interface{}) bool
	Swap(index_1, index_2 int)
	Insert(index int, elements ...interface{}) bool
	Set(index int, element interface{}) bool
	Empty() bool
	Size() int
	Clear()
	List() []interface{}
	ToString() string
}
