package op

import (
	"context"

	etcd "go.etcd.io/etcd/client/v3"
)

type BoolOp = ForType[bool]

// BoolMapper converts an etcd response to true/false value.
type BoolMapper func(ctx context.Context, r etcd.OpResponse) (bool, error)

// NewBoolOp wraps an operation, the result of which us true/false value.
// True means success of the operation.
func NewBoolOp(factory Factory, mapper BoolMapper) BoolOp {
	return ForType[bool]{factory: factory, mapper: mapper}
}
