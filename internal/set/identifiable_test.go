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

func (s someID) Identifier() interface{} {
	return s.name
}

func Test_SetDifference(t *testing.T) {
	oneLeft := newSome("1", "left")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")
	fourRight := newSome("4", "right")

	testCases := []struct {
		left  []Identifiable
		right []Identifiable
		out   []Identifiable
	}{
		{left: []Identifiable{oneLeft, twoLeft}, right: []Identifiable{twoRight, threeRight}, out: []Identifiable{oneLeft}},
		{left: []Identifiable{twoRight, threeRight}, right: []Identifiable{oneLeft, twoLeft}, out: []Identifiable{threeRight}},
		{left: []Identifiable{oneLeft, twoLeft}, right: []Identifiable{threeRight, fourRight}, out: []Identifiable{oneLeft, twoLeft}},
		// Empty
		{left: []Identifiable{}, right: []Identifiable{threeRight, fourRight}, out: []Identifiable{}},
		{left: []Identifiable{threeRight, fourRight}, right: []Identifiable{}, out: []Identifiable{threeRight, fourRight}},
		// Nil
		{left: nil, right: []Identifiable{threeRight, fourRight}, out: []Identifiable{}},
		{left: []Identifiable{threeRight, fourRight}, right: nil, out: []Identifiable{threeRight, fourRight}},
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

	assert.Equal(t, []Identifiable{oneLeft}, Difference(leftNotIdentifiable, rightNotIdentifiable))
	assert.Equal(t, []Identifiable{threeRight}, Difference(rightNotIdentifiable, leftNotIdentifiable))
}

func Test_SetIntersection(t *testing.T) {
	oneLeft := newSome("1", "left")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")
	fourRight := newSome("4", "right")

	testCases := []struct {
		left  []Identifiable
		right []Identifiable
		out   [][]Identifiable
	}{
		// intersectionIdentifiable on "2"
		{left: []Identifiable{oneLeft, twoLeft}, right: []Identifiable{twoRight, threeRight}, out: [][]Identifiable{pair(twoLeft, twoRight)}},
		{left: []Identifiable{twoRight, threeRight}, right: []Identifiable{oneLeft, twoLeft}, out: [][]Identifiable{pair(twoRight, twoLeft)}},
		// No intersection
		{left: []Identifiable{oneLeft, twoLeft}, right: []Identifiable{threeRight, fourRight}, out: [][]Identifiable{}},
		{left: []Identifiable{threeRight, fourRight}, right: []Identifiable{oneLeft, twoLeft}, out: [][]Identifiable{}},
		// Empty
		{left: []Identifiable{}, right: []Identifiable{threeRight, fourRight}, out: [][]Identifiable{}},
		{left: []Identifiable{threeRight, fourRight}, right: []Identifiable{}, out: [][]Identifiable{}},
		// Nil
		{left: nil, right: []Identifiable{threeRight, fourRight}, out: [][]Identifiable{}},
		{left: []Identifiable{threeRight, fourRight}, right: nil, out: [][]Identifiable{}},
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

	assert.Equal(t, [][]Identifiable{pair(oneLeft, oneRight), pair(twoLeft, twoRight)}, Intersection(leftNotIdentifiable, rightNotIdentifiable))
	assert.Equal(t, [][]Identifiable{pair(oneRight, oneLeft), pair(twoRight, twoLeft)}, Intersection(rightNotIdentifiable, leftNotIdentifiable))
}

func pair(left, right Identifiable) []Identifiable {
	return []Identifiable{left, right}
}
