package status

type Prometheus struct {
	// +optional
	Scheme string `json:"scheme,omitempty"`
	// +optional
	DiscoveryURL string `json:"prometheusDiscoveryURL,omitempty"`
}
