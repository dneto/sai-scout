package maps

func MapBy[S ~[]E, E any, K comparable](s S, mapfunc func(E) K) map[K]E {
	m := make(map[K]E)
	for _, e := range s {
		m[mapfunc(e)] = e
	}
	return m
}
