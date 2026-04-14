// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type someID struct {
	// name is a "key" field used for merging
	name string
	// some other property. Indicates which exactly object was returned by an aggregation operation
	property string
}

func newSome(name, property string) someID {
	return someID{
		name:     name,
		property: property,
	}
}

func (s someID) Identifier() any {
	return s.name
}

func Test_SetDifference(t *testing.T) {
	oneLeft := newSome("1", "left")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")
	fourRight := newSome("4", "right")

	testCases := []struct {
		left  []DeprecatedIdentifiable
		right []DeprecatedIdentifiable
		out   []DeprecatedIdentifiable
	}{
		{left: []DeprecatedIdentifiable{oneLeft, twoLeft}, right: []DeprecatedIdentifiable{twoRight, threeRight}, out: []DeprecatedIdentifiable{oneLeft}},
		{left: []DeprecatedIdentifiable{twoRight, threeRight}, right: []DeprecatedIdentifiable{oneLeft, twoLeft}, out: []DeprecatedIdentifiable{threeRight}},
		{left: []DeprecatedIdentifiable{oneLeft, twoLeft}, right: []DeprecatedIdentifiable{threeRight, fourRight}, out: []DeprecatedIdentifiable{oneLeft, twoLeft}},
		// Empty
		{left: []DeprecatedIdentifiable{}, right: []DeprecatedIdentifiable{threeRight, fourRight}, out: []DeprecatedIdentifiable{}},
		{left: []DeprecatedIdentifiable{threeRight, fourRight}, right: []DeprecatedIdentifiable{}, out: []DeprecatedIdentifiable{threeRight, fourRight}},
		// Nil
		{left: nil, right: []DeprecatedIdentifiable{threeRight, fourRight}, out: []DeprecatedIdentifiable{}},
		{left: []DeprecatedIdentifiable{threeRight, fourRight}, right: nil, out: []DeprecatedIdentifiable{threeRight, fourRight}},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, testCase.out, differenceIdentifiable(testCase.left, testCase.right))
		})
	}
}

func Test_SetDifferenceCovariant(t *testing.T) {
	// check reflection magic to solve lack of covariance in go. The arrays are declared as '[]someId' instead of
	// '[]Identifiable'
	oneLeft := newSome("1", "left")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")
	leftNotIdentifiable := []someID{oneLeft, twoLeft}
	rightNotIdentifiable := []someID{twoRight, threeRight}

	assert.Equal(t, []DeprecatedIdentifiable{oneLeft}, DeprecatedDifference(leftNotIdentifiable, rightNotIdentifiable))
	assert.Equal(t, []DeprecatedIdentifiable{threeRight}, DeprecatedDifference(rightNotIdentifiable, leftNotIdentifiable))
}

func Test_SetIntersection(t *testing.T) {
	oneLeft := newSome("1", "left")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")
	fourRight := newSome("4", "right")

	testCases := []struct {
		left  []DeprecatedIdentifiable
		right []DeprecatedIdentifiable
		out   [][]DeprecatedIdentifiable
	}{
		// intersectionIdentifiable on "2"
		{left: []DeprecatedIdentifiable{oneLeft, twoLeft}, right: []DeprecatedIdentifiable{twoRight, threeRight}, out: [][]DeprecatedIdentifiable{pair(twoLeft, twoRight)}},
		{left: []DeprecatedIdentifiable{twoRight, threeRight}, right: []DeprecatedIdentifiable{oneLeft, twoLeft}, out: [][]DeprecatedIdentifiable{pair(twoRight, twoLeft)}},
		// No intersection
		{left: []DeprecatedIdentifiable{oneLeft, twoLeft}, right: []DeprecatedIdentifiable{threeRight, fourRight}, out: [][]DeprecatedIdentifiable{}},
		{left: []DeprecatedIdentifiable{threeRight, fourRight}, right: []DeprecatedIdentifiable{oneLeft, twoLeft}, out: [][]DeprecatedIdentifiable{}},
		// Empty
		{left: []DeprecatedIdentifiable{}, right: []DeprecatedIdentifiable{threeRight, fourRight}, out: [][]DeprecatedIdentifiable{}},
		{left: []DeprecatedIdentifiable{threeRight, fourRight}, right: []DeprecatedIdentifiable{}, out: [][]DeprecatedIdentifiable{}},
		// Nil
		{left: nil, right: []DeprecatedIdentifiable{threeRight, fourRight}, out: [][]DeprecatedIdentifiable{}},
		{left: []DeprecatedIdentifiable{threeRight, fourRight}, right: nil, out: [][]DeprecatedIdentifiable{}},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, testCase.out, intersectionIdentifiable(testCase.left, testCase.right))
		})
	}
}

func Test_SetIntersectionCovariant(t *testing.T) {
	oneLeft := newSome("1", "left")
	oneRight := newSome("1", "right")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")

	// check reflection magic to solve lack of covariance in go. The arrays are declared as '[]someId' instead of
	// '[]Identifiable'
	leftNotIdentifiable := []someID{oneLeft, twoLeft}
	rightNotIdentifiable := []someID{oneRight, twoRight, threeRight}

	assert.Equal(t, [][]DeprecatedIdentifiable{pair(oneLeft, oneRight), pair(twoLeft, twoRight)}, DeprecatedIntersection(leftNotIdentifiable, rightNotIdentifiable))
	assert.Equal(t, [][]DeprecatedIdentifiable{pair(oneRight, oneLeft), pair(twoRight, twoLeft)}, DeprecatedIntersection(rightNotIdentifiable, leftNotIdentifiable))
}

func pair(left, right DeprecatedIdentifiable) []DeprecatedIdentifiable {
	return []DeprecatedIdentifiable{left, right}
}
