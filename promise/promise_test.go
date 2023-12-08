package promise_test

import (
	"errors"
	"testing"

	"github.com/JureBevc/cyanic/promise"
)

func TestPromiseError(t *testing.T) {
	_, err := promise.NewPromise[int](func() (int, error) {
		return 123, errors.New("Test error")
	}).Await()
	if err == nil {
		t.Fatalf("Expected promise error, got %s\n", err)
	}
}

func TestPromiseAwaitValue(t *testing.T) {

	expected := 2

	result, err := promise.NewPromise[int](func() (int, error) {
		return expected, nil
	}).Await()
	if err != nil {
		t.Fatalf("Expected no promise error, got %s\n", err)
	}
	if result != expected {
		t.Fatalf("Expected %d, got %d\n", expected, result)
	}
}

func TestPromiseThenValue(t *testing.T) {
	expected := 2
	promise.NewPromise[int](func() (int, error) {
		return 2, nil
	}).Then(func(result int) {
		if expected != result {
			t.Fatalf("Expected %d, got %d\n", expected, result)
		}
	}, func(err error) {
		t.Fatalf("Expected no promise error, got %s\n", err)
	})
}

func TestPromiseThenError(t *testing.T) {
	promise.NewPromise[int](func() (int, error) {
		return 123, errors.New("Test error")
	}).Then(func(result int) {
		t.Fatalf("Expected error, 'Then' function should not execute")
	}, func(err error) {
		if err == nil {
			t.Fatalf("Expected promise error, got %s\n", err)
		}
	})
}
