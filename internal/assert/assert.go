package assert

import "testing"

// use generic function to test actual and expected results
// as long as test used has same type, this should help with testing before we run
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}
