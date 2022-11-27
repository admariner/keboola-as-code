// Package etcdop provides a framework on top of etcd low-level operations.
//
// At first, create a custom prefix using NewPrefix/NewTypedPrefix functions.
// Examples can be found in the tests. See also Key, KeyT[T], Prefix and PrefixT[T] types.
//
// Goals:
// - Reduce the risk of an error when defining an operation.
// - Distinguish between operations over one key (Key type) and several keys (Prefix type).
// - Provides Serialization composed of encode, decode and validate operations.
//
// A new operation can be defined as a:
// New<Operation>(<operation factory>, <response processor>) function.
//
// On Operation.Do(ctx, etcdClient) call :
//  - The <operation factory> is executed, result is an etcd operation.
//  - The etcd operation is executed, result is an etcd response.
//  - The <response processor> is executed, to process the etcd response.
//
// If an error occurs, it will be returned and the operation will stop.

package etcdop
