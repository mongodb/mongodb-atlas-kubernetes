package status

type AtlasSearchIndexConfigStatus struct {
	Common `json:",inline"`
}

// +kubebuilder:object:generate=false

type AtlasSearchIndexConfigStatusOption func(s *AtlasSearchIndexConfigStatus)
