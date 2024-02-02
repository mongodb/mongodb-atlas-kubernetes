package pointer

// SetOrNil returns the address of the given value or nil if it equals defaultValue
func SetOrNil[T comparable](val T, defaultValue T) *T {
	if val == defaultValue {
		return nil
	}
	return &val
}

// GetOrDefault returns the value of a pointer or a default value
func GetOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// MakePtr returns a pointer to the given value
func MakePtr[T any](value T) *T {
	return &value
}
