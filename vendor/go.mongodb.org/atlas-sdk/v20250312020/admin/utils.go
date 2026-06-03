// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"reflect"
	"time"
)

// PtrBool is a helper routine that returns a pointer to given boolean value.
func PtrBool(v bool) *bool { return &v }

// PtrInt is a helper routine that returns a pointer to given integer value.
func PtrInt(v int) *int { return &v }

// PtrInt32 is a helper routine that returns a pointer to given integer value.
func PtrInt32(v int32) *int32 { return &v }

// PtrInt64 is a helper routine that returns a pointer to given integer value.
func PtrInt64(v int64) *int64 { return &v }

// PtrFloat32 is a helper routine that returns a pointer to given float value.
func PtrFloat32(v float32) *float32 { return &v }

// PtrFloat64 is a helper routine that returns a pointer to given float value.
func PtrFloat64(v float64) *float64 { return &v }

// PtrString is a helper routine that returns a pointer to given string value.
func PtrString(v string) *string { return &v }

// PtrTime is helper routine that returns a pointer to given Time value.
func PtrTime(v time.Time) *time.Time { return &v }

// IsNil checks if an input is nil
func IsNil(i any) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	case reflect.Array:
		return reflect.ValueOf(i).IsZero()
	}
	return false
}

// SetOrNil returns the address of the given value or nil if it equals defaultValue
func SetOrNil[T comparable](val T, defaultValue T) *T {
	if val == defaultValue {
		return nil
	}
	return &val
}

// GetOrDefault returns the value of a pointer or a default value
func GetOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// TimePtrToStringPtr converts the time to the string pointer.
func TimePtrToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	res := TimeToString(*t)
	return &res
}

// TimeToString returns a RFC3339 date time string format.
// The resulting format is identical to the format returned by Atlas API, documented as ISO 8601 timestamp format in UTC.
// It also returns decimals in seconds (up to nanoseconds) if available.
// Example formats: "2023-07-18T16:12:23Z", "2023-07-18T16:12:23.456Z"
func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

// StringToTime parses the given string and returns the resulting time.
// The expected format is identical to the format returned by Atlas API, documented as ISO 8601 timestamp format in UTC.
// Example formats: "2023-07-18T16:12:23Z", "2023-07-18T16:12:23.456Z"
func StringToTime(val string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, val)
}

// Int64PtrToIntPtr convert int64 to int pointer.
func Int64PtrToIntPtr(i64 *int64) *int {
	if i64 == nil {
		return nil
	}

	i := int(*i64)
	return &i
}
