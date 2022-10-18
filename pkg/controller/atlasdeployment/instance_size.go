package atlasdeployment

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
		return InstanceSize{}, fmt.Errorf("instance size is invalid. %e", err)
	}

	return InstanceSize{
		Family: string(pieces[0][0]),
		Size:   number,
		IsNVME: len(pieces) == 2 && pieces[1] == "NVME",
	}, nil
}
