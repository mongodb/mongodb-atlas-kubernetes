package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFromInstanceSizeName(t *testing.T) {
	t.Run("should return error when instance size name is invalid", func(t *testing.T) {
		is, err := NewFromInstanceSizeName("a")

		assert.EqualError(t, err, "instance size is invalid")
		assert.Empty(t, is)
	})

	t.Run("should return error when instance is from wrong family", func(t *testing.T) {
		is, err := NewFromInstanceSizeName("Z10")

		assert.EqualError(t, err, "instance size is invalid. instance family should be M or R")
		assert.Empty(t, is)
	})

	t.Run("should return error when instance is malformed", func(t *testing.T) {
		is, err := NewFromInstanceSizeName("MZ")

		assert.EqualError(t, err, "instance size is invalid. &{%!e(string=Atoi) %!e(string=Z) %!e(*errors.errorString=&{invalid syntax})}")
		assert.Empty(t, is)
	})

	t.Run("should return a general instance size parsed", func(t *testing.T) {
		is, err := NewFromInstanceSizeName("M10")

		assert.NoError(t, err)
		assert.Equal(
			t,
			InstanceSize{
				Family: "M",
				Size:   10,
				IsNVME: false,
			},
			is,
		)
	})

	t.Run("should return a NVME instance size parsed", func(t *testing.T) {
		is, err := NewFromInstanceSizeName("M10_NVME")

		assert.NoError(t, err)
		assert.Equal(
			t,
			InstanceSize{
				Family: "M",
				Size:   10,
				IsNVME: true,
			},
			is,
		)
	})
}

func TestCompareInstanceSizes(t *testing.T) {
	t.Run("should return -1 when instance 1 family is less than instance 2 family", func(t *testing.T) {
		assert.Equal(
			t,
			-1,
			CompareInstanceSizes(
				InstanceSize{
					Family: "M",
					Size:   10,
				},
				InstanceSize{
					Family: "R",
					Size:   10,
				},
			),
		)
	})

	t.Run("should return 1 when instance 1 family is greater than instance 2 family", func(t *testing.T) {
		assert.Equal(
			t,
			1,
			CompareInstanceSizes(
				InstanceSize{
					Family: "R",
					Size:   10,
				},
				InstanceSize{
					Family: "M",
					Size:   10,
				},
			),
		)
	})

	t.Run("should return -1 when instance 1 size is less than instance 2 size", func(t *testing.T) {
		assert.Equal(
			t,
			-1,
			CompareInstanceSizes(
				InstanceSize{
					Family: "M",
					Size:   10,
				},
				InstanceSize{
					Family: "M",
					Size:   20,
				},
			),
		)
	})

	t.Run("should return 1 when instance 1 size is greater than instance 2 size", func(t *testing.T) {
		assert.Equal(
			t,
			1,
			CompareInstanceSizes(
				InstanceSize{
					Family: "M",
					Size:   20,
				},
				InstanceSize{
					Family: "M",
					Size:   10,
				},
			),
		)
	})

	t.Run("should return -1 when instance 1 is not NVME and instance 2 is NVME", func(t *testing.T) {
		assert.Equal(
			t,
			-1,
			CompareInstanceSizes(
				InstanceSize{
					Family: "M",
					Size:   50,
				},
				InstanceSize{
					Family: "M",
					Size:   50,
					IsNVME: true,
				},
			),
		)
	})

	t.Run("should return -1 when instance 1 is NVME and instance 2 is not NVME", func(t *testing.T) {
		assert.Equal(
			t,
			1,
			CompareInstanceSizes(
				InstanceSize{
					Family: "M",
					Size:   50,
					IsNVME: true,
				},
				InstanceSize{
					Family: "M",
					Size:   50,
				},
			),
		)
	})

	t.Run("should return 0 when instance 1 and 2 are equal", func(t *testing.T) {
		assert.Equal(
			t,
			0,
			CompareInstanceSizes(
				InstanceSize{
					Family: "M",
					Size:   50,
					IsNVME: true,
				},
				InstanceSize{
					Family: "M",
					Size:   50,
					IsNVME: true,
				},
			),
		)
	})
}
