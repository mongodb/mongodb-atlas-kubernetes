package identifiable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type someId struct {
	// name is a "key" field used for merging
	name string
	// some other property. Indicates which exactly object was returned by an aggregation operation
	property string
}

func newSome(name, property string) someId {
	return someId{
		name:     name,
		property: property,
	}
}

func (s someId) Identifier() interface{} {
	return s.name
}

func Test_SetDifference(t *testing.T) {
	oneLeft := newSome("1", "left")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")
	fourRight := newSome("4", "right")

	left := []Identifiable{oneLeft, twoLeft}
	right := []Identifiable{twoRight, threeRight}

	assert.Equal(t, []Identifiable{oneLeft}, SetDifference(left, right))
	assert.Equal(t, []Identifiable{threeRight}, SetDifference(right, left))

	left = []Identifiable{oneLeft, twoLeft}
	right = []Identifiable{threeRight, fourRight}
	assert.Equal(t, left, SetDifference(left, right))

	left = []Identifiable{}
	right = []Identifiable{threeRight, fourRight}
	assert.Empty(t, SetDifference(left, right))
	assert.Equal(t, right, SetDifference(right, left))

	left = nil
	right = []Identifiable{threeRight, fourRight}
	assert.Empty(t, SetDifference(left, right))
	assert.Equal(t, right, SetDifference(right, left))
}

func Test_SetDifferenceCovariant(t *testing.T) {
	// check reflection magic to solve lack of covariance in go. The arrays are declared as '[]someId' instead of
	// '[]Identifiable'
	oneLeft := newSome("1", "left")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")
	leftNotIdentifiable := []someId{oneLeft, twoLeft}
	rightNotIdentifiable := []someId{twoRight, threeRight}

	assert.Equal(t, []Identifiable{oneLeft}, SetDifferenceGeneric(leftNotIdentifiable, rightNotIdentifiable))
	assert.Equal(t, []Identifiable{threeRight}, SetDifferenceGeneric(rightNotIdentifiable, leftNotIdentifiable))
}

func Test_SetIntersection(t *testing.T) {
	oneLeft := newSome("1", "left")
	oneRight := newSome("1", "right")
	twoLeft := newSome("2", "left")
	twoRight := newSome("2", "right")
	threeRight := newSome("3", "right")
	fourRight := newSome("4", "right")

	left := []Identifiable{oneLeft, twoLeft}
	right := []Identifiable{twoRight, threeRight}

	assert.Equal(t, [][]Identifiable{pair(twoLeft, twoRight)}, SetIntersection(left, right))
	assert.Equal(t, [][]Identifiable{pair(twoRight, twoLeft)}, SetIntersection(right, left))

	left = []Identifiable{oneLeft, twoLeft}
	right = []Identifiable{threeRight, fourRight}
	assert.Empty(t, SetIntersection(left, right))
	assert.Empty(t, SetIntersection(right, left))

	left = []Identifiable{}
	right = []Identifiable{threeRight, fourRight}
	assert.Empty(t, SetIntersection(left, right))
	assert.Empty(t, SetIntersection(right, left))

	left = nil
	right = []Identifiable{threeRight, fourRight}
	assert.Empty(t, SetIntersection(left, right))
	assert.Empty(t, SetIntersection(right, left))

	// check reflection magic to solve lack of covariance in go. The arrays are declared as '[]someId' instead of
	// '[]Identifiable'
	leftNotIdentifiable := []someId{oneLeft, twoLeft}
	rightNotIdentifiable := []someId{oneRight, twoRight, threeRight}

	assert.Equal(t, [][]Identifiable{pair(oneLeft, oneRight), pair(twoLeft, twoRight)}, SetIntersectionGeneric(leftNotIdentifiable, rightNotIdentifiable))
	assert.Equal(t, [][]Identifiable{pair(oneRight, oneLeft), pair(twoRight, twoLeft)}, SetIntersectionGeneric(rightNotIdentifiable, leftNotIdentifiable))
}

func pair(left, right Identifiable) []Identifiable {
	return []Identifiable{left, right}
}
