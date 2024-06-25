package utils

import (
	"context"
	"fmt"
	"time"
)

type result[T any] struct {
	value T
	err   error
}

func ExecuteWithTimeout[T any](timeout time.Duration, f func() (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan result[T], 1)
	go func() {
		value, err := f()
		done <- result[T]{value, err}
	}()

	select {
	case res := <-done:
		return res.value, res.err
	case <-ctx.Done():
		var zero T
		return zero, fmt.Errorf("function execution timed out after %v", timeout)
	}
}
