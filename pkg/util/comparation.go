package util

func IsEqualWithoutOrder[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[T]bool, len(a))
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		if _, ok := m[item]; !ok {
			return false
		}
	}
	return true
}

func PtrValuesEqual[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func Contains[T comparable](a []T, b T) bool {
	for _, item := range a {
		if item == b {
			return true
		}
	}
	return false
}
