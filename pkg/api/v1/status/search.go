package status

import (
	"go.mongodb.org/atlas/mongodbatlas"
)

type indexStatus string

const (
	IndexStatusReady      indexStatus = "ready"
	IndexStatusInProgress indexStatus = "inProgress"
	IndexStatusFailed     indexStatus = "failed"
)

type AtlasSearch struct {
	CustomAnalyzers []string      `json:"customAnalyzers,omitempty"`
	Indexes         []*AtlasIndex `json:"indexes,omitempty"`
}

type AtlasIndex struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Database       string      `json:"database"`
	CollectionName string      `json:"collectionName"`
	Status         indexStatus `json:"status"`
	Error          string      `json:"error,omitempty"`
}

func NewStatusFromAtlas(index *mongodbatlas.SearchIndex, err error) *AtlasIndex {
	if index == nil {
		return &AtlasIndex{
			Status: IndexStatusFailed,
			Error:  err.Error(),
		}
	}

	var status indexStatus
	var errMessage string

	switch index.Status {
	case "IN_PROGRESS", "MIGRATING":
		status = IndexStatusInProgress
	case "FAILED":
		status = IndexStatusFailed
		errMessage = err.Error()
	case "STEADY":
		status = IndexStatusReady
	}

	return &AtlasIndex{
		ID:             index.IndexID,
		Name:           index.Name,
		Database:       index.Database,
		CollectionName: index.CollectionName,
		Status:         status,
		Error:          errMessage,
	}
}
