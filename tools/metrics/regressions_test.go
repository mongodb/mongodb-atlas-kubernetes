package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryRegressions(t *testing.T) {
	srs, err := QueryRegressions(newTestClient(), lastRecordingTime, Weekly, 5)
	assert.NoError(t, err)
	require.NotNil(t, srs)
}
