package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnSet(t *testing.T) {
	testCases := []struct {
		title    string
		inputs   [][]Connection
		expected []Connection
	}{
		{
			title: "Disjoint lists concatenate",
			inputs: [][]Connection{
				{{Name: "A"}, {Name: "B"}, {Name: "C"}},
				{{Name: "D"}, {Name: "E"}, {Name: "F"}},
			},
			expected: []Connection{
				{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}, {Name: "E"}, {Name: "F"},
			},
		},
		{
			title: "Common items get merged away",
			inputs: [][]Connection{
				{{Name: "A"}, {Name: "B"}, {Name: "C"}},
				{{Name: "B"}, {Name: "C"}, {Name: "D"}},
			},
			expected: []Connection{
				{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result := connectionSet(tc.inputs...)
			assert.Equal(t, tc.expected, result)
		})
	}
}
