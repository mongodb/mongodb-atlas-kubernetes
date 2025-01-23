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

func MapDiff[K comparable, V any](a, b map[K]V) map[K]V {
	d := make(map[K]V, len(a))
	for i, val := range a {
		if _, ok := b[i]; !ok {
			d[i] = val
		}
	}

	return d
}
