package contract

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	for _, tc := range []struct {
		title    string
		name     string
		focus    string
		enabled  bool
		expected error
	}{
		{
			title:    "empty name, focus and disabled skips all",
			expected: errors.New("AKO_CONTRACT_TEST is unset"),
		},
		{
			title:   "empty name, focus and enabled does not skip",
			enabled: true,
		},
		{
			title:    "disabled skips regardles of focus matching",
			name:     "target",
			focus:    "target",
			expected: errors.New("AKO_CONTRACT_TEST is unset"),
		},
		{
			title:   "enabled with no focus does not skip",
			enabled: true,
			name:    "target",
		},
		{
			title:   "enabled with no focus does not skip",
			enabled: true,
			name:    "target",
		},
		{
			title:    "enabled with non matching focus skips",
			enabled:  true,
			name:     "something else",
			focus:    "target",
			expected: errors.New("test \"something else\" does not contain focus string \"target\""),
		},
		{
			title:   "enabled with matching focus does not skip",
			enabled: true,
			name:    "target",
			focus:   "target",
		},
		{
			title:   "enabled matching a sub-target focus does not skip",
			enabled: true,
			name:    "some target phrase",
			focus:   "target",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			err := skipCheck(tc.name, tc.focus, tc.enabled)
			assert.Equal(t, tc.expected, err)
		})
	}
}
