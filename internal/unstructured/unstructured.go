package unstructured

import (
	"encoding/json"
)

func TypedFromUnstructured[U any, T any](unstructured U) (*T, error) {
	asJSON, err := json.Marshal(unstructured)
	if err != nil {
		return nil, err
	}

	var asObj T
	err = json.Unmarshal(asJSON, &asObj)
	if err != nil {
		return nil, err
	}

	return &asObj, nil
}
