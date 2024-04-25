package atlasdeployment

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func Test_getIndexesFromAnnotations(t *testing.T) {
	t.Run("Should return empty map if annotations are empty", func(t *testing.T) {
		in := map[string]string{}
		assert.Nil(t, getIndexesFromAnnotations(in))
	})

	t.Run("Should return empty map if annotations are nil", func(t *testing.T) {
		assert.Nil(t, getIndexesFromAnnotations(nil))
	})

	t.Run("Should return valid IndexName:IndexID pairs", func(t *testing.T) {
		in := map[string]string{
			DeploymentIndexesAnnotation: "IndexOne:1,IndexTwo:2,IndexThree:3",
		}
		result := getIndexesFromAnnotations(in)
		assert.Len(t, result, 3)
		assert.Equal(t, "1", result["IndexOne"])
		assert.Equal(t, "2", result["IndexTwo"])
		assert.Equal(t, "3", result["IndexThree"])
	})

	t.Run("Should return ONLY valid IndexName:IndexID pairs", func(t *testing.T) {
		in := map[string]string{
			DeploymentIndexesAnnotation: "IndexOne:1,IndexTwo:2,IndexThree:3,IndexName4,:5",
		}
		result := getIndexesFromAnnotations(in)
		assert.Len(t, result, 4)
		assert.Equal(t, "1", result["IndexOne"])
		assert.Equal(t, "2", result["IndexTwo"])
		assert.Equal(t, "3", result["IndexThree"])
		assert.Equal(t, "5", result["IndexFive"])
	})
}

func Test_verifyAllIndexesNamesAreUnique(t *testing.T) {
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
		assert.True(t, verifyAllIndexesNamesAreUnique(in))
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
		assert.False(t, verifyAllIndexesNamesAreUnique(in))
	})
}

func Test_handleSearchIndexes(t *testing.T) {

}

func Test_findIndexesIntersection(t *testing.T) {
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
			assert.Equalf(t, tt.want, findIndexesIntersection(tt.args.akoIndices, tt.args.atlasIndices, tt.args.intersection), "findIndexesIntersection(%v, %v, %v)", tt.args.akoIndices, tt.args.atlasIndices, tt.args.intersection)
		})
	}
}

func Test_getIndexesFromDeploymentStatus(t *testing.T) {
	tests := []struct {
		name             string
		deploymentStatus status.AtlasDeploymentStatus
		want             map[string]string
	}{
		{
			name: "Should return valid indexes for some valid indexes in the status",
			deploymentStatus: status.AtlasDeploymentStatus{
				SearchIndexes: []status.DeploymentSearchIndexStatus{
					{
						Name:      "FirstIndex",
						ID:        "1",
						Status:    "",
						ConfigRef: common.ResourceRefNamespaced{},
						Message:   "",
					},
					{
						Name:      "SecondIndex",
						ID:        "2",
						Status:    "",
						ConfigRef: common.ResourceRefNamespaced{},
						Message:   "",
					},
				},
			},
			want: map[string]string{
				"FirstIndex":  "1",
				"SecondIndex": "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getIndexesFromDeploymentStatus(tt.deploymentStatus), "getIndexesFromDeploymentStatus(%v)", tt.deploymentStatus)
		})
	}
}
