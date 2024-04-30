package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnSet(t *testing.T) {
	testCases := []struct {
		title    string
		inputs   [][]Conn
		expected []Conn
	}{
		{
			title: "Disjoint lists concatenate",
			inputs: [][]Conn{
				{{Name: "A"}, {Name: "B"}, {Name: "C"}},
				{{Name: "D"}, {Name: "E"}, {Name: "F"}},
			},
			expected: []Conn{
				{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}, {Name: "E"}, {Name: "F"},
			},
		},
		{
			title: "Common items get merged away",
			inputs: [][]Conn{
				{{Name: "A"}, {Name: "B"}, {Name: "C"}},
				{{Name: "B"}, {Name: "C"}, {Name: "D"}},
			},
			expected: []Conn{
				{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result := connSet(tc.inputs...)
			assert.Equal(t, tc.expected, result)
		})
	}
}
