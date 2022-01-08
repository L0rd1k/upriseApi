package lists

type List interface {
	Get(index int) (interface{}, bool)
	Remove(index int)
	Add(elements ...interface{})
	Contains(elements ...interface{}) bool
	Swap(index_1, index_2 int)
	Insert(index int, elements ...interface{})
	Set(index int, element interface{})
	Empty() bool
	Size() int
	Clear()
	List() []interface{}
	ToString() string
}
