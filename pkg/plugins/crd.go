package plugins

type Crd struct {
	NoOp
}

func (c Crd) Name() string {
	return "crd"
}
