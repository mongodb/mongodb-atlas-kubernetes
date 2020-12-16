package status

import "reflect"

// Option is the generic container for any information that needs to be put into the status of the Custom Resource.
// This is the way to handle some random data that need to be written to status
type Option interface {
	Value() interface{}
}

// GetOption finds the option by its type
func GetOption(statusOptions []Option, targetOption Option) (Option, bool) {
	for _, s := range statusOptions {
		if reflect.TypeOf(s) == reflect.TypeOf(targetOption) {
			return s, true
		}
	}
	return nil, false
}

// IDOption describes the ID of some Atlas resource
type IDOption struct {
	ID string
}

func NewIDOption(id string) IDOption {
	return IDOption{ID: id}
}

func (o IDOption) Value() interface{} {
	return o.ID
}
