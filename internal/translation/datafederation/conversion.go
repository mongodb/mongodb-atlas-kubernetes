package datafederation

import (
	"fmt"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
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

	// skip "SkipRoleValidation" field as it is a request parameter, not a returned body from/to Atlas.
	spec.SkipRoleValidation = false

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
					DatasetName:         atlasDataSource.GetDatasetName(),
					DatasetPrefix:       atlasDataSource.GetDatasetPrefix(),
					DefaultFormat:       atlasDataSource.GetDefaultFormat(),
					Path:                atlasDataSource.GetPath(),
					ProvenanceFieldName: atlasDataSource.GetProvenanceFieldName(),
					StoreName:           atlasDataSource.GetStoreName(),
					TrimLevel:           atlasDataSource.GetTrimLevel(),
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
			Name:           atlasStore.GetName(),
			Provider:       atlasStore.GetProvider(),
			Bucket:         atlasStore.GetBucket(),
			Delimiter:      atlasStore.GetDelimiter(),
			IncludeTags:    atlasStore.GetIncludeTags(),
			Prefix:         atlasStore.GetPrefix(),
			Public:         atlasStore.GetPublic(),
			Region:         atlasStore.GetRegion(),
			ClusterName:    atlasStore.GetClusterName(),
			AllowInsecure:  atlasStore.GetAllowInsecure(),
			DefaultFormat:  atlasStore.GetDefaultFormat(),
			ReadConcern:    readConcernFromAtlas(atlasStore.ReadConcern),
			ReadPreference: readPreferenceFromAtlas(atlasStore.ReadPreference),
		}
		store.Urls = append(store.Urls, atlasStore.GetUrls()...)
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
					DatasetName:         pointer.MakePtrOrNil(dataSource.DatasetName),
					DatasetPrefix:       pointer.MakePtrOrNil(dataSource.DatasetPrefix),
					DefaultFormat:       pointer.MakePtrOrNil(dataSource.DefaultFormat),
					Path:                pointer.MakePtrOrNil(dataSource.Path),
					ProvenanceFieldName: pointer.MakePtrOrNil(dataSource.ProvenanceFieldName),
					StoreName:           pointer.MakePtrOrNil(dataSource.StoreName),
					TrimLevel:           pointer.MakePtrOrNil(dataSource.TrimLevel),
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
			Name:           pointer.MakePtrOrNil(store.Name),
			Provider:       store.Provider,
			Bucket:         pointer.MakePtrOrNil(store.Bucket),
			Delimiter:      pointer.MakePtrOrNil(store.Delimiter),
			IncludeTags:    pointer.MakePtr(store.IncludeTags),
			Prefix:         pointer.MakePtrOrNil(store.Prefix),
			Public:         pointer.MakePtr(store.Public),
			Region:         pointer.MakePtrOrNil(store.Region),
			ClusterName:    pointer.MakePtrOrNil(store.ClusterName),
			AllowInsecure:  pointer.MakePtr(store.AllowInsecure),
			DefaultFormat:  pointer.MakePtrOrNil(store.DefaultFormat),
			ReadConcern:    readConcernToAtlas(store.ReadConcern),
			ReadPreference: readPreferenceToAtlas(store.ReadPreference),
		}
		atlasStore.Urls = pointer.GetOrNilIfEmpty(append([]string{}, store.Urls...))
		additionalStorageClasses := make([]string, 0, len(store.AdditionalStorageClasses))
		additionalStorageClasses = append(additionalStorageClasses, store.AdditionalStorageClasses...)
		atlasStore.AdditionalStorageClasses = pointer.GetOrNilIfEmpty(additionalStorageClasses)
		stores = append(stores, atlasStore)
	}
	result.Stores = pointer.GetOrNilIfEmpty(stores)
	return result
}

func readPreferenceFromAtlas(preference *admin.DataLakeAtlasStoreReadPreference) *akov2.ReadPreference {
	if preference == nil {
		return nil
	}
	result := &akov2.ReadPreference{
		MaxStalenessSeconds: preference.GetMaxStalenessSeconds(),
		Mode:                preference.GetMode(),
	}
	for _, tagset := range preference.GetTagSets() {
		var akoTags []akov2.ReadPreferenceTag
		if len(tagset) > 0 {
			akoTags = make([]akov2.ReadPreferenceTag, 0, len(tagset))
			for _, tag := range tagset {
				akoTags = append(akoTags, akov2.ReadPreferenceTag{
					Name:  tag.GetName(),
					Value: tag.GetValue(),
				})
			}
		}
		result.TagSets = append(result.TagSets, akoTags)
	}
	return result
}

func readPreferenceToAtlas(preference *akov2.ReadPreference) *admin.DataLakeAtlasStoreReadPreference {
	if preference == nil {
		return nil
	}

	var atlasTagSets [][]admin.DataLakeAtlasStoreReadPreferenceTag
	if len(preference.TagSets) > 0 {
		atlasTagSets = make([][]admin.DataLakeAtlasStoreReadPreferenceTag, 0, len(preference.TagSets))
		for _, tagset := range preference.TagSets {
			var atlasTags []admin.DataLakeAtlasStoreReadPreferenceTag
			if len(tagset) > 0 {
				atlasTags = make([]admin.DataLakeAtlasStoreReadPreferenceTag, 0, len(tagset))
				for _, tag := range tagset {
					atlasTags = append(atlasTags, admin.DataLakeAtlasStoreReadPreferenceTag{
						Name:  pointer.MakePtrOrNil(tag.Name),
						Value: pointer.MakePtrOrNil(tag.Value),
					})
				}
			}
			atlasTagSets = append(atlasTagSets, atlasTags)
		}
	}

	return &admin.DataLakeAtlasStoreReadPreference{
		MaxStalenessSeconds: pointer.MakePtrOrNil(preference.MaxStalenessSeconds),
		Mode:                pointer.MakePtrOrNil(preference.Mode),
		TagSets:             pointer.GetOrNilIfEmpty(atlasTagSets),
	}
}

func readConcernFromAtlas(concern *admin.DataLakeAtlasStoreReadConcern) *akov2.ReadConcern {
	if concern == nil {
		return nil
	}

	return &akov2.ReadConcern{
		Level: concern.GetLevel(),
	}
}

func readConcernToAtlas(concern *akov2.ReadConcern) *admin.DataLakeAtlasStoreReadConcern {
	if concern == nil {
		return nil
	}

	return &admin.DataLakeAtlasStoreReadConcern{
		Level: pointer.MakePtrOrNil(concern.Level),
	}
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
		Aws: admin.DataLakeAWSCloudProviderConfig{
			RoleId:       config.AWS.RoleID,
			TestS3Bucket: config.AWS.TestS3Bucket,
		},
	}
}
