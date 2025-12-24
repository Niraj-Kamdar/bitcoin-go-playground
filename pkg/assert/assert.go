package assert

import (
	"errors"
	"testing"
)

func Equal[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("assert.Equal failed: Expected %v, got %v", expected, actual)
	}
}

func NotEqual[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()

	if actual == expected {
		t.Errorf("assert.NotEqual failed: Expected %v, got %v", expected, actual)
	}
}

func Raises[E error](t *testing.T, fn func() (any, E), expected E) {
	t.Helper()

	_, err := fn()

	if !errors.Is(err, expected) {
		t.Errorf("assert.Raises failed: Expected %v, got %v", expected, err)
	}
}

func Less(t *testing.T, actual int, expected int) {
	t.Helper()

	if actual >= expected {
		t.Errorf("assert.Less failed: Expected %v, got %v", expected, actual)
	}
}

func GreaterOrEqual(t *testing.T, actual int, expected int) {
	t.Helper()

	if actual < expected {
		t.Errorf("assert.GreaterOrEqual failed: Expected %v, got %v", expected, actual)
	}
}
