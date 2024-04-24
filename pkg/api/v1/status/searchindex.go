package status

type AtlasSearchIndexConfigStatus struct {
	Projects []SearchIndexConfigProject
}

type SearchIndexConfigProject struct {
	// Unique identifier of the project inside atlas
	ID string `json:"id,omitempty"`
	// Name given to the project
	Name string `json:"name,omitempty"`
}
