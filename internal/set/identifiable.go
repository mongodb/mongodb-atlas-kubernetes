package set

import (
	"reflect"
)

// DeprecatedIdentifiable is a simple interface wrapping any object which has some key field which can be used for later
// aggregation operations (grouping, intersection, difference etc)
//
// Note: this construct is DEPRECATED. Instead, use the translation layer for comparing types.
type DeprecatedIdentifiable interface {
	Identifier() interface{}
}

// DeprecatedDifference returns all 'Identifiable' elements that are in left slice and not in the right one.
// Note, that despite the parameters being declared as 'interface{}' they must contain only the elements implementing
// 'Identifiable' interface (this is needed to solve lack of covariance in Go).
//
// Note: this construct is DEPRECATED. Instead, use the translation layer for comparing types.
func DeprecatedDifference(left, right interface{}) []DeprecatedIdentifiable {
	leftIdentifiers := toIdentifiableSlice(left)
	rightIdentifiers := toIdentifiableSlice(right)

	return differenceIdentifiable(leftIdentifiers, rightIdentifiers)
}

// DeprecatedIntersection returns all 'Identifiable' elements from 'left' and 'right' slice that intersect by 'Identifier()' value.
// Each intersection is represented as a tuple of two elements - matching elements from 'left' and 'right'.
// Note, that despite the parameters being declared as 'interface{}' they must contain only the elements implementing
// 'Identifiable' interface (this is needed to solve lack of covariance in Go)
//
// Note: this construct is DEPRECATED. Instead, use the translation layer for comparing types.
func DeprecatedIntersection(left, right interface{}) [][]DeprecatedIdentifiable {
	leftIdentifiers := toIdentifiableSlice(left)
	rightIdentifiers := toIdentifiableSlice(right)

	return intersectionIdentifiable(leftIdentifiers, rightIdentifiers)
}

func differenceIdentifiable(left, right []DeprecatedIdentifiable) []DeprecatedIdentifiable {
	result := make([]DeprecatedIdentifiable, 0)
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

func intersectionIdentifiable(left, right []DeprecatedIdentifiable) [][]DeprecatedIdentifiable {
	result := make([][]DeprecatedIdentifiable, 0)
	for _, l := range left {
		for _, r := range right {
			if r.Identifier() == l.Identifier() {
				result = append(result, []DeprecatedIdentifiable{l, r})
			}
		}
	}
	return result
}

// toIdentifiableSlice uses reflection to cast the array
func toIdentifiableSlice(data interface{}) []DeprecatedIdentifiable {
	value := reflect.ValueOf(data)

	result := make([]DeprecatedIdentifiable, value.Len())
	for i := 0; i < value.Len(); i++ {
		result[i] = value.Index(i).Interface().(DeprecatedIdentifiable)
	}
	return result
}
