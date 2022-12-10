// Package iterator provides iterator for etcd keys.
package iterator

import (
	"context"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/op"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/serde"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

type DefinitionT[T any] struct {
	config
}

type IteratorT[T any] struct {
	config
	ctx          context.Context
	client       *etcd.Client
	err          error
	start        string          // page start prefix
	page         int             // page number, start from 1
	lastIndex    int             // lastIndex in the page, 0 means empty
	currentIndex int             // currentIndex in the page, start from 0
	values       []*op.KeyValue  // values in the page
	header       *Header         // page response header
	currentValue op.KeyValueT[T] // currentValue in the page, match currentIndex
}

func NewTyped[R any](start string, serde *serde.Serde, opts ...Option) DefinitionT[R] {
	return DefinitionT[R]{config: newConfig(start, serde, opts)}
}

// Do converts iterator definition to the iterator.
func (v DefinitionT[T]) Do(ctx context.Context, client *etcd.Client) *IteratorT[T] {
	return &IteratorT[T]{ctx: ctx, client: client, config: v.config, start: v.config.prefix, page: 0, currentIndex: 0}
}

// Next returns true if there is a next value.
// False is returned if there is no next value or an error occurred.
func (v *IteratorT[T]) Next() bool {
	select {
	case <-v.ctx.Done():
		// Stop iteration if the context is done
		v.err = v.ctx.Err()
		return false
	default:
		// Is there one more item?
		if !v.nextItem() && !v.nextPage() {
			return false
		}

		// Decode item
		v.currentValue = op.KeyValueT[T]{KV: v.values[v.currentIndex]}
		if err := v.serde.Decode(v.ctx, v.currentValue.KV, &v.currentValue.Value); err != nil {
			v.err = errors.Errorf(`etcd iterator failed: cannot decode key "%s", page=%d, index=%d: %w`, v.currentValue.KV.Key, v.page, v.currentIndex, err)
		}
		return v.err == nil
	}
}

// Value returns the current value.
// It must be called after Next method.
func (v *IteratorT[T]) Value() op.KeyValueT[T] {
	if v.page == 0 {
		panic(errors.New("unexpected Value() call: Next() must be called first"))
	}
	if v.err != nil {
		panic(errors.Errorf("unexpected Value() call: %w", v.err))
	}
	return v.currentValue
}

// Header returns header of the page etcd response.
func (v *IteratorT[T]) Header() *Header {
	return v.header
}

// Err returns error. It must be checked after iterations (Next() == false).
func (v *IteratorT[T]) Err() error {
	return v.err
}

// All returns all values as a slice.
//
// The values are sorted by key in ascending order.
func (v *IteratorT[T]) All() (out op.KeyValuesT[T], err error) {
	if err = v.AllTo(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// AllTo resets the slice and add all values to the slice.
//
// The values are sorted by key in ascending order.
func (v *IteratorT[T]) AllTo(out *op.KeyValuesT[T]) (err error) {
	*out = (*out)[:0]
	for v.Next() {
		*out = append(*out, v.Value())
	}
	if err = v.Err(); err != nil {
		return err
	}
	return nil
}

// ForEachKV iterates the KVs using a callback.
func (v *IteratorT[T]) ForEachKV(fn func(value op.KeyValueT[T], header *Header) error) (err error) {
	for v.Next() {
		if err = fn(v.Value(), v.Header()); err != nil {
			return err
		}
	}
	if err = v.Err(); err != nil {
		return err
	}
	return nil
}

// ForEachValue iterates the typed values using a callback.
func (v *IteratorT[T]) ForEachValue(fn func(value T, header *Header) error) (err error) {
	for v.Next() {
		if err = fn(v.Value().Value, v.Header()); err != nil {
			return err
		}
	}
	if err = v.Err(); err != nil {
		return err
	}
	return nil
}

func (v *IteratorT[T]) nextItem() bool {
	if v.lastIndex > v.currentIndex {
		v.currentIndex++
		return true
	}
	return false
}

func (v *IteratorT[T]) nextPage() bool {
	// Is there one more page?
	if v.start == end {
		return false
	}

	// Range options
	ops := []etcd.OpOption{
		etcd.WithFromKey(),
		etcd.WithRange(etcd.GetPrefixRangeEnd(v.prefix)), // iterate to the end of the prefix
		etcd.WithLimit(int64(v.pageSize)),
		etcd.WithSort(etcd.SortByKey, etcd.SortAscend),
	}

	// Ensure atomicity
	if v.revision > 0 {
		ops = append(ops, etcd.WithRev(v.revision))
	}

	// Get page
	v.page++
	r, err := v.client.Get(v.ctx, v.start, ops...)
	if err != nil {
		v.err = errors.Errorf(`etcd iterator failed: cannot get page "%s", page=%d: %w`, v.start, v.page, err)
		return false
	}

	// Handle empty result
	v.values = r.Kvs
	v.header = r.Header
	v.lastIndex = len(v.values) - 1
	if v.lastIndex == -1 {
		return false
	}

	// Prepare next page
	if r.More {
		// Start of the next page is one key after the last key
		lastKey := string(v.values[v.lastIndex].Key)
		v.start = etcd.GetPrefixRangeEnd(lastKey)
	} else {
		v.start = end
	}

	v.currentIndex = 0
	v.revision = r.Header.Revision
	return true
}
