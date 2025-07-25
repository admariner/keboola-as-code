package op

import (
	"context"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

const (
	txnOpIf      = txnPartType("if")
	txnOpThen    = txnPartType("then")
	txnOpThenTxn = txnPartType("thenTxn") // thenTxnOps are separated from thenOps to avoid confusing between Then and Merge
	txnOpElse    = txnPartType("else")
	txnOpMerge   = txnPartType("merge")
)

// TxnOp provides a high-level interface built on top of etcd.Txn., see If, Then and Else methods.
//
// For more information on etcd transactions, please refer to:
// https://etcd.io/docs/v3.6/learning/api/#transaction
//
// Like other operations in this package, you can define processors using the AddProcessor method.
//
// High-level TxnOp is composed of high-level Op operations.
// This allows you to easily combine operations into an atomic transaction
// Processors defined in the operations will be executed.
//
// Another advantage is the ability to combine several TxnOp transactions into one, see Merge method.
//
// R is type of the transaction result, use NoResult type if you don't need it.
// The results of individual sub-operations can be obtained using TxnResult.SubResults,
// but a better way is to set the OnResult callback on the sub-transaction.
type TxnOp[R any] struct {
	result *R
	client etcd.KV
	errs   errors.MultiError
	parts  []txnPart
	// processors are callbacks that can react on or modify the result of an operation.
	// processors are invoked only if the etcd operation completed without server error.
	// Result methods Err/AddErr/ResetErr are used for logical errors, e.g. some unexpected value is found.
	processors []func(ctx context.Context, r *TxnResult[R])
}

// txnInterface is a marker interface for generic type TxnOp.
type txnInterface interface {
	Op
	txnParts() []txnPart
	txnInitError() error
	txnInvokeProcessors(ctx context.Context, result *resultBase)
}

type txnPartType string

type txnPart struct {
	Type txnPartType
	If   etcd.Cmp
	Op   Op
}

type lowLevelTxn[R any] struct {
	result     *R
	client     etcd.KV
	processors []func(ctx context.Context, r *TxnResult[R])
	ifs        []etcd.Cmp
	thenOps    []etcd.Op
	elseOps    []etcd.Op
}

// Txn creates an empty transaction with NoResult.
func Txn(client etcd.KV) *TxnOp[NoResult] {
	return &TxnOp[NoResult]{client: client, result: &NoResult{}, errs: errors.NewMultiError()}
}

// TxnWithResult creates an empty transaction with the result.
func TxnWithResult[R any](client etcd.KV, result *R) *TxnOp[R] {
	return &TxnOp[R]{client: client, result: result, errs: errors.NewMultiError()}
}

// MergeToTxn merges listed operations into a transaction using And method.
func MergeToTxn(client etcd.KV, ops ...Op) *TxnOp[NoResult] {
	return Txn(client).Merge(ops...)
}

func (v *TxnOp[R]) txnParts() []txnPart {
	return v.parts
}

func (v *TxnOp[R]) txnInitError() error {
	return v.errs.ErrorOrNil()
}

func (v *TxnOp[R]) txnInvokeProcessors(ctx context.Context, base *resultBase) {
	txnResult := newTxnResult(base, v.result)
	for _, p := range v.processors {
		p(ctx, txnResult)
	}
}

func (v *TxnOp[R]) Empty() bool {
	return len(v.parts) == 0
}

// If takes a list of comparison.
// If all comparisons succeed, the Then branch will be executed; otherwise, the Else branch will be executed.
func (v *TxnOp[R]) If(cs ...etcd.Cmp) *TxnOp[R] {
	for _, c := range cs {
		v.parts = append(v.parts, txnPart{Type: txnOpIf, If: c})
	}

	return v
}

// Then takes a list of operations.
// The Then operations will be executed if all If comparisons succeed.
// To add a transaction to the Then branch, use ThenTxn, or use Merge to merge transactions.
func (v *TxnOp[R]) Then(ops ...Op) *TxnOp[R] {
	for i, op := range ops {
		// Check common high-level transaction types.
		// Bulletproof check of the low-level transaction is in the "lowLevelTxn" method.
		if _, ok := op.(txnInterface); ok {
			panic(errors.Errorf(`invalid operation[%d]: op is a transaction, use Merge or ThenTxn, not Then`, i))
		}
		v.parts = append(v.parts, txnPart{Type: txnOpThen, Op: op})
	}

	return v
}

// ThenTxn adds the transaction to the Then branch.
// To merge a transaction, use Merge.
func (v *TxnOp[R]) ThenTxn(ops ...Op) *TxnOp[R] {
	for _, op := range ops {
		v.parts = append(v.parts, txnPart{Type: txnOpThenTxn, Op: op})
	}
	return v
}

// Else takes a list of operations.
// The operations in the Else branch will be executed if any of the If comparisons fail.
func (v *TxnOp[R]) Else(ops ...Op) *TxnOp[R] {
	for _, op := range ops {
		v.parts = append(v.parts, txnPart{Type: txnOpElse, Op: op})
	}
	return v
}

// Merge merges the transaction with one or more other transactions.
// If comparisons from all transactions are merged.
// The processors from all transactions are preserved and executed.
//
// For non-transaction operations, the method behaves the same as Then.
// To add a transaction to the Then branch without merging, use ThenTxn.
func (v *TxnOp[R]) Merge(ops ...Op) *TxnOp[R] {
	for _, op := range ops {
		v.parts = append(v.parts, txnPart{Type: txnOpMerge, Op: op})
	}
	return v
}

// AddError - all static errors are returned when the low level txn is composed.
// It makes error handling easier and move it to one place.
func (v *TxnOp[R]) AddError(errs ...error) *TxnOp[R] {
	v.errs.Append(errs...)
	return v
}

// AddProcessor adds a processor callback which is always executed after the transaction.
func (v *TxnOp[R]) AddProcessor(p func(ctx context.Context, r *TxnResult[R])) *TxnOp[R] {
	v.processors = append(v.processors, p)
	return v
}

// OnResult is a shortcut for the AddProcessor.
// If no error occurred yet, then the callback is executed with the result.
func (v *TxnOp[R]) OnResult(fn func(result *TxnResult[R])) *TxnOp[R] {
	return v.AddProcessor(func(_ context.Context, r *TxnResult[R]) {
		if r.Err() == nil {
			fn(r)
		}
	})
}

// SetResultTo is a shortcut for the AddProcessor.
// If no error occurred, the result of the operation is written to the target pointer,
// otherwise an empty value is written.
func (v *TxnOp[R]) SetResultTo(ptr *R) *TxnOp[R] {
	v.AddProcessor(func(ctx context.Context, r *TxnResult[R]) {
		if r.Err() == nil {
			*ptr = r.Result()
		} else {
			var empty R
			*ptr = empty
		}
	})
	return v
}

// OnFailed is a shortcut for the AddProcessor.
// If no error occurred yet and the transaction is failed, then the callback is executed.
func (v *TxnOp[R]) OnFailed(fn func(result *TxnResult[R])) *TxnOp[R] {
	return v.AddProcessor(func(_ context.Context, r *TxnResult[R]) {
		if r.Err() == nil && !r.Succeeded() {
			fn(r)
		}
	})
}

// OnSucceeded is a shortcut for the AddProcessor.
// If no error occurred yet and the transaction is succeeded, then the callback is executed.
func (v *TxnOp[R]) OnSucceeded(fn func(result *TxnResult[R])) *TxnOp[R] {
	return v.AddProcessor(func(_ context.Context, r *TxnResult[R]) {
		if r.Err() == nil && r.Succeeded() {
			fn(r)
		}
	})
}

func (v *TxnOp[R]) Do(ctx context.Context, opts ...Option) *TxnResult[R] {
	if lowLevel, err := v.lowLevelTxn(ctx); err == nil {
		result := lowLevel.Do(ctx, opts...)
		result.setMaxOps(v.maxOps(lowLevel.op().Op))
		return result
	} else {
		return newErrorTxnResult[R](err)
	}
}

// https://github.com/etcd-io/etcd/blob/e289ba30780208c1a464ea18f2f08259c75ac23a/server/etcdserver/api/v3rpc/key.go#L150-L179
func (v *TxnOp[R]) maxOps(op etcd.Op) int {
	if !op.IsTxn() {
		return 0
	}

	cmps, thenOps, elseOps := op.Txn()
	maxOps := max(len(elseOps), max(len(thenOps), len(cmps)))

	maxTxn := 0
	for _, subOp := range append(thenOps, elseOps...) {
		maxOps := v.maxOps(subOp)
		if maxOps > maxTxn {
			maxTxn = maxOps
		}
	}

	return maxOps + maxTxn
}

func (v *TxnOp[R]) Op(ctx context.Context) (LowLevelOp, error) {
	if lowLevel, err := v.lowLevelTxn(ctx); err == nil {
		return lowLevel.Op(ctx)
	} else {
		return LowLevelOp{}, err
	}
}

func (v *TxnOp[R]) lowLevelTxn(ctx context.Context) (*lowLevelTxn[R], error) {
	out := &lowLevelTxn[R]{result: v.result, client: v.client}

	if err := out.addParts(ctx, v); err != nil {
		return nil, err
	}

	// Add top-level processors
	out.processors = append(out.processors, v.processors...)

	return out, nil
}

func (v *lowLevelTxn[R]) Op(_ context.Context) (LowLevelOp, error) {
	return v.op(), nil
}

func (v *lowLevelTxn[R]) op() LowLevelOp {
	return LowLevelOp{
		Op: etcd.OpTxn(v.ifs, v.thenOps, v.elseOps),
		MapResponse: func(ctx context.Context, raw *RawResponse) (result any, err error) {
			txnResult := v.mapResponse(ctx, raw)
			return txnResult, txnResult.Err()
		},
	}
}

func (v *lowLevelTxn[R]) Do(ctx context.Context, opts ...Option) *TxnResult[R] {
	// Create low-level operation
	op := v.op()

	// Do with retry
	response, err := DoWithRetry(ctx, v.client, op.Op, opts...)
	if err != nil {
		return newErrorTxnResult[R](err)
	}

	return v.mapResponse(ctx, response)
}

func (v *lowLevelTxn[R]) addParts(ctx context.Context, txn txnInterface) error {
	if err := txn.txnInitError(); err != nil {
		return err
	}

	errs := errors.NewMultiError()
	opIndex := make(map[txnPartType]int)
	for _, part := range txn.txnParts() {
		if err := v.addPart(ctx, part); err != nil {
			err = errors.PrefixErrorf(err, `cannot create operation [%s][%d]:`, part.Type, opIndex[part.Type])
			errs.Append(err)
		}
		opIndex[part.Type]++
	}

	return errs.ErrorOrNil()
}

func (v *lowLevelTxn[R]) addPart(ctx context.Context, part txnPart) error {
	// Add if
	if part.Type == txnOpIf {
		v.ifs = append(v.ifs, part.If)
		return nil
	}

	// Merge
	if part.Type == txnOpMerge {
		return v.mergeTxn(ctx, part.Op)
	}

	// Create low-level operation
	lowLevel, err := part.Op.Op(ctx)
	if err != nil {
		return err
	}
	switch part.Type {
	case txnOpThen:
		if lowLevel.Op.IsTxn() {
			return errors.New("operation is a transaction, use Merge or ThenTxn, not Then")
		}
		v.addThen(lowLevel.Op, lowLevel.MapResponse)
	case txnOpThenTxn:
		if !lowLevel.Op.IsTxn() {
			return errors.New("operation is not a transaction, use Then, not ThenTxn")
		}
		v.addThen(lowLevel.Op, lowLevel.MapResponse)
	case txnOpElse:
		v.addElse(lowLevel.Op, lowLevel.MapResponse)
	default:
		panic(errors.Errorf(`unexpected operation type "%s"`, part.Type))
	}

	return nil
}

func (v *lowLevelTxn[R]) addThen(op etcd.Op, mapper MapFn) {
	v.thenOps = append(v.thenOps, op)

	// If the transaction is successful (then branch), the OP result will be available in the Responses[INDEX],
	// so we can call original mapper with the sub-response.
	if mapper != nil {
		index := len(v.thenOps) - 1
		v.processors = append(v.processors, func(ctx context.Context, r *TxnResult[R]) {
			if r.Succeeded() {
				rawSubResponse := r.Response().Txn().Responses[index]
				subResponse := r.Response().WithOpResponse(mapRawResponse(rawSubResponse))
				if _, err := mapper(ctx, subResponse); err != nil {
					r.AddErr(err)
				}
			}
		})
	}
}

func (v *lowLevelTxn[R]) addElse(op etcd.Op, mapper MapFn) {
	v.elseOps = append(v.elseOps, op)

	// If the transaction is NOT successful (else branch), the OP result will be available in the Responses[INDEX],
	// so we can call original mapper with the sub-response.
	if mapper != nil {
		index := len(v.elseOps) - 1
		v.processors = append(v.processors, func(ctx context.Context, r *TxnResult[R]) {
			if !r.Succeeded() {
				subResponse := r.Response().WithOpResponse(mapRawResponse(r.Response().OpResponse().Txn().Responses[index]))
				if _, err := mapper(ctx, subResponse); err != nil {
					r.AddErr(err)
				}
			}
		})
	}
}

func (v *lowLevelTxn[R]) mergeTxn(ctx context.Context, op Op) error {
	// Step down to a nested Merge operation, if there is no processor/else branch.
	if txn, ok := op.(txnInterface); ok && isSimpleTxn(txn) {
		if err := v.addParts(ctx, txn); err != nil {
			return err
		}
		v.processors = append(v.processors, func(ctx context.Context, r *TxnResult[R]) {
			txn.txnInvokeProcessors(ctx, r.resultBase)
		})
		return nil
	}

	// Create low level operation
	lowLevel, err := op.Op(ctx)
	if err != nil {
		return err
	}

	// If it is not a transaction, process it as Then
	if !lowLevel.Op.IsTxn() {
		v.addThen(lowLevel.Op, lowLevel.MapResponse)
		return nil
	}

	// Get transaction parts
	ifs, thenOps, elseOps := lowLevel.Op.Txn()

	// Merge IFs
	v.ifs = append(v.ifs, ifs...)

	// Merge THEN operations
	// The THEN branch will be applied, if all conditions (from all sub-transactions) are met.
	thenStart := len(v.thenOps)
	thenEnd := thenStart + len(thenOps)
	for _, item := range thenOps {
		v.addThen(item, nil)
	}

	// Merge ELSE operations
	// The ELSE branch will be applied only if the conditions of the sub-transaction are not met
	elseTxnPos := -1
	if len(elseOps) > 0 || len(ifs) > 0 {
		elseTxnPos = len(v.elseOps)
		v.addElse(etcd.OpTxn(ifs, nil, elseOps), nil)
	}

	// There may be a situation where neither THEN nor ELSE branch is executed:
	// It occurs, when the root transaction fails, but the reason is not in this sub-transaction.

	// On result, compose and map response that corresponds to the original sub-transaction
	// Processor from nested transactions must be invoked first.
	v.processors = append(v.processors, func(ctx context.Context, r *TxnResult[R]) {
		// Get sub-transaction response
		var subTxnResponse *etcd.TxnResponse
		switch {
		case r.Succeeded():
			subTxnResponse = &etcd.TxnResponse{
				// The entire transaction succeeded, which means that a partial transaction succeeded as well
				Succeeded: true,
				// Compose responses that corresponds to the original sub-transaction
				Responses: r.Response().Txn().Responses[thenStart:thenEnd],
			}
		case elseTxnPos >= 0:
			subTxnResponse = (*etcd.TxnResponse)(r.Response().Txn().Responses[elseTxnPos].GetResponseTxn())
			if subTxnResponse.Succeeded {
				// Skip mapper bellow, the transaction failed, but not due to a condition in the sub-transaction
				return
			}
		default:
			// Skip mapper bellow, the transaction failed, but there is no condition in the sub-transaction
			return
		}

		// Call original mapper of the sub transaction
		if _, err := lowLevel.MapResponse(ctx, r.Response().WithOpResponse(subTxnResponse.OpResponse())); err != nil {
			r.AddErr(err)
		}
	})

	return nil
}

func (v *lowLevelTxn[R]) mapResponse(ctx context.Context, raw *RawResponse) *TxnResult[R] {
	result := newTxnResult(newResultBase(raw), v.result)
	for _, p := range v.processors {
		p(ctx, result)
	}
	return result
}

// isSimpleTxn returns true if the transaction has no IF/ELSE branch.
// In that case, it is possible to merge it more easily: IF/THEN parts are added to the parent transaction.
func isSimpleTxn(txn txnInterface) bool {
	for _, part := range txn.txnParts() {
		if part.Type == txnOpIf {
			return false
		}
		if part.Type == txnOpElse {
			return false
		}
	}
	return true
}
