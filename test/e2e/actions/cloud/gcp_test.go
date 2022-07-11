package cloud

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormAddressNameByRuleName_NoPrefixErr(t *testing.T) {
	ruleName := "invalid-rule-name"
	addressName, err := FormAddressNameByRuleName(ruleName)
	assert.Error(t, err)
	assert.Empty(t, addressName)
}

func TestFormAddressNameByRuleName_NoMinusErr(t *testing.T) {
	name := "name"
	i := 1
	ruleName := fmt.Sprintf("%s%s%d", googleConnectPrefix, name, i)
	addressName, err := FormAddressNameByRuleName(ruleName)
	assert.Error(t, err)
	assert.Empty(t, addressName)
}

func TestFormAddressNameByRuleName_AtoiErr(t *testing.T) {
	name := "name"
	i := "one"
	ruleName := fmt.Sprintf("%s%s%s", googleConnectPrefix, name, i)
	addressName, err := FormAddressNameByRuleName(ruleName)
	assert.Error(t, err)
	assert.Empty(t, addressName)
}

func TestFormAddressNameByRuleName_OK(t *testing.T) {
	name := "name"
	i := 1
	ruleName := formRuleName(name, i)
	expectedAddress := formAddressName(name, i)
	addressName, err := FormAddressNameByRuleName(ruleName)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, addressName)
}
