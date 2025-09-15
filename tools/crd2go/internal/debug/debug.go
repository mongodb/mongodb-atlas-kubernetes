package debug

import "encoding/json"

func JSONize(obj any) string {
	js, err := json.MarshalIndent(obj, "  ", "  ")
	if err != nil {
		return err.Error()
	}
	return string(js)
}
