package atlasdeployment

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func Test_getIndicesFromAnnotations(t *testing.T) {
	t.Run("Should return empty map if annotations are empty", func(t *testing.T) {
		in := map[string]string{}
		assert.Nil(t, getIndicesFromAnnotations(in))
	})

	t.Run("Should return empty map if annotations are nil", func(t *testing.T) {
		assert.Nil(t, getIndicesFromAnnotations(nil))
	})

	t.Run("Should return valid IndexName:IndexID pairs", func(t *testing.T) {
		in := map[string]string{
			DeploymentIndicesAnnotation: "IndexOne:1,IndexTwo:2,IndexThree:3",
		}
		result := getIndicesFromAnnotations(in)
		assert.Len(t, result, 3)
		assert.Equal(t, "1", result["IndexOne"])
		assert.Equal(t, "2", result["IndexTwo"])
		assert.Equal(t, "3", result["IndexThree"])
	})

	t.Run("Should return ONLY valid IndexName:IndexID pairs", func(t *testing.T) {
		in := map[string]string{
			DeploymentIndicesAnnotation: "IndexOne:1,IndexTwo:2,IndexThree:3,IndexName4,:5",
		}
		result := getIndicesFromAnnotations(in)
		assert.Len(t, result, 4)
		assert.Equal(t, "1", result["IndexOne"])
		assert.Equal(t, "2", result["IndexTwo"])
		assert.Equal(t, "3", result["IndexThree"])
		assert.Equal(t, "5", result["IndexFive"])
	})
}

func Test_verifyAllIndicesNamesAreUnique(t *testing.T) {
	t.Run("Should return true if all indices names are unique", func(t *testing.T) {
		in := []akov2.SearchIndex{
			{
				Name: "Index-One",
			},
			{
				Name: "Index-Two",
			},
			{
				Name: "Index-Three",
			},
		}
		assert.True(t, verifyAllIndicesNamesAreUnique(in))
	})
	t.Run("Should return false if one index name appeared twice", func(t *testing.T) {
		in := []akov2.SearchIndex{
			{
				Name: "Index-One",
			},
			{
				Name: "Index-Two",
			},
			{
				Name: "Index-One",
			},
		}
		assert.False(t, verifyAllIndicesNamesAreUnique(in))
	})
}

func Test_handleSearchIndices(t *testing.T) {

}

func Test_findIndicesIntersection(t *testing.T) {
	type args struct {
		akoIndices   []*searchindex.SearchIndex
		atlasIndices []*searchindex.SearchIndex
		intersection IntersectionType
	}
	tests := []struct {
		name string
		args args
		want []searchindex.SearchIndex
	}{
		{
			name: "Should find indices to create with no indices in Atlas (exclusive left join)",
			args: args{
				akoIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				atlasIndices: nil,
				intersection: ToCreate,
			},
			want: []searchindex.SearchIndex{
				{SearchIndex: akov2.SearchIndex{Name: "1"}},
				{SearchIndex: akov2.SearchIndex{Name: "2"}},
				{SearchIndex: akov2.SearchIndex{Name: "3"}},
			},
		},
		{
			name: "Should find indices to create with some indices in Atlas (exclusive left join)",
			args: args{
				akoIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				atlasIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
				},
				intersection: ToCreate,
			},
			want: []searchindex.SearchIndex{
				{SearchIndex: akov2.SearchIndex{Name: "1"}},
				{SearchIndex: akov2.SearchIndex{Name: "3"}},
			},
		},
		{
			name: "Shouldn't find indices to create when no indices in AKO (exclusive left join)",
			args: args{
				akoIndices: nil,
				atlasIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
				},
				intersection: ToCreate,
			},
			want: nil,
		},
		{
			name: "Should find indices to update with some indices in Atlas (inner join)",
			args: args{
				akoIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				atlasIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
				},
				intersection: ToUpdate,
			},
			want: []searchindex.SearchIndex{
				{SearchIndex: akov2.SearchIndex{Name: "2"}},
			},
		},
		{
			name: "Shouldn't find indices to update with no indices in Atlas (inner join)",
			args: args{
				akoIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				atlasIndices: nil,
				intersection: ToUpdate,
			},
			want: nil,
		},
		{
			name: "Shouldn't find indices to update when no indices in AKO (inner join)",
			args: args{
				akoIndices: nil,
				atlasIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				intersection: ToUpdate,
			},
			want: nil,
		},
		{
			name: "Should find indices to delete with some indices in Atlas (exclusive right join)",
			args: args{
				akoIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				atlasIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "2"}},
				},
				intersection: ToDelete,
			},
			want: []searchindex.SearchIndex{
				{SearchIndex: akov2.SearchIndex{Name: "2"}},
			},
		},
		{
			name: "Shouldn't find indices to delete with equal indices in Atlas and in AKO (exclusive right join)",
			args: args{
				akoIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				atlasIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				intersection: ToDelete,
			},
			want: nil,
		},
		{
			name: "Shouldn't find indices to delete with equal indices in Atlas and in AKO (exclusive right join)",
			args: args{
				akoIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				atlasIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				intersection: ToDelete,
			},
			want: nil,
		},
		{
			name: "Shouldn't find indices to delete with no indices in Atlas (exclusive right join)",
			args: args{
				akoIndices: []*searchindex.SearchIndex{
					{SearchIndex: akov2.SearchIndex{Name: "1"}},
					{SearchIndex: akov2.SearchIndex{Name: "3"}},
				},
				atlasIndices: nil,
				intersection: ToDelete,
			},
			want: nil,
		},
		{
			name: "Shouldn't find indices to delete with no indices in AKO (exclusive right join)",
			args: args{
				akoIndices:   nil,
				atlasIndices: nil,
				intersection: ToDelete,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, findIndicesIntersection(tt.args.akoIndices, tt.args.atlasIndices, tt.args.intersection), "findIndicesIntersection(%v, %v, %v)", tt.args.akoIndices, tt.args.atlasIndices, tt.args.intersection)
		})
	}
}
