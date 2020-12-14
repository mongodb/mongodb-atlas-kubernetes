package set

import (
	"reflect"
)

// Identifiable is a simple interface wrapping any object which has some key field which can be used for later
// aggregation operations (grouping, intersection, difference etc)
type Identifiable interface {
	Identifier() interface{}
}

// Difference returns all 'Identifiable' elements that are in left slice and not in the right one.
// Note, that despite the parameters being declared as 'interface{}' they must contain only the elements implementing
// 'Identifiable' interface (this is needed to solve lack of covariance in Go).
func Difference(left, right interface{}) []Identifiable {
	leftIdentifiers := toIdentifiableSlice(left)
	rightIdentifiers := toIdentifiableSlice(right)

	return differenceIdentifiable(leftIdentifiers, rightIdentifiers)
}

// Intersection returns all 'Identifiable' elements from 'left' and 'right' slice that intersect by 'Identifier()' value.
// Each intersection is represented as a tuple of two elements - matching elements from 'left' and 'right'.
// Note, that despite the parameters being declared as 'interface{}' they must contain only the elements implementing
// 'Identifiable' interface (this is needed to solve lack of covariance in Go)
func Intersection(left, right interface{}) [][]Identifiable {
	leftIdentifiers := toIdentifiableSlice(left)
	rightIdentifiers := toIdentifiableSlice(right)

	return intersectionIdentifiable(leftIdentifiers, rightIdentifiers)
}

func differenceIdentifiable(left, right []Identifiable) []Identifiable {
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

func intersectionIdentifiable(left, right []Identifiable) [][]Identifiable {
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

// toIdentifiableSlice uses reflection to cast the array
func toIdentifiableSlice(data interface{}) []Identifiable {
	value := reflect.ValueOf(data)

	result := make([]Identifiable, value.Len())
	for i := 0; i < value.Len(); i++ {
		result[i] = value.Index(i).Interface().(Identifiable)
	}
	return result
}
