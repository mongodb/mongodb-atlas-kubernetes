package translate

import (
	"encoding/base64"
	"fmt"
)

func secretDecode(value string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}
	return string(bytes), nil
}

func secretEncode(value string) (string, error) {
	return base64.StdEncoding.EncodeToString(([]byte)(value)), nil
}
