package toptr

func MakePtr[T comparable](value T) *T {
	return &value
}
