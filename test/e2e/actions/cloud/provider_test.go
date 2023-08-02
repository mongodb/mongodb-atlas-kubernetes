package cloud

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prefix(cidr string) (*regexp.Regexp, error) {
	parts := strings.SplitN(cidr, ".", 4)
	prefix := strings.Join(parts[:3], ".")
	return regexp.Compile(fmt.Sprintf("^%s", prefix))
}

func TestDefaultAWSConfig(t *testing.T) {
	cfg := getAWSConfigDefaults()

	prefixRegexp, err := prefix(cfg.CIDR)
	require.NoError(t, err)
	assert.Regexp(t, prefixRegexp, cfg.Subnets[Subnet1Name])
	assert.Regexp(t, prefixRegexp, cfg.Subnets[Subnet2Name])
	assert.NotEqual(t, cfg.Subnets[Subnet1Name], cfg.Subnets[Subnet2Name])
}

func TestDefaultGoogleConfig(t *testing.T) {
	cfg := getGCPConfigDefaults()

	prefixRegexp, err := prefix(cfg.Subnets[Subnet1Name])
	require.NoError(t, err)
	assert.Regexp(t, prefixRegexp, cfg.Subnets[Subnet2Name])
	assert.NotEqual(t, cfg.Subnets[Subnet1Name], cfg.Subnets[Subnet2Name])
}

func TestDefaultAzureConfig(t *testing.T) {
	cfg := getAzureConfigDefaults()

	prefixRegexp, err := prefix(cfg.CIDR)
	require.NoError(t, err)
	assert.Regexp(t, prefixRegexp, cfg.Subnets[Subnet1Name])
	assert.Regexp(t, prefixRegexp, cfg.Subnets[Subnet2Name])
	assert.NotEqual(t, cfg.Subnets[Subnet1Name], cfg.Subnets[Subnet2Name])
}
