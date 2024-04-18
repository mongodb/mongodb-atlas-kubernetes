package status

type AtlasSearchIndexStatus string

const (
	SearchIndexStatusReady   = "Ready"
	SearchIndexStatusError   = "Error"
	SearchIndexStatusPending = "Updating"
)

type AtlasSearchIndex struct {
	Name   string                 `json:"name"`
	ID     string                 `json:"ID"`
	Status AtlasSearchIndexStatus `json:"status"`
}
