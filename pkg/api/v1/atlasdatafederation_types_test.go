package v1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestAtlasDataFederation_ToAtlas(t *testing.T) {
	type fields struct {
		Spec DataFederationSpec
	}
	tests := []struct {
		name    string
		fields  fields
		want    *mongodbatlas.DataFederationInstance
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Should convert all fields",
			fields: struct{ Spec DataFederationSpec }{Spec: DataFederationSpec{
				Project: common.ResourceRefNamespaced{
					Name:      "testName",
					Namespace: "testNamespace",
				},
				Name: "testName",
				CloudProviderConfig: &CloudProviderConfig{AWS: &AWSProviderConfig{
					RoleID:       "testRoleID",
					TestS3Bucket: "testS3Bucket",
				}},
				DataProcessRegion: &DataProcessRegion{
					CloudProvider: "AWS",
					Region:        "SYDNEY_AUS",
				},
				Storage: &Storage{
					Databases: []Database{
						{
							Collections: []Collection{
								{
									DataSources: []DataSource{
										{
											AllowInsecure:       true,
											Collection:          "test-collection-1",
											CollectionRegex:     "test-collection-regex",
											Database:            "test-db-1",
											DatabaseRegex:       "test-db-regex",
											DefaultFormat:       "test-format",
											Path:                "test-path",
											ProvenanceFieldName: "test-field-name",
											StoreName:           "http-test",
											Urls:                []string{"https://data.cityofnewyork.us/api/views/vfnx-vebw/rows.csv"},
										},
									},
									Name: "test-collection-1",
								},
							},
							MaxWildcardCollections: 0,
							Name:                   "test-db-1",
							Views: []View{
								{
									Name:     "test-view-1",
									Pipeline: "test-pipeline-1",
									Source:   "test-store-source",
								},
							},
						},
					},
					Stores: []Store{
						{
							Name:                     "http-test",
							Provider:                 "http",
							AdditionalStorageClasses: []string{"test-storage-class"},
							Bucket:                   "test-bucket",
							Delimiter:                ",",
							IncludeTags:              true,
							Prefix:                   "test-prefix",
							Public:                   true,
							Region:                   "SYDNEY_AUS",
						},
					},
				},
				PrivateEndpoints: []DataFederationPE{
					{
						EndpointID: "test-id",
						Provider:   "AWS",
						Type:       "DATA_LAKE",
					},
				},
			}},
			want: &mongodbatlas.DataFederationInstance{
				CloudProviderConfig: &mongodbatlas.CloudProviderConfig{AWSConfig: mongodbatlas.AwsCloudProviderConfig{
					ExternalID:        "",
					IAMAssumedRoleARN: "",
					IAMUserARN:        "",
					RoleID:            "testRoleID",
					TestS3Bucket:      "testS3Bucket",
				}},
				DataProcessRegion: &mongodbatlas.DataProcessRegion{
					CloudProvider: "AWS",
					Region:        "SYDNEY_AUS",
				},
				Storage: &mongodbatlas.DataFederationStorage{
					Databases: []*mongodbatlas.DataFederationDatabase{
						{
							Collections: []*mongodbatlas.DataFederationCollection{
								{
									DataSources: []*mongodbatlas.DataFederationDataSource{
										{
											AllowInsecure:       toptr.MakePtr(true),
											Collection:          "test-collection-1",
											CollectionRegex:     "test-collection-regex",
											Database:            "test-db-1",
											DatabaseRegex:       "test-db-regex",
											DefaultFormat:       "test-format",
											Path:                "test-path",
											ProvenanceFieldName: "test-field-name",
											StoreName:           "http-test",
											Urls:                []*string{toptr.MakePtr[string]("https://data.cityofnewyork.us/api/views/vfnx-vebw/rows.csv")},
										},
									},
									Name: "test-collection-1",
								},
							},
							MaxWildcardCollections: 0,
							Name:                   "test-db-1",
							Views: []*mongodbatlas.DataFederationDatabaseView{
								{
									Name:     "test-view-1",
									Pipeline: "test-pipeline-1",
									Source:   "test-store-source",
								},
							},
						},
					},
					Stores: []*mongodbatlas.DataFederationStore{
						{
							Name:                     "http-test",
							Provider:                 "http",
							AdditionalStorageClasses: []*string{toptr.MakePtr[string]("test-storage-class")},
							Bucket:                   "test-bucket",
							Delimiter:                ",",
							IncludeTags:              toptr.MakePtr(true),
							Prefix:                   "test-prefix",
							Region:                   "SYDNEY_AUS",
							Public:                   toptr.MakePtr(true),
						},
					},
				},
				Name: "testName",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AtlasDataFederation{
				Spec: tt.fields.Spec,
			}
			got, err := c.ToAtlas()
			if !tt.wantErr(t, err, "ToAtlas()") {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				g, _ := json.MarshalIndent(got, "", " ")
				w, _ := json.MarshalIndent(tt.want, "", " ")
				fmt.Println("GOT", string(g))
				fmt.Println("WANT", string(w))
			}

			assert.Equalf(t, tt.want, got, "ToAtlas()")
		})
	}
}
