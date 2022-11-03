package status

type CustomRoleStatus string

const (
	CustomRoleStatusOK     CustomRoleStatus = "OK"
	CustomRoleStatusFailed CustomRoleStatus = "FAILED"
)

type CustomRole struct {
	// Role name which is unique
	Name string `json:"name"`
	// The status of the given custom role (OK or FAILED)
	Status CustomRoleStatus `json:"status"`
	// The message when the custom role is in the FAILED status
	Error string `json:"error,omitempty"`
}
