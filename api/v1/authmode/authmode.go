package authmode

type AuthMode string

const (
	Scram AuthMode = "SCRAM"
	X509  AuthMode = "X509"
)

type AuthModes []AuthMode

func (authModes AuthModes) CheckAuthMode(modeToCheck AuthMode) bool {
	for _, mode := range authModes {
		if mode == modeToCheck {
			return true
		}
	}

	return false
}

func (authModes *AuthModes) AddAuthMode(modeToAdd AuthMode) {
	found := false
	for _, mode := range *authModes {
		if mode == modeToAdd {
			found = true
			break
		}
	}

	if !found {
		*authModes = append(*authModes, modeToAdd)
	}
}

func (authModes *AuthModes) RemoveAuthMode(modeToRemove AuthMode) {
	var result AuthModes
	for _, mode := range *authModes {
		if mode != modeToRemove {
			result = append(result, mode)
		}
	}
	*authModes = result
}
