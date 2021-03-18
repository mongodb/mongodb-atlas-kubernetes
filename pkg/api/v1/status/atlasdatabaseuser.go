package status

// +k8s:deepcopy-gen=false

// AtlasDatabaseUserStatusOption is the option that is applied to Atlas Project Status
type AtlasDatabaseUserStatusOption func(s *AtlasDatabaseUserStatus)

func AtlasDatabaseUserSecretsOption(clusters2Secrets map[string]string) AtlasDatabaseUserStatusOption {
	return func(s *AtlasDatabaseUserStatus) {
		s.ConnectionSecrets = clusters2Secrets
	}
}

func AtlasDatabaseUserPasswordVersion(passwordVersion string) AtlasDatabaseUserStatusOption {
	return func(s *AtlasDatabaseUserStatus) {
		s.PasswordVersion = passwordVersion
	}
}

// AtlasDatabaseUserStatus defines the observed state of AtlasProject
type AtlasDatabaseUserStatus struct {
	Common            `json:",inline"`
	ConnectionSecrets map[string]string `json:"connectionSecrets,omitempty"`
	PasswordVersion   string            `json:"passwordVersion,omitempty"`
}
