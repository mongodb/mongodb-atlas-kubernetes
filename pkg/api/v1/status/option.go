package status

// Option is the generic container for any information that needs to be put into the status of the Custom Resource
type Option interface {
	Value() interface{}
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
