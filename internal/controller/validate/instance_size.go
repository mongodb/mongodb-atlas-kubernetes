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

type InstanceSize struct {
	Family string
	Size   int
	IsNVME bool
}

func (i *InstanceSize) String() string {
	if i.IsNVME {
		return fmt.Sprintf("%s%d_NVME", i.Family, i.Size)
	}

	return fmt.Sprintf("%s%d", i.Family, i.Size)
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

	pieces := strings.Split(instanceSizeName, "_")

	if pieces[0][0] != 'M' && pieces[0][0] != 'R' {
		return InstanceSize{}, errors.New("instance size is invalid. instance family should be M or R")
	}

	number, err := strconv.Atoi(pieces[0][1:])
	if err != nil {
		return InstanceSize{}, fmt.Errorf("instance size is invalid. %w", err)
	}

	return InstanceSize{
		Family: string(pieces[0][0]),
		Size:   number,
		IsNVME: len(pieces) == 2 && pieces[1] == "NVME",
	}, nil
}
