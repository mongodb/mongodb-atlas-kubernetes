package plugins

import (
	"github.com/mongodb/atlas2crd/pkg/processor"
)

type Plugin interface {
	processor.Processor

	Name() string
}
