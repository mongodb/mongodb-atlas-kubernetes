package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanBadInputs(t *testing.T) {
	t.Run("Calling without arguments fails", func(t *testing.T) {
		assert.ErrorContains(t, Clean([]string{}), "Wrong number of arguments")
	})
	t.Run("Calling just with the command fails", func(t *testing.T) {
		assert.ErrorContains(t, Clean([]string{"clean"}), "expected 1 got 0")
	})
	t.Run("Calling with too many arguments fails", func(t *testing.T) {
		assert.ErrorContains(t, Clean([]string{"clean", "atlas", "duh"}), "expected 1 got 2")
	})
}

func TestCleanCalls(t *testing.T) {
	t.Run("Calling atlas works", func(t *testing.T) {
		cleanAtlas = func() {}
		assert.NoError(t, Clean([]string{"clean", "Atlas"}))
	})
	t.Run("Calling vpc works", func(t *testing.T) {
		cleanVPCs = func() {}
		assert.NoError(t, Clean([]string{"clean", "VPC"}))
	})
	t.Run("Calling atlas fails due to missing credentials", func(t *testing.T) {
		cleanPEs = func() {}
		assert.NoError(t, Clean([]string{"clean", "Pe"}))
	})
}
