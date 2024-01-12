package compat_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
)

func TestJSONSliceMerge(t *testing.T) {
	type Item struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	type OtherItem struct {
		OtherID   string `json:"id,omitempty"`
		OtherName string `json:"name,omitempty"`
	}

	tests := []struct {
		name               string
		dst, src, expected interface{}
		expectedError      error
	}{
		{
			name: "src is longer",
			dst: &[]*Item{
				{"00001", "dst1"},
				{"00002", "dst2"},
				{"00003", "dst3"},
			},
			src: []OtherItem{ // copying from different element type
				{"99999", "src1"},  // different key, different value
				{"", "src2"},       // no key, different value
				{"", ""},           // no key, no value
				{"12345", "extra"}, // extra value
			},
			expected: &[]*Item{ // kept dst element type
				{"99999", "src1"},  // key & value replaced by src
				{"00002", "src2"},  // only value replaced by src
				{"00003", "dst3"},  // untouched
				{"12345", "extra"}, // appended from src
			},
		},
		{
			name: "dst is longer",
			dst: &[]*Item{
				{"00001", "dst1"},
				{"00002", "dst2"},
				{"00003", "dst3"},
			},
			src: []OtherItem{
				{"99999", "src1"},
			},
			expected: &[]*Item{
				{"99999", "src1"}, // key & value replaced by src
				{"00002", "dst2"}, // untouched
				{"00003", "dst3"}, // untouched
			},
		},
		{
			name: "src is nil",
			dst: &[]*Item{
				{"00001", "dst1"},
				{"00002", "dst2"},
				{"00003", "dst3"},
			},
			src:           nil,
			expectedError: errors.New("src must be a slice or a pointer to slice"),
			expected: &[]*Item{
				{"00001", "dst1"}, // untouched
				{"00002", "dst2"}, // untouched
				{"00003", "dst3"}, // untouched
			},
		},
		{
			name:          "dst is nil",
			dst:           nil,
			expectedError: errors.New("dst must be a pointer to slice"),
			src: []OtherItem{
				{"99999", "src1"},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			err := JSONSliceMerge(tt.dst, tt.src)
			require.Equal(tt.expectedError, err)
			require.Equal(tt.expected, tt.dst)
		})
	}
}
