package compat

import "encoding/json"

// JSONCopy will copy src to dst via JSON serialization/deserialization.
func JSONCopy(dst, src interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &dst)
	if err != nil {
		return err
	}

	return nil
}
