package toptr

func MakePtr[T comparable](value T) *T {
	return &value
}

func PtrValOrDefault[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}
