// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datafederation

import (
	"fmt"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type DataFederation struct {
	*akov2.DataFederationSpec
	ProjectID string
	Hostnames []string
}

// SpecEqualsTo returns true if the spec of the data federation instance semantically equals to the given one.
// Note: it assumes the spec is already normalized.
func (df *DataFederation) SpecEqualsTo(target *DataFederation) bool {
	var dfSpecCopy, targetSpecCopy *akov2.DataFederationSpec
	if df != nil {
		dfSpecCopy = df.DataFederationSpec.DeepCopy()
	}
	if target != nil {
		targetSpecCopy = target.DataFederationSpec.DeepCopy()
	}
	return reflect.DeepEqual(pruneSpec(dfSpecCopy), pruneSpec(targetSpecCopy))
}

func pruneSpec(spec *akov2.DataFederationSpec) *akov2.DataFederationSpec {
	if spec == nil {
		return nil
	}

	// Atlas embeds AWS config as a value, AKO embeds AWS config as a pointer,
	// hence treat the absence of both in AKO equally.
	if spec.CloudProviderConfig == nil || spec.CloudProviderConfig.AWS == nil {
		spec.CloudProviderConfig = nil
	}

	// ignore project references, they are not sent to/from Atlas.
	var emptyRef common.ResourceRefNamespaced
	spec.Project = emptyRef

	// private endpoints are sub-resources, they have their own conversion and are not part of the data federation entity.
	spec.PrivateEndpoints = nil

	// normalize nested empty stores/database slices
	if spec.Storage != nil && (len(spec.Storage.Stores) == 0 && len(spec.Storage.Databases) == 0) {
		spec.Storage = nil
	}

	return spec
}

func NewDataFederation(spec *akov2.DataFederationSpec, projectID string, hostnames []string) (*DataFederation, error) {
	if spec == nil {
		return nil, nil
	}

	specCopy := spec.DeepCopy()
	if err := cmp.Normalize(specCopy); err != nil {
		return nil, fmt.Errorf("failed to normalize data federation spec: %w", err)
	}

	return &DataFederation{
		DataFederationSpec: specCopy,
		ProjectID:          projectID,
		Hostnames:          hostnames,
	}, nil
}

func toAtlas(df *DataFederation) *admin.DataLakeTenant {
	if df == nil || df.DataFederationSpec == nil {
		return nil
	}

	return &admin.DataLakeTenant{
		GroupId:             pointer.MakePtrOrNil(df.ProjectID),
		CloudProviderConfig: cloudProviderConfigToAtlas(df.CloudProviderConfig),
		DataProcessRegion:   dataProcessRegionToAtlas(df.DataProcessRegion),
		Name:                pointer.MakePtrOrNil(df.Name),
		Storage:             storageToAtlas(df.Storage),
	}
}

func fromAtlas(federation *admin.DataLakeTenant) (*DataFederation, error) {
	if federation == nil {
		return nil, nil
	}

	return NewDataFederation(
		&akov2.DataFederationSpec{
			CloudProviderConfig: cloudProviderConfigFromAtlas(federation.CloudProviderConfig),
			DataProcessRegion:   dataProcessRegionFromAtlas(federation.DataProcessRegion),
			Name:                federation.GetName(),
			Storage:             storageFromAtlas(federation.Storage),
		},
		federation.GetGroupId(),
		federation.GetHostnames(),
	)
}

func storageFromAtlas(storage *admin.DataLakeStorage) *akov2.Storage {
	if storage == nil {
		return nil
	}

	result := &akov2.Storage{}
	for _, atlasDB := range storage.GetDatabases() {
		db := akov2.Database{
			MaxWildcardCollections: atlasDB.GetMaxWildcardCollections(),
			Name:                   atlasDB.GetName(),
		}
		for _, atlasCollection := range atlasDB.GetCollections() {
			collection := akov2.Collection{
				Name: atlasCollection.GetName(),
			}
			for _, atlasDataSource := range atlasCollection.GetDataSources() {
				dataSource := akov2.DataSource{
					AllowInsecure:       atlasDataSource.GetAllowInsecure(),
					Collection:          atlasDataSource.GetCollection(),
					CollectionRegex:     atlasDataSource.GetCollectionRegex(),
					Database:            atlasDataSource.GetDatabase(),
					DatabaseRegex:       atlasDataSource.GetDatabaseRegex(),
					DefaultFormat:       atlasDataSource.GetDefaultFormat(),
					Path:                atlasDataSource.GetPath(),
					ProvenanceFieldName: atlasDataSource.GetProvenanceFieldName(),
					StoreName:           atlasDataSource.GetStoreName(),
				}
				dataSource.Urls = append(dataSource.Urls, atlasDataSource.GetUrls()...)
				collection.DataSources = append(collection.DataSources, dataSource)
			}
			db.Collections = append(db.Collections, collection)
		}
		for _, atlasView := range atlasDB.GetViews() {
			db.Views = append(db.Views, akov2.View{
				Name:     atlasView.GetName(),
				Pipeline: atlasView.GetPipeline(),
				Source:   atlasView.GetSource(),
			})
		}
		result.Databases = append(result.Databases, db)
	}

	for _, atlasStore := range storage.GetStores() {
		store := akov2.Store{
			Name:        atlasStore.GetName(),
			Provider:    atlasStore.GetProvider(),
			Bucket:      atlasStore.GetBucket(),
			Delimiter:   atlasStore.GetDelimiter(),
			IncludeTags: atlasStore.GetIncludeTags(),
			Prefix:      atlasStore.GetPrefix(),
			Public:      atlasStore.GetPublic(),
			Region:      atlasStore.GetRegion(),
		}
		store.AdditionalStorageClasses = append(store.AdditionalStorageClasses, atlasStore.GetAdditionalStorageClasses()...)
		result.Stores = append(result.Stores, store)
	}
	return result
}

func storageToAtlas(storage *akov2.Storage) *admin.DataLakeStorage {
	if storage == nil {
		return nil
	}

	result := &admin.DataLakeStorage{}
	databases := make([]admin.DataLakeDatabaseInstance, 0, len(storage.Databases))
	for _, db := range storage.Databases {
		atlasDB := admin.DataLakeDatabaseInstance{
			MaxWildcardCollections: pointer.MakePtrOrNil(db.MaxWildcardCollections),
			Name:                   pointer.MakePtrOrNil(db.Name),
		}
		atlasCollections := make([]admin.DataLakeDatabaseCollection, 0, len(db.Collections))
		for _, collection := range db.Collections {
			atlasCollection := admin.DataLakeDatabaseCollection{
				Name: pointer.MakePtrOrNil(collection.Name),
			}
			atlasDataSources := make([]admin.DataLakeDatabaseDataSourceSettings, 0, len(collection.DataSources))
			for _, dataSource := range collection.DataSources {
				atlasDataSource := admin.DataLakeDatabaseDataSourceSettings{
					AllowInsecure:       pointer.MakePtr(dataSource.AllowInsecure),
					Collection:          pointer.MakePtrOrNil(dataSource.Collection),
					CollectionRegex:     pointer.MakePtrOrNil(dataSource.CollectionRegex),
					Database:            pointer.MakePtrOrNil(dataSource.Database),
					DatabaseRegex:       pointer.MakePtrOrNil(dataSource.DatabaseRegex),
					DefaultFormat:       pointer.MakePtrOrNil(dataSource.DefaultFormat),
					Path:                pointer.MakePtrOrNil(dataSource.Path),
					ProvenanceFieldName: pointer.MakePtrOrNil(dataSource.ProvenanceFieldName),
					StoreName:           pointer.MakePtrOrNil(dataSource.StoreName),
				}
				atlasDataSource.Urls = pointer.GetOrNilIfEmpty(append([]string{}, dataSource.Urls...))
				atlasDataSources = append(atlasDataSources, atlasDataSource)
			}
			atlasCollection.DataSources = pointer.GetOrNilIfEmpty(atlasDataSources)
			atlasCollections = append(atlasCollections, atlasCollection)
		}
		atlasDB.Collections = pointer.GetOrNilIfEmpty(atlasCollections)
		atlasViews := make([]admin.DataLakeApiBase, 0, len(db.Views))
		for _, view := range db.Views {
			atlasViews = append(atlasViews, admin.DataLakeApiBase{
				Name:     pointer.MakePtrOrNil(view.Name),
				Pipeline: pointer.MakePtrOrNil(view.Pipeline),
				Source:   pointer.MakePtrOrNil(view.Source),
			})
		}
		atlasDB.Views = pointer.GetOrNilIfEmpty(atlasViews)
		databases = append(databases, atlasDB)
	}
	result.Databases = pointer.GetOrNilIfEmpty(databases)

	stores := make([]admin.DataLakeStoreSettings, 0, len(storage.Stores))
	for _, store := range storage.Stores {
		atlasStore := admin.DataLakeStoreSettings{
			Name:        pointer.MakePtrOrNil(store.Name),
			Provider:    store.Provider,
			Bucket:      pointer.MakePtrOrNil(store.Bucket),
			Delimiter:   pointer.MakePtrOrNil(store.Delimiter),
			IncludeTags: pointer.MakePtr(store.IncludeTags),
			Prefix:      pointer.MakePtrOrNil(store.Prefix),
			Public:      pointer.MakePtr(store.Public),
			Region:      pointer.MakePtrOrNil(store.Region),
		}
		additionalStorageClasses := make([]string, 0, len(store.AdditionalStorageClasses))
		additionalStorageClasses = append(additionalStorageClasses, store.AdditionalStorageClasses...)
		atlasStore.AdditionalStorageClasses = pointer.GetOrNilIfEmpty(additionalStorageClasses)
		stores = append(stores, atlasStore)
	}
	result.Stores = pointer.GetOrNilIfEmpty(stores)
	return result
}

func dataProcessRegionFromAtlas(region *admin.DataLakeDataProcessRegion) *akov2.DataProcessRegion {
	if region == nil {
		return nil
	}
	return &akov2.DataProcessRegion{
		CloudProvider: region.GetCloudProvider(),
		Region:        region.GetRegion(),
	}
}

func dataProcessRegionToAtlas(region *akov2.DataProcessRegion) *admin.DataLakeDataProcessRegion {
	if region == nil {
		return nil
	}
	return &admin.DataLakeDataProcessRegion{
		CloudProvider: region.CloudProvider,
		Region:        region.Region,
	}
}

func cloudProviderConfigFromAtlas(config *admin.DataLakeCloudProviderConfig) *akov2.CloudProviderConfig {
	if config == nil {
		return nil
	}
	result := &akov2.CloudProviderConfig{}
	if aws, ok := config.GetAwsOk(); ok {
		result.AWS = &akov2.AWSProviderConfig{
			RoleID:       aws.GetRoleId(),
			TestS3Bucket: aws.GetTestS3Bucket(),
		}
	}
	return result
}

func cloudProviderConfigToAtlas(config *akov2.CloudProviderConfig) *admin.DataLakeCloudProviderConfig {
	if config == nil || config.AWS == nil {
		return nil
	}
	return &admin.DataLakeCloudProviderConfig{
		Aws: &admin.DataLakeAWSCloudProviderConfig{
			RoleId:       config.AWS.RoleID,
			TestS3Bucket: config.AWS.TestS3Bucket,
		},
	}
}
