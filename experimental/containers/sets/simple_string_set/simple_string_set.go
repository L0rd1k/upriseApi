package simple_string_set

import (
	"errors"
	"fmt"
)

type Set struct {
	Elements map[string]struct{}
}

func NewSet() *Set {
	var set Set
	set.Elements = make(map[string]struct{})
	return &set
}

func (s *Set) Add(elem string) {
	s.Elements[elem] = struct{}{}
}

func (s *Set) Delete(elem string) error {
	if _, exists := s.Elements[elem]; !exists {
		return errors.New("no elements in set")
	}
	delete(s.Elements, elem)
	return nil
}

func (s *Set) Contains(elem string) bool {
	_, exists := s.Elements[elem]
	return exists
}

func (s *Set) List() {
	for itr, _ := range s.Elements {
		fmt.Println(itr)
	}
}
