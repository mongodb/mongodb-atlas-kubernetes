package set

func FromSlice[K comparable, T any](list []T, hash func(item T) K) map[K]T {
	m := map[K]T{}

	for _, item := range list {
		m[hash(item)] = item
	}

	return m
}
