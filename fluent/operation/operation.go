package operation

import (
	"github.com/MaiMee1/go-apispec/oas/v3"
)

func New(id string, opts ...Option) *oas.Operation {
	operation := &oas.Operation{
		OperationId: id,
	}
	for _, opt := range opts {
		opt.apply(operation)
	}
	return operation
}
