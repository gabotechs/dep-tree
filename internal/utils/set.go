package utils

type Set[T comparable] map[T]struct{}

func SetFromSlice[T comparable](arr []T) Set[T] {
	s := make(map[T]struct{})
	for _, v := range arr {
		s[v] = struct{}{}
	}
	return s
}

func (s *Set[T]) Has(el T) bool {
	_, ok := (*s)[el]
	return ok
}
