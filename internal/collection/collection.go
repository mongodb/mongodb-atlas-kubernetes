package collection

func CopyWithSkip[T comparable](list []T, skip T) []T {
	newList := make([]T, 0, len(list))

	for _, item := range list {
		if item != skip {
			newList = append(newList, item)
		}
	}

	return newList
}

func Keys[K comparable, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))

	for k := range m {
		s = append(s, k)
	}

	return s
}

func FirstKeyFromMap[K comparable, V any](m map[K]V) K {
	for k := range m {
		return k
	}

	return *new(K)
}
