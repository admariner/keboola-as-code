package op

import "time"

type AtomicResult[R any] struct {
	result      *R
	error       error
	attempt     int
	elapsedTime time.Duration
}

func (v AtomicResult[R]) Result() R {
	var empty R
	if v.error != nil || v.result == nil {
		return empty
	}
	return *v.result
}

func (v AtomicResult[R]) Err() error {
	return v.error
}

func (v AtomicResult[R]) ResultOrErr() (R, error) {
	var empty R
	if v.error != nil {
		return empty, v.error
	}
	if v.result == nil {
		return empty, nil
	}
	return *v.result, nil
}

func (v AtomicResult[R]) Attempt() int {
	return v.attempt
}

func (v AtomicResult[R]) ElapsedTime() time.Duration {
	return v.elapsedTime
}
