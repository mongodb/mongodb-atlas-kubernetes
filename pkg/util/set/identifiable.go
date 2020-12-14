package set

import (
	"reflect"
)

// Identifiable is a simple interface wrapping any object which has some key field which can be used for later
// aggregation operations (grouping, intersection, difference etc)
type Identifiable interface {
	Identifier() interface{}
}

// Difference returns all 'Identifiable' elements that are in left slice and not in the right one
func Difference(left, right []Identifiable) []Identifiable {
	result := make([]Identifiable, 0)
	for _, l := range left {
		found := false
		for _, r := range right {
			if r.Identifier() == l.Identifier() {
				found = true
				break
			}
		}
		if !found {
			result = append(result, l)
		}
	}
	return result
}

// Intersection returns all 'Identifiable' elements from 'left' and 'right' slice that intersect by 'Identifier()'
//value. Each intersection is represented as a tuple of two elements - matching elements from 'left' and 'right'
func Intersection(left, right []Identifiable) [][]Identifiable {
	result := make([][]Identifiable, 0)
	for _, l := range left {
		for _, r := range right {
			if r.Identifier() == l.Identifier() {
				result = append(result, []Identifiable{l, r})
			}
		}
	}
	return result
}

// SetDifferenceGeneric is a convenience function solving lack of covariance in Go: it allows to pass the arrays declared
// as some types implementing 'Identifiable' and find the difference between them
// Important: the arrays past must declare types implementing 'Identifiable'!
func SetDifferenceGeneric(left, right interface{}) []Identifiable {
	leftIdentifiers := toIdentifiableSlice(left)
	rightIdentifiers := toIdentifiableSlice(right)

	return Difference(leftIdentifiers, rightIdentifiers)
}

// SetIntersectionGeneric is a convenience function solving lack of covariance in Go: it allows to pass the arrays declared
// as some types implementing 'Identifiable' and find the intersection between them
// Important: the arrays past must declare types implementing 'Identifiable'!
func SetIntersectionGeneric(left, right interface{}) [][]Identifiable {
	leftIdentifiers := toIdentifiableSlice(left)
	rightIdentifiers := toIdentifiableSlice(right)

	return Intersection(leftIdentifiers, rightIdentifiers)
}

// toIdentifiableSlice uses reflection to cast the array
func toIdentifiableSlice(data interface{}) []Identifiable {
	value := reflect.ValueOf(data)

	result := make([]Identifiable, value.Len())
	for i := 0; i < value.Len(); i++ {
		result[i] = value.Index(i).Interface().(Identifiable)
	}
	return result
}
