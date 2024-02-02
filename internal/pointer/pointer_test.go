package pointer

import (
	"testing"
)

type testCase[T comparable] struct {
	name         string
	val          T
	defaultValue T
	wantNil      bool
}

func assertSetOrNil[T comparable](t *testing.T, tc testCase[T]) {
	t.Run(tc.name, func(t *testing.T) {
		ptr := SetOrNil(tc.val, tc.defaultValue)
		if gotNil := ptr == nil; gotNil != tc.wantNil {
			t.Errorf("got nil %t, want %t", gotNil, tc.wantNil)
		}
	})
}

func TestSetOrNil(t *testing.T) {
	assertSetOrNil(t, testCase[int]{
		name:         "non default int",
		val:          1,
		defaultValue: 0,
		wantNil:      false,
	})

	assertSetOrNil(t, testCase[int]{
		name:         "default int",
		val:          0,
		defaultValue: 0,
		wantNil:      true,
	})

	assertSetOrNil(t, testCase[string]{
		name:         "non default string",
		val:          "hello",
		defaultValue: "",
		wantNil:      false,
	})

	assertSetOrNil(t, testCase[string]{
		name:         "default string",
		val:          "",
		defaultValue: "",
		wantNil:      true,
	})
}
