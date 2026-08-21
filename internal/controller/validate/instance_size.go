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
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const gen2Suffix = "_GEN_2"

// Generation is the hardware generation an instance size belongs to. Sizes of different
// generations are not interchangeable, but they are not ordered relative to each other
// either, so generation is a compatibility dimension rather than a comparable one.
type Generation int

const (
	// GenerationAny covers the sizes offered on every generation, so it is compatible
	// with all of them. Atlas treats M10 and M20 this way.
	GenerationAny Generation = iota
	Generation1
	Generation2
)

func (g Generation) CompatibleWith(other Generation) bool {
	return g == GenerationAny || other == GenerationAny || g == other
}

type InstanceSize struct {
	Family     string
	Size       int
	IsNVME     bool
	Generation Generation
}

func (i *InstanceSize) String() string {
	name := fmt.Sprintf("%s%d", i.Family, i.Size)

	if i.IsNVME {
		name += "_NVME"
	}

	if i.Generation == Generation2 {
		name += gen2Suffix
	}

	return name
}

func CompareInstanceSizes(is1 InstanceSize, is2 InstanceSize) int {
	if is1.Family != is2.Family {
		if is1.Family == "M" {
			return -1
		} else {
			return 1
		}
	}

	if is1.Size != is2.Size {
		if is1.Size < is2.Size {
			return -1
		} else {
			return 1
		}
	}

	if is1.IsNVME != is2.IsNVME {
		if !is1.IsNVME {
			return -1
		} else {
			return 1
		}
	}

	return 0
}

func NewFromInstanceSizeName(instanceSizeName string) (InstanceSize, error) {
	if len(instanceSizeName) < 2 {
		return InstanceSize{}, errors.New("instance size is invalid")
	}

	// The generation suffix is stripped first because NVMe sizes carry both suffixes,
	// as in M40_NVME_GEN_2.
	generation := Generation1
	base := instanceSizeName

	if strings.HasSuffix(base, gen2Suffix) {
		generation = Generation2
		base = strings.TrimSuffix(base, gen2Suffix)
	}

	pieces := strings.Split(base, "_")

	if pieces[0][0] != 'M' && pieces[0][0] != 'R' {
		return InstanceSize{}, errors.New("instance size is invalid. instance family should be M or R")
	}

	number, err := strconv.Atoi(pieces[0][1:])
	if err != nil {
		return InstanceSize{}, fmt.Errorf("instance size is invalid. %w", err)
	}

	if generation == Generation1 && pieces[0][0] == 'M' && (number == 10 || number == 20) {
		generation = GenerationAny
	}

	return InstanceSize{
		Family:     string(pieces[0][0]),
		Size:       number,
		IsNVME:     len(pieces) == 2 && pieces[1] == "NVME",
		Generation: generation,
	}, nil
}
