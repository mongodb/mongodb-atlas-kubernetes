package debug

import "encoding/json"

// PrettyString is a utility function for debugging.
func PrettyString(obj interface{}) (string, error) {
	bytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MustPrettyString returns an indented string and panics if there is an error.
// the use case for this function is debugging.
func MustPrettyString(obj interface{}) string {
	s, err := PrettyString(obj)
	if err != nil {
		panic(err)
	}
	return s
}
