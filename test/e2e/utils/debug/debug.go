package debug

import "encoding/json"

// PrettyString is a utility function for displaying an indented json structure as a sring.
func PrettyString(obj interface{}) string {
	return string(PrettyBytes(obj))
}

// PrettyBytes is a utility function for displaying an indented json structure as a byte array.
func PrettyBytes(obj interface{}) []byte {
	bytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil
	}
	return bytes
}
