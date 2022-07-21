package stringutil

// Contains returns true if there is at least one string in `slice`
// that is equal to `s`.
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
