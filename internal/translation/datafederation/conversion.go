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

	"go.mongodb.org/atlas-sdk/v20250312021/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type DataFederation struct {
	*akov2.DataFederationSpec
	ProjectID           string
	Hostnames           []string
	CloudProviderStatus *status.DataFederationCloudProviderConfigStatus
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

	// Treat a config with no providers as absent (Atlas may return an empty
	// struct even when none is configured).
	if spec.CloudProviderConfig != nil &&
		spec.CloudProviderConfig.AWS == nil &&
		spec.CloudProviderConfig.Azure == nil &&
		spec.CloudProviderConfig.GCP == nil {
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

	// nil inner tag-sets cannot roundtrip through Atlas (they become empty slices),
	// so strip them before comparison.
	if spec.Storage != nil {
		for i := range spec.Storage.Stores {
			if rp := spec.Storage.Stores[i].ReadPreference; rp != nil {
				normalized := rp.TagSets[:0]
				for _, ts := range rp.TagSets {
					if len(ts) > 0 {
						normalized = append(normalized, ts)
					}
				}
				if len(normalized) == 0 {
					rp.TagSets = nil
				} else {
					rp.TagSets = normalized
				}
			}
		}
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

	df, err := NewDataFederation(
		&akov2.DataFederationSpec{
			CloudProviderConfig: cloudProviderConfigFromAtlas(federation.CloudProviderConfig),
			DataProcessRegion:   dataProcessRegionFromAtlas(federation.DataProcessRegion),
			Name:                federation.GetName(),
			Storage:             storageFromAtlas(federation.Storage),
		},
		federation.GetGroupId(),
		federation.GetHostnames(),
	)
	if err != nil {
		return nil, err
	}
	if df != nil {
		df.CloudProviderStatus = cloudProviderConfigStatusFromAtlas(federation.CloudProviderConfig)
	}
	return df, nil
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
			ClusterName: atlasStore.GetClusterName(),
		}
		store.AdditionalStorageClasses = append(store.AdditionalStorageClasses, atlasStore.GetAdditionalStorageClasses()...)
		if rc := atlasStore.ReadConcern; rc != nil {
			store.ReadConcern = &akov2.ReadConcern{Level: rc.GetLevel()}
		}
		if rp := atlasStore.ReadPreference; rp != nil {
			store.ReadPreference = readPreferenceFromAtlas(rp)
		}
		result.Stores = append(result.Stores, store)
	}
	return result
}

func readPreferenceFromAtlas(rp *admin.DataLakeAtlasStoreReadPreference) *akov2.ReadPreference {
	result := &akov2.ReadPreference{
		Mode:                rp.GetMode(),
		MaxStalenessSeconds: rp.GetMaxStalenessSeconds(),
	}
	for _, atlasTagSet := range rp.GetTagSets() {
		tagSet := make([]akov2.ReadPreferenceTag, 0, len(atlasTagSet))
		for _, atlasTag := range atlasTagSet {
			tagSet = append(tagSet, akov2.ReadPreferenceTag{
				Name:  atlasTag.GetName(),
				Value: atlasTag.GetValue(),
			})
		}
		result.TagSets = append(result.TagSets, tagSet)
	}
	return result
}

func readPreferenceToAtlas(rp *akov2.ReadPreference) *admin.DataLakeAtlasStoreReadPreference {
	result := &admin.DataLakeAtlasStoreReadPreference{
		Mode:                pointer.MakePtrOrNil(rp.Mode),
		MaxStalenessSeconds: pointer.MakePtrOrNil(rp.MaxStalenessSeconds),
	}
	var atlasSets [][]admin.DataLakeAtlasStoreReadPreferenceTag
	for _, tagSet := range rp.TagSets {
		if len(tagSet) == 0 {
			continue // nil/empty inner slices don't roundtrip through Atlas
		}
		atlasTagSet := make([]admin.DataLakeAtlasStoreReadPreferenceTag, 0, len(tagSet))
		for _, tag := range tagSet {
			atlasTagSet = append(atlasTagSet, admin.DataLakeAtlasStoreReadPreferenceTag{
				Name:  pointer.MakePtrOrNil(tag.Name),
				Value: pointer.MakePtrOrNil(tag.Value),
			})
		}
		atlasSets = append(atlasSets, atlasTagSet)
	}
	result.TagSets = pointer.GetOrNilIfEmpty(atlasSets)
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
					AllowInsecure:       new(dataSource.AllowInsecure),
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
			IncludeTags: new(store.IncludeTags),
			Prefix:      pointer.MakePtrOrNil(store.Prefix),
			Public:      new(store.Public),
			Region:      pointer.MakePtrOrNil(store.Region),
			ClusterName: pointer.MakePtrOrNil(store.ClusterName),
		}
		additionalStorageClasses := make([]string, 0, len(store.AdditionalStorageClasses))
		additionalStorageClasses = append(additionalStorageClasses, store.AdditionalStorageClasses...)
		atlasStore.AdditionalStorageClasses = pointer.GetOrNilIfEmpty(additionalStorageClasses)
		if store.ReadConcern != nil {
			atlasStore.ReadConcern = &admin.DataLakeAtlasStoreReadConcern{
				Level: pointer.MakePtrOrNil(store.ReadConcern.Level),
			}
		}
		if store.ReadPreference != nil {
			atlasStore.ReadPreference = readPreferenceToAtlas(store.ReadPreference)
		}
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
	if azure, ok := config.GetAzureOk(); ok {
		result.Azure = &akov2.AzureProviderConfig{
			RoleID: azure.GetRoleId(),
		}
	}
	if gcp, ok := config.GetGcpOk(); ok {
		result.GCP = &akov2.GCPProviderConfig{
			RoleID: gcp.GetRoleId(),
		}
	}
	return result
}

func cloudProviderConfigStatusFromAtlas(config *admin.DataLakeCloudProviderConfig) *status.DataFederationCloudProviderConfigStatus {
	if config == nil {
		return nil
	}
	cs := &status.DataFederationCloudProviderConfigStatus{}
	hasAny := false
	if aws, ok := config.GetAwsOk(); ok {
		cs.AWS = &status.AWSProviderConfigStatus{
			ExternalID:        aws.GetExternalId(),
			IAMAssumedRoleARN: aws.GetIamAssumedRoleARN(),
			IAMUserARN:        aws.GetIamUserARN(),
		}
		hasAny = true
	}
	if azure, ok := config.GetAzureOk(); ok {
		cs.Azure = &status.AzureProviderConfigStatus{
			AtlasAppID:         azure.GetAtlasAppId(),
			ServicePrincipalID: azure.GetServicePrincipalId(),
			TenantID:           azure.GetTenantId(),
		}
		hasAny = true
	}
	if gcp, ok := config.GetGcpOk(); ok {
		cs.GCP = &status.GCPProviderConfigStatus{
			GCPServiceAccount: gcp.GetGcpServiceAccount(),
		}
		hasAny = true
	}
	if !hasAny {
		return nil
	}
	return cs
}

// NewCloudProviderConfigStatusOption returns a status option that writes
// the Atlas-assigned read-only cloud provider fields into the DataFederation status.
func NewCloudProviderConfigStatusOption(cs *status.DataFederationCloudProviderConfigStatus) status.DataFederationStatusOption {
	return func(s *status.DataFederationStatus) {
		s.CloudProviderConfig = cs
	}
}

func cloudProviderConfigToAtlas(config *akov2.CloudProviderConfig) *admin.DataLakeCloudProviderConfig {
	if config == nil {
		return nil
	}
	result := &admin.DataLakeCloudProviderConfig{}
	hasAny := false
	if config.AWS != nil {
		result.Aws = &admin.DataLakeAWSCloudProviderConfig{
			RoleId:       config.AWS.RoleID,
			TestS3Bucket: config.AWS.TestS3Bucket,
		}
		hasAny = true
	}
	if config.Azure != nil {
		result.Azure = &admin.DataFederationAzureCloudProviderConfig{
			RoleId: config.Azure.RoleID,
		}
		hasAny = true
	}
	if config.GCP != nil {
		result.Gcp = &admin.DataFederationGCPCloudProviderConfig{
			RoleId: config.GCP.RoleID,
		}
		hasAny = true
	}
	if !hasAny {
		return nil
	}
	return result
}
