package provider

import "time"

type operationError interface {
	error

	Retryable() bool
	RetryAfter() time.Duration
}

type FatalError struct {
	Err error
}

func (e FatalError) Error() string {
	return e.Err.Error()
}

func (e FatalError) Unwrap() error {
	return e.Err
}

func (e FatalError) Retryable() bool {
	return false
}

func (e FatalError) RetryAfter() time.Duration {
	return -1
}

type RetryableError struct {
	Err   error
	After time.Duration
}

func (e RetryableError) Error() string {
	return e.Err.Error()
}

func (e RetryableError) Unwrap() error {
	return e.Err
}

func (e RetryableError) Retryable() bool {
	return true
}

func (e RetryableError) RetryAfter() time.Duration {
	return e.After
}
