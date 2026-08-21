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

	t.Run("should parse generations", func(t *testing.T) {
		for name, expected := range map[string]InstanceSize{
			"M30":            {Family: "M", Size: 30, Generation: Generation1},
			"M30_GEN_2":      {Family: "M", Size: 30, Generation: Generation2},
			"R40_GEN_2":      {Family: "R", Size: 40, Generation: Generation2},
			"M40_NVME":       {Family: "M", Size: 40, IsNVME: true, Generation: Generation1},
			"M40_NVME_GEN_2": {Family: "M", Size: 40, IsNVME: true, Generation: Generation2},
			// M10 and M20 are offered on every generation.
			"M10": {Family: "M", Size: 10, Generation: GenerationAny},
			"M20": {Family: "M", Size: 20, Generation: GenerationAny},
		} {
			is, err := NewFromInstanceSizeName(name)

			assert.NoError(t, err, name)
			assert.Equal(t, expected, is, name)
			assert.Equal(t, name, is.String(), name)
		}
	})
}

func TestGenerationCompatibleWith(t *testing.T) {
	t.Run("should only be compatible within a generation or with the agnostic sizes", func(t *testing.T) {
		assert.True(t, Generation2.CompatibleWith(Generation2))
		assert.True(t, Generation1.CompatibleWith(Generation1))
		assert.True(t, Generation2.CompatibleWith(GenerationAny))
		assert.True(t, GenerationAny.CompatibleWith(Generation1))
		assert.False(t, Generation1.CompatibleWith(Generation2))
		assert.False(t, Generation2.CompatibleWith(Generation1))
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

	t.Run("should order a Gen2 NVME size above its non-NVME counterpart", func(t *testing.T) {
		nvme, err := NewFromInstanceSizeName("M40_NVME_GEN_2")
		assert.NoError(t, err)

		plain, err := NewFromInstanceSizeName("M40_GEN_2")
		assert.NoError(t, err)

		assert.Equal(t, 1, CompareInstanceSizes(nvme, plain))
		assert.Equal(t, -1, CompareInstanceSizes(plain, nvme))
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
