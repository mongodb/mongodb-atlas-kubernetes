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

package validate

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

		assert.EqualError(t, err, "instance size is invalid. strconv.Atoi: parsing \"Z\": invalid syntax")
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
