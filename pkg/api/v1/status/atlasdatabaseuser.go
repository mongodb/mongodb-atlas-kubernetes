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
func AtlasDatabaseUserNameOption(name string) AtlasDatabaseUserStatusOption {
	return func(s *AtlasDatabaseUserStatus) {
		s.UserName = name
	}
}

// AtlasDatabaseUserStatus defines the observed state of AtlasProject
type AtlasDatabaseUserStatus struct {
	Common `json:",inline"`

	// ConnectionSecrets defines all the Secrets for the current database user that were created by the Atlas Operator
	// and which provide the Atlas clusters connectivity information
	ConnectionSecrets map[string]string `json:"connectionSecrets,omitempty"`

	// PasswordVersion is the 'ResourceVersion' of the password Secret that the Atlas Operator is aware of
	PasswordVersion string `json:"passwordVersion,omitempty"`

	// UserName is the current name of database user.
	UserName string `json:"name,omitempty"`
}
