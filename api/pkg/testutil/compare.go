package testutil

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Compare allows comparing expected vs actual excluding ignoreFields and shows difference between both
func Compare(t *testing.T, expected interface{}, actual interface{}, model interface{}, ignoreFields ...string) {
	// Add time comparison options
	opts := []cmp.Option{
		cmpopts.IgnoreFields(model, ignoreFields...),
		cmp.Comparer(func(x, y time.Time) bool {
			return x.Sub(y) < 1*time.Second // Allow 1 second difference
		}),
	}

	if !cmp.Equal(expected, actual, opts...) {
		t.Errorf("\n model mismatched. \n expected: %+v \n got: %+v \n diff: %+v",
			expected, actual, cmp.Diff(expected, actual, opts...))
		t.FailNow()
	}
}
