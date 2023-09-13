package slices

import "slices"

func Sort[S ~[]E, E any](s S, sortFunc func(a, b E) int) []E {
	clone := slices.Clone(s)
	slices.SortFunc(clone, sortFunc)
	return clone
}
