package model

import (
	"encoding/json"
)

type PatchList struct {
	PatchList []Patch
}

type Patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func NewPatchList() *PatchList {
	return &PatchList{}
}

func (l *PatchList) PatchApplicationName(name string) *PatchList {
	l.PatchList = append(l.PatchList,
		Patch{
			Op:    "replace",
			Path:  "/metadata/name",
			Value: name,
		},
	)
	return l
}

func (l *PatchList) PatchSecret(keyName string) *PatchList {
	l.PatchList = append(l.PatchList,
		Patch{
			Op:    "replace",
			Path:  "/spec/template/spec/volumes/0/secret/secretName",
			Value: keyName,
		},
		Patch{
			Op:    "replace",
			Path:  "/spec/template/spec/containers/0/env/0/valueFrom/secretKeyRef/key",
			Value: keyName,
		},
	)
	return l
}

func (l *PatchList) GetData() []byte {
	data, _ := json.Marshal(l.PatchList)
	return data
}
