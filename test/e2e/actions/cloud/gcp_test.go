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
	ruleName := fmt.Sprintf("%s%s-%s", googleConnectPrefix, name, i)
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

func TestFormForwardRuleNameByAddressName_NoPrefixErr(t *testing.T) {
	addressName := "invalid-address-name"
	ruleName, err := FormForwardRuleNameByAddressName(addressName)
	assert.Error(t, err)
	assert.Empty(t, ruleName)
}

func TestFormForwardRuleNameByAddressName_NoMinusErr(t *testing.T) {
	name := "name"
	i := 1
	ruleName := fmt.Sprintf("%s%s%d", googleConnectPrefix, name, i)
	ruleName, err := FormForwardRuleNameByAddressName(ruleName)
	assert.Error(t, err)
	assert.Empty(t, ruleName)
}

func TestFormForwardRuleNameByAddressName_NoIPErr(t *testing.T) {
	name := "name"
	i := "one"
	addressName := fmt.Sprintf("%s-%s-%s", googleConnectPrefix, name, i)
	ruleName, err := FormForwardRuleNameByAddressName(addressName)
	assert.Error(t, err)
	assert.Empty(t, ruleName)
}

func TestFormForwardRuleNameByAddressName_AtoiErr(t *testing.T) {
	name := "name"
	i := "one"
	addressName := fmt.Sprintf("%s%s-ip-%s", googleConnectPrefix, name, i)
	ruleName, err := FormForwardRuleNameByAddressName(addressName)
	assert.Error(t, err)
	assert.Empty(t, ruleName)
}

func TestFormForwardRuleNameByAddressName_OK(t *testing.T) {
	name := "name"
	i := 1
	addressName := formAddressName(name, i)
	expectedAddress := formRuleName(name, i)
	ruleName, err := FormForwardRuleNameByAddressName(addressName)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, ruleName)
}
