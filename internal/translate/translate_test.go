package translate_test

import (
	"bufio"
	"embed"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin2025 "go.mongodb.org/atlas-sdk/v20250312005/admin"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/josvazg/crd2go/internal/crd2go"
	v1 "github.com/josvazg/crd2go/internal/crd2go/samples/v1"
	"github.com/josvazg/crd2go/internal/pointer"
	"github.com/josvazg/crd2go/internal/translate"
)

const (
	version = "v1"
)

//go:embed samples/*
var samples embed.FS

// NetworkPermissions is a required sturct wrapper to match the API structure
// TODO: do we need a mapping option? for this case a rename would suffice to
// load the entry array field as results in a PaginatedNetworkAccess.
// On the other hand, is extracting the whole list the proper way interact with the API?
type NetworkPermissions struct {
	Entry []admin2025.NetworkPermissionEntry `json:"entry"`
}

func TestToAPI(t *testing.T) {
	for _, tc := range []struct {
		name       string
		crd        string
		sdkVersion string
		spec       any
		deps       []client.Object
		target     any
		want       any
	}{
		{
			name:       "simple group",
			crd:        "Group",
			sdkVersion: "v20250312",
			spec: v1.GroupSpec{
				V20250312: &v1.GroupSpecV20250312{
					Entry: &v1.V20250312Entry{
						Name:                      "project-name",
						OrgId:                     "60987654321654321",
						RegionUsageRestrictions:   pointer.Get("fake-restriction"),
						WithDefaultAlertsSettings: pointer.Get(true),
						Tags: &[]v1.Tags{
							{Key: "key", Value: "value"},
						},
					},
					// read only field, not translated back to the API
					ProjectOwnerId: "61234567890123456",
				},
			},
			target: &admin2025.Group{},
			want: &admin2025.Group{
				Name:                      "project-name",
				OrgId:                     "60987654321654321",
				RegionUsageRestrictions:   pointer.Get("fake-restriction"),
				WithDefaultAlertsSettings: pointer.Get(true),
				Tags: &[]admin2025.ResourceTag{
					{Key: "key", Value: "value"},
				},
			},
		},
		{
			name:       "group alert config with project and credential references",
			crd:        "GroupAlertsConfig",
			sdkVersion: "v20250312",
			spec: v1.GroupAlertsConfigSpec{
				V20250312: &v1.GroupAlertsConfigSpecV20250312{
					Entry: &v1.GroupAlertsConfigSpecV20250312Entry{
						Enabled:       pointer.Get(true),
						EventTypeName: pointer.Get("event-type"),
						Matchers: &[]v1.Matchers{
							{
								FieldName: "field-name-1",
								Operator:  "operator-1",
								Value:     "value-1",
							},
							{
								FieldName: "field-name-2",
								Operator:  "operator-2",
								Value:     "value-2",
							},
						},
						MetricThreshold: &v1.MetricThreshold{
							MetricName: "metric-1",
							Mode:       pointer.Get("mode"),
							Operator:   pointer.Get("operator"),
							Threshold:  pointer.Get(1.1),
							Units:      pointer.Get("units"),
						},
						Threshold: &v1.MetricThreshold{
							MetricName: "metric-t",
							Mode:       pointer.Get("mode-t"),
							Operator:   pointer.Get("operator-t"),
							Threshold:  pointer.Get(2.2),
							Units:      pointer.Get("units-t"),
						},
						Notifications: &[]v1.Notifications{
							{
								DatadogApiKeySecretRef: &v1.ApiTokenSecretRef{
									Name: pointer.Get("datadog-secret"),
								},
								DatadogRegion: pointer.Get("US"),
							},
						},
						SeverityOverride: pointer.Get("some-severity-override"),
					},
					GroupId: pointer.Get("60965432187654321"),
				},
			},
			deps: []client.Object{
				&corev1.Secret{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Secret",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "datadog-secret",
						Namespace: "ns",
					},
					Data: map[string][]byte{
						"datadogApiKey": ([]byte)("sample-password"),
					},
				},
			},
			target: &admin2025.GroupAlertsConfig{},
			want: &admin2025.GroupAlertsConfig{
				Enabled:       pointer.Get(true),
				EventTypeName: pointer.Get("event-type"),
				Matchers: &[]admin2025.StreamsMatcher{
					{
						FieldName: "field-name-1",
						Operator:  "operator-1",
						Value:     "value-1",
					},
					{
						FieldName: "field-name-2",
						Operator:  "operator-2",
						Value:     "value-2",
					},
				},
				MetricThreshold: &admin2025.FlexClusterMetricThreshold{
					MetricName: "metric-1",
					Mode:       pointer.Get("mode"),
					Operator:   pointer.Get("operator"),
					Threshold:  pointer.Get(1.1),
					Units:      pointer.Get("units"),
				},
				Threshold: &admin2025.StreamProcessorMetricThreshold{
					MetricName: pointer.Get("metric-t"),
					Mode:       pointer.Get("mode-t"),
					Operator:   pointer.Get("operator-t"),
					Threshold:  pointer.Get(2.2),
					Units:      pointer.Get("units-t"),
				},
				Notifications: &[]admin2025.AlertsNotificationRootForGroup{
					{
						DatadogApiKey: pointer.Get("sample-password"),
						DatadogRegion: pointer.Get("US"),
					},
				},
				GroupId:          pointer.Get("60965432187654321"),
				SeverityOverride: pointer.Get("some-severity-override"),
			},
		},
		{
			name:       "sample organization",
			crd:        "Organization",
			sdkVersion: "v20250312",
			spec: v1.OrganizationSpec{
				V20250312: &v1.V20250312{
					Entry: &v1.Entry{
						ApiKey: &v1.ApiKey{
							Desc:  "description",
							Roles: []string{"role-1", "role-2"},
						},
						FederationSettingsId:      pointer.Get("fed-id"),
						Name:                      "org-name",
						OrgOwnerId:                pointer.Get("org-owner-id"),
						SkipDefaultAlertsSettings: pointer.Get(true),
					},
				},
			},
			target: &admin2025.AtlasOrganization{},
			want: &admin2025.AtlasOrganization{
				Name:                      "org-name",
				SkipDefaultAlertsSettings: pointer.Get(true),
			},
		},
		{
			name:       "sample network peering connection",
			crd:        "NetworkPeeringConnection",
			sdkVersion: "v20250312",
			spec: v1.NetworkPeeringConnectionSpec{
				V20250312: &v1.NetworkPeeringConnectionSpecV20250312{
					Entry: &v1.NetworkPeeringConnectionSpecV20250312Entry{
						AccepterRegionName:  pointer.Get("accepter-region-name"),
						AwsAccountId:        pointer.Get("aws-account-id"),
						AzureDirectoryId:    pointer.Get("azure-dir-id"),
						AzureSubscriptionId: pointer.Get("azure-subcription-id"),
						ContainerId:         "container-id",
						GcpProjectId:        pointer.Get("azure-subcription-id"),
						NetworkName:         pointer.Get("net-name"),
						ProviderName:        pointer.Get("provider-name"),
						ResourceGroupName:   pointer.Get("resource-group-name"),
						RouteTableCidrBlock: pointer.Get("cidr"),
						VnetName:            pointer.Get("vnet-name"),
						VpcId:               pointer.Get("vpc-id"),
					},
					GroupId: pointer.Get("32b6e34b3d91647abb20e7b8"),
				},
			},
			target: &admin2025.BaseNetworkPeeringConnectionSettings{},
			want: &admin2025.BaseNetworkPeeringConnectionSettings{
				ContainerId:         "container-id",
				ProviderName:        pointer.Get("provider-name"),
				AccepterRegionName:  pointer.Get("accepter-region-name"),
				AwsAccountId:        pointer.Get("aws-account-id"),
				RouteTableCidrBlock: pointer.Get("cidr"),
				VpcId:               pointer.Get("vpc-id"),
				AzureDirectoryId:    pointer.Get("azure-dir-id"),
				AzureSubscriptionId: pointer.Get("azure-subcription-id"),
				ResourceGroupName:   pointer.Get("resource-group-name"),
				VnetName:            pointer.Get("vnet-name"),
				GcpProjectId:        pointer.Get("azure-subcription-id"),
				NetworkName:         pointer.Get("net-name"),
			},
		},
		{
			name:       "sample database user",
			crd:        "DatabaseUser",
			sdkVersion: "v20250312",
			spec: v1.DatabaseUserSpec{
				V20250312: &v1.DatabaseUserSpecV20250312{
					Entry: &v1.DatabaseUserSpecV20250312Entry{
						Username:     "test-user",
						DatabaseName: "admin",
						GroupId:      "32b6e34b3d91647abb20e7b8",
						Roles: &[]v1.Roles{
							{DatabaseName: "admin", RoleName: "readWrite"},
						},
						AwsIAMType:      pointer.Get("aws-iam-type"),
						DeleteAfterDate: pointer.Get("2025-07-01T00:00:00Z"),
						Description:     pointer.Get("description"),
						Labels: &[]v1.Tags{
							{Key: "key-1", Value: "value-1"},
							{Key: "key-2", Value: "value-2"},
						},
						LdapAuthType: pointer.Get("ldap-auth-type"),
						OidcAuthType: pointer.Get("oidc-auth-type"),
						Password:     pointer.Get("password"),
						Scopes: &[]v1.Scopes{
							{Name: "scope-1", Type: "type-1"},
							{Name: "scope-2", Type: "type-2"},
						},
						X509Type: pointer.Get("x509-type"),
					},
					GroupId: pointer.Get("32b6e34b3d91647abb20e7b8"),
				},
			},
			target: &admin2025.CloudDatabaseUser{},
			want: &admin2025.CloudDatabaseUser{
				Username:     "test-user",
				DatabaseName: "admin",
				GroupId:      "32b6e34b3d91647abb20e7b8",
				Roles: &[]admin2025.DatabaseUserRole{
					{DatabaseName: "admin", RoleName: "readWrite"},
				},
				AwsIAMType:      pointer.Get("aws-iam-type"),
				DeleteAfterDate: pointer.Get(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)),
				Description:     pointer.Get("description"),
				Labels: &[]admin2025.ComponentLabel{
					{Key: pointer.Get("key-1"), Value: pointer.Get("value-1")},
					{Key: pointer.Get("key-2"), Value: pointer.Get("value-2")},
				},
				LdapAuthType: pointer.Get("ldap-auth-type"),
				OidcAuthType: pointer.Get("oidc-auth-type"),
				Password:     pointer.Get("password"),
				Scopes: &[]admin2025.UserScope{
					{Name: "scope-1", Type: "type-1"},
					{Name: "scope-2", Type: "type-2"},
				},
				X509Type: pointer.Get("x509-type"),
			},
		},
		{
			name:       "sample backup compliance policy",
			crd:        "BackupCompliancePolicy",
			sdkVersion: "v20250312",
			spec: v1.BackupCompliancePolicySpec{
				V20250312: &v1.BackupCompliancePolicySpecV20250312{
					Entry: &v1.BackupCompliancePolicySpecV20250312Entry{
						AuthorizedEmail:         "user@example.com",
						CopyProtectionEnabled:   pointer.Get(true),
						EncryptionAtRestEnabled: pointer.Get(true),
						AuthorizedUserFirstName: "first-name",
						AuthorizedUserLastName:  "last-name",
						OnDemandPolicyItem: &v1.OnDemandPolicyItem{
							FrequencyInterval: 1,
							FrequencyType:     "some-freq",
							RetentionUnit:     "some-unit",
							RetentionValue:    2,
						},
						PitEnabled:        pointer.Get(true),
						ProjectId:         pointer.Get("project-id"),
						RestoreWindowDays: pointer.Get(3),
						ScheduledPolicyItems: &[]v1.OnDemandPolicyItem{
							{
								FrequencyInterval: 3,
								FrequencyType:     "some-freq",
								RetentionUnit:     "some-unit",
								RetentionValue:    4,
							},
							{
								FrequencyInterval: 5,
								FrequencyType:     "some-other-freq",
								RetentionUnit:     "some-other-unit",
								RetentionValue:    6,
							},
						},
					},
					GroupId:                 pointer.Get("32b6e34b3d91647abb20e7b8"),
					OverwriteBackupPolicies: true,
				},
			},
			target: &admin2025.DataProtectionSettings20231001{},
			want: &admin2025.DataProtectionSettings20231001{
				AuthorizedEmail:         "user@example.com",
				CopyProtectionEnabled:   pointer.Get(true),
				EncryptionAtRestEnabled: pointer.Get(true),
				AuthorizedUserFirstName: "first-name",
				AuthorizedUserLastName:  "last-name",
				OnDemandPolicyItem: &admin2025.BackupComplianceOnDemandPolicyItem{
					FrequencyInterval: 1,
					FrequencyType:     "some-freq",
					RetentionUnit:     "some-unit",
					RetentionValue:    2,
				},
				PitEnabled:        pointer.Get(true),
				ProjectId:         pointer.Get("project-id"),
				RestoreWindowDays: pointer.Get(3),
				ScheduledPolicyItems: &[]admin2025.BackupComplianceScheduledPolicyItem{
					{
						FrequencyInterval: 3,
						FrequencyType:     "some-freq",
						RetentionUnit:     "some-unit",
						RetentionValue:    4,
					},
					{
						FrequencyInterval: 5,
						FrequencyType:     "some-other-freq",
						RetentionUnit:     "some-other-unit",
						RetentionValue:    6,
					},
				},
			},
		},
		{
			name:       "backup schedule all fields",
			crd:        "BackupSchedule",
			sdkVersion: "v20250312",
			spec: v1.BackupScheduleSpec{
				V20250312: &v1.BackupScheduleSpecV20250312{
					Entry: &v1.BackupScheduleSpecV20250312Entry{
						ReferenceHourOfDay:    pointer.Get(2),
						ReferenceMinuteOfHour: pointer.Get(30),
						RestoreWindowDays:     pointer.Get(7),
						UpdateSnapshots:       pointer.Get(true),
						AutoExportEnabled:     pointer.Get(true),
						CopySettings: &[]v1.CopySettings{
							{
								CloudProvider:    pointer.Get("AWS"),
								Frequencies:      &[]string{"freq-1", "freq-2"},
								RegionName:       pointer.Get("us-east-1"),
								ShouldCopyOplogs: pointer.Get(true),
								ZoneId:           "zone-id",
							},
							{
								CloudProvider:    pointer.Get("GCE"),
								Frequencies:      &[]string{"freq-3", "freq-4"},
								RegionName:       pointer.Get("us-east-3"),
								ShouldCopyOplogs: pointer.Get(true),
								ZoneId:           "zone-id-0",
							},
						},
						DeleteCopiedBackups: &[]v1.DeleteCopiedBackups{
							{
								CloudProvider: pointer.Get("Azure"),
								RegionName:    pointer.Get("us-west-2"),
								ZoneId:        pointer.Get("zone-id"),
							},
						},
						Export: &v1.Export{
							ExportBucketId: pointer.Get("ExportBucketId"),
							FrequencyType:  pointer.Get("FrequencyType"),
						},
						ExtraRetentionSettings: &[]v1.ExtraRetentionSettings{
							{
								FrequencyType: pointer.Get("FrequencyType0"),
								RetentionDays: pointer.Get(1),
							},
							{
								FrequencyType: pointer.Get("FrequencyType1"),
								RetentionDays: pointer.Get(2),
							},
						},
						Policies: &[]v1.Policies{
							{
								Id: pointer.Get("id0"),
								PolicyItems: &[]v1.OnDemandPolicyItem{
									{
										FrequencyInterval: 1,
										FrequencyType:     "freq-type0",
										RetentionUnit:     "ret-unit0",
										RetentionValue:    2,
									},
									{
										FrequencyInterval: 3,
										FrequencyType:     "freq-type1",
										RetentionUnit:     "ret-unit1",
										RetentionValue:    4,
									},
								},
							},
						},
						UseOrgAndGroupNamesInExportPrefix: pointer.Get(true),
					},
					GroupId:     pointer.Get("group-id-101"),
					ClusterName: "cluster-name",
				},
			},
			target: &admin2025.DiskBackupSnapshotSchedule20240805{},
			want: &admin2025.DiskBackupSnapshotSchedule20240805{
				ReferenceHourOfDay:    pointer.Get(2),
				ReferenceMinuteOfHour: pointer.Get(30),
				RestoreWindowDays:     pointer.Get(7),
				UpdateSnapshots:       pointer.Get(true),
				AutoExportEnabled:     pointer.Get(true),
				CopySettings: &[]admin2025.DiskBackupCopySetting20240805{
					{
						CloudProvider:    pointer.Get("AWS"),
						Frequencies:      &[]string{"freq-1", "freq-2"},
						RegionName:       pointer.Get("us-east-1"),
						ShouldCopyOplogs: pointer.Get(true),
						ZoneId:           "zone-id",
					},
					{
						CloudProvider:    pointer.Get("GCE"),
						Frequencies:      &[]string{"freq-3", "freq-4"},
						RegionName:       pointer.Get("us-east-3"),
						ShouldCopyOplogs: pointer.Get(true),
						ZoneId:           "zone-id-0",
					},
				},
				DeleteCopiedBackups: &[]admin2025.DeleteCopiedBackups20240805{
					{
						CloudProvider: pointer.Get("Azure"),
						RegionName:    pointer.Get("us-west-2"),
						ZoneId:        pointer.Get("zone-id"),
					},
				},
				Export: &admin2025.AutoExportPolicy{
					ExportBucketId: pointer.Get("ExportBucketId"),
					FrequencyType:  pointer.Get("FrequencyType"),
				},
				ExtraRetentionSettings: &[]admin2025.ExtraRetentionSetting{
					{
						FrequencyType: pointer.Get("FrequencyType0"),
						RetentionDays: pointer.Get(1),
					},
					{
						FrequencyType: pointer.Get("FrequencyType1"),
						RetentionDays: pointer.Get(2),
					},
				},
				Policies: &[]admin2025.AdvancedDiskBackupSnapshotSchedulePolicy{
					{
						Id: pointer.Get("id0"),
						PolicyItems: &[]admin2025.DiskBackupApiPolicyItem{
							{
								FrequencyInterval: 1,
								FrequencyType:     "freq-type0",
								RetentionUnit:     "ret-unit0",
								RetentionValue:    2,
							},
							{
								FrequencyInterval: 3,
								FrequencyType:     "freq-type1",
								RetentionUnit:     "ret-unit1",
								RetentionValue:    4,
							},
						},
					},
				},
				UseOrgAndGroupNamesInExportPrefix: pointer.Get(true),
				ClusterName:                       pointer.Get("cluster-name"),
			},
		},
		{
			name:       "data federation all fields",
			crd:        "DataFederation",
			sdkVersion: "v20250312",
			spec: v1.DataFederationSpec{
				V20250312: &v1.DataFederationSpecV20250312{
					Entry: &v1.DataFederationSpecV20250312Entry{
						CloudProviderConfig: &v1.CloudProviderConfig{
							Aws: &v1.Aws{
								RoleId:       "aws-role-id-123",
								TestS3Bucket: "my-s3-bucket",
							},
							Azure: &v1.Azure{
								RoleId: "azure-role-id-456",
							},
							Gcp: &v1.Azure{
								RoleId: "gcp-role-id-789",
							},
						},
						DataProcessRegion: &v1.DataProcessRegion{
							CloudProvider: "GCE",
							Region:        "eu-north-2",
						},
						Name: pointer.Get("some-name"),
						Storage: &v1.Storage{
							Databases: &[]v1.Databases{
								{
									Collections: &[]v1.Collections{
										{
											DataSources: &[]v1.DataSources{
												{
													AllowInsecure:       pointer.Get(true),
													Collection:          pointer.Get("some-name"),
													CollectionRegex:     pointer.Get("collection-regex"),
													Database:            pointer.Get("db"),
													DatabaseRegex:       pointer.Get("db-regex"),
													DatasetName:         pointer.Get("dataset-name"),
													DatasetPrefix:       pointer.Get("dataset-prefix"),
													DefaultFormat:       pointer.Get("default-format"),
													Path:                pointer.Get("path"),
													ProvenanceFieldName: pointer.Get("provenqance-field-name"),
													StoreName:           pointer.Get("store-name"),
													TrimLevel:           pointer.Get(1),
													Urls:                &[]string{"url1", "url2"},
												},
											},
											Name: pointer.Get("collection0"),
										},
									},
									MaxWildcardCollections: pointer.Get(3),
									Name:                   pointer.Get("db0"),
									Views: &[]v1.Views{
										{
											Name:     pointer.Get("view0"),
											Pipeline: pointer.Get("pipeline0"),
											Source:   pointer.Get("source0"),
										},
									},
								},
							},
							Stores: &[]v1.Stores{
								{
									AdditionalStorageClasses: &[]string{"stc1", "stc2"},
									AllowInsecure:            pointer.Get(true),
									Bucket:                   pointer.Get("bucket-name"),
									ClusterName:              pointer.Get("cluster-name"),
									ContainerName:            pointer.Get("container-name"),
									DefaultFormat:            pointer.Get("default-format"),
									Delimiter:                pointer.Get("delimiter"),
									IncludeTags:              pointer.Get(true),
									Name:                     pointer.Get("store-name"),
									Prefix:                   pointer.Get("prefix"),
									Provider:                 "AWS",
									Public:                   pointer.Get(true),
									ReadConcern: &v1.ReadConcern{
										Level: pointer.Get("local"),
									},
									ReadPreference: &v1.ReadPreference{
										Mode: pointer.Get("primary"),
									},
									Region:               pointer.Get("us-east-1"),
									ReplacementDelimiter: pointer.Get("replacement-delimiter"),
									ServiceURL:           pointer.Get("https://service-url.com"),
									Urls:                 &[]string{"url1", "url2"},
								},
							},
						},
					},
				},
			},
			target: &admin2025.DataLakeTenant{},
			want: &admin2025.DataLakeTenant{
				Name: pointer.Get("some-name"),
				CloudProviderConfig: &admin2025.DataLakeCloudProviderConfig{
					Aws: &admin2025.DataLakeAWSCloudProviderConfig{
						RoleId:       "aws-role-id-123",
						TestS3Bucket: "my-s3-bucket",
					},
					Azure: &admin2025.DataFederationAzureCloudProviderConfig{
						RoleId: "azure-role-id-456",
					},
					Gcp: &admin2025.DataFederationGCPCloudProviderConfig{
						RoleId: "gcp-role-id-789",
					},
				},
				DataProcessRegion: &admin2025.DataLakeDataProcessRegion{
					CloudProvider: "GCE",
					Region:        "eu-north-2",
				},
				Storage: &admin2025.DataLakeStorage{
					Databases: &[]admin2025.DataLakeDatabaseInstance{
						{
							Collections: &[]admin2025.DataLakeDatabaseCollection{
								{
									DataSources: &[]admin2025.DataLakeDatabaseDataSourceSettings{
										{
											AllowInsecure:       pointer.Get(true),
											Collection:          pointer.Get("some-name"),
											CollectionRegex:     pointer.Get("collection-regex"),
											Database:            pointer.Get("db"),
											DatabaseRegex:       pointer.Get("db-regex"),
											DatasetName:         pointer.Get("dataset-name"),
											DatasetPrefix:       pointer.Get("dataset-prefix"),
											DefaultFormat:       pointer.Get("default-format"),
											Path:                pointer.Get("path"),
											ProvenanceFieldName: pointer.Get("provenqance-field-name"),
											StoreName:           pointer.Get("store-name"),
											TrimLevel:           pointer.Get(1),
											Urls:                &[]string{"url1", "url2"},
										},
									},
									Name: pointer.Get("collection0"),
								},
							},
							MaxWildcardCollections: pointer.Get(3),
							Name:                   pointer.Get("db0"),
							Views: &[]admin2025.DataLakeApiBase{
								{
									Name:     pointer.Get("view0"),
									Pipeline: pointer.Get("pipeline0"),
									Source:   pointer.Get("source0"),
								},
							},
						},
					},
					Stores: &[]admin2025.DataLakeStoreSettings{
						{
							AdditionalStorageClasses: &[]string{"stc1", "stc2"},
							AllowInsecure:            pointer.Get(true),
							Bucket:                   pointer.Get("bucket-name"),
							ClusterName:              pointer.Get("cluster-name"),
							ContainerName:            pointer.Get("container-name"),
							DefaultFormat:            pointer.Get("default-format"),
							Delimiter:                pointer.Get("delimiter"),
							IncludeTags:              pointer.Get(true),
							Name:                     pointer.Get("store-name"),
							Prefix:                   pointer.Get("prefix"),
							Provider:                 "AWS",
							Public:                   pointer.Get(true),
							ReadConcern: &admin2025.DataLakeAtlasStoreReadConcern{
								Level: pointer.Get("local"),
							},
							ReadPreference: &admin2025.DataLakeAtlasStoreReadPreference{
								Mode: pointer.Get("primary"),
							},
							Region:               pointer.Get("us-east-1"),
							ReplacementDelimiter: pointer.Get("replacement-delimiter"),
							ServiceURL:           pointer.Get("https://service-url.com"),
							Urls:                 &[]string{"url1", "url2"},
						},
					},
				},
			},
		},
		{
			name:       "Organization setting with all fields",
			crd:        "OrganizationSetting",
			sdkVersion: "v20250312",
			spec: v1.OrganizationSettingSpec{
				V20250312: &v1.OrganizationSettingSpecV20250312{
					Entry: &v1.OrganizationSettingSpecV20250312Entry{
						ApiAccessListRequired:                  pointer.Get(true),
						GenAIFeaturesEnabled:                   pointer.Get(true),
						MaxServiceAccountSecretValidityInHours: pointer.Get(24),
						MultiFactorAuthRequired:                pointer.Get(true),
						RestrictEmployeeAccess:                 pointer.Get(true),
						SecurityContact:                        pointer.Get("contact-info"),
						StreamsCrossGroupEnabled:               pointer.Get(true),
					},
					OrgId: "org-id",
				},
			},
			target: &admin2025.OrganizationSettings{},
			want: &admin2025.OrganizationSettings{
				ApiAccessListRequired:                  pointer.Get(true),
				GenAIFeaturesEnabled:                   pointer.Get(true),
				MaxServiceAccountSecretValidityInHours: pointer.Get(24),
				MultiFactorAuthRequired:                pointer.Get(true),
				RestrictEmployeeAccess:                 pointer.Get(true),
				SecurityContact:                        pointer.Get("contact-info"),
				StreamsCrossGroupEnabled:               pointer.Get(true),
			},
		},
		{
			name:       "team all fields",
			crd:        "Team",
			sdkVersion: "v20250312",
			spec: v1.TeamSpec{
				V20250312: &v1.TeamSpecV20250312{
					Entry: &v1.TeamSpecV20250312Entry{
						Name:      "team-name",
						Usernames: []string{"user1", "user2"},
					},
					OrgId: "org-id",
				},
			},
			target: &admin2025.Team{},
			want: &admin2025.Team{
				Name: "team-name",
				Usernames: []string{
					"user1", "user2",
				},
			},
		},
		{
			name:       "network permission entries all fields",
			crd:        "NetworkPermissionEntries",
			sdkVersion: "v20250312",
			spec: v1.NetworkPermissionEntriesSpec{
				V20250312: &v1.NetworkPermissionEntriesSpecV20250312{
					Entry: &[]v1.NetworkPermissionEntriesSpecV20250312Entry{
						{
							AwsSecurityGroup: pointer.Get("sg-12345678"),
							CidrBlock:        pointer.Get("cird"),
							Comment:          pointer.Get("comment"),
							DeleteAfterDate:  pointer.Get("2025-07-01T00:00:00Z"),
							IpAddress:        pointer.Get("1.1.1.1"),
						},
					},
					GroupId: pointer.Get("32b6e34b3d91647abb20e7b8"),
				},
			},
			target: &NetworkPermissions{},
			want: &NetworkPermissions{
				Entry: []admin2025.NetworkPermissionEntry{
					{
						AwsSecurityGroup: pointer.Get("sg-12345678"),
						CidrBlock:        pointer.Get("cird"),
						Comment:          pointer.Get("comment"),
						DeleteAfterDate:  pointer.Get(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)),
						IpAddress:        pointer.Get("1.1.1.1"),
					},
				},
			},
		},
		{
			name:       "cluster all fields",
			crd:        "Cluster",
			sdkVersion: "v20250312",
			spec: v1.ClusterSpec{
				V20250312: &v1.ClusterSpecV20250312{
					Entry: &v1.ClusterSpecV20250312Entry{
						AcceptDataRisksAndForceReplicaSetReconfig: pointer.Get("2025-01-01T00:00:00Z"),
						AdvancedConfiguration: &v1.AdvancedConfiguration{
							CustomOpensslCipherConfigTls12: &[]string{
								"TLS_AES_256_GCM_SHA384", "TLS_CHACHA20_POLY1305_SHA256",
							},
							MinimumEnabledTlsProtocol: pointer.Get("TLS1.2"),
							TlsCipherConfigMode:       pointer.Get("Custom"),
						},
						BackupEnabled:                             pointer.Get(true),
						BiConnector:                               &v1.BiConnector{Enabled: pointer.Get(true)},
						ClusterType:                               pointer.Get("ReplicaSet"),
						ConfigServerManagementMode:                pointer.Get("Managed"),
						ConfigServerType:                          pointer.Get("ReplicaSet"),
						DiskWarmingMode:                           pointer.Get("Enabled"),
						EncryptionAtRestProvider:                  pointer.Get("AWS-KMS"),
						FeatureCompatibilityVersion:               pointer.Get("7.0"),
						FeatureCompatibilityVersionExpirationDate: pointer.Get("2025-12-31T00:00:00Z"),
						GlobalClusterSelfManagedSharding:          pointer.Get(true),
						Labels: &[]v1.Tags{
							{Key: "key1", Value: "value1"},
							{Key: "key2", Value: "value2"},
						},
						MongoDBEmployeeAccessGrant: &v1.MongoDBEmployeeAccessGrant{
							ExpirationTime: "2025-12-31T00:00:00Z",
							GrantType:      "Temporary",
						},
						MongoDBMajorVersion:       pointer.Get("8.0"),
						Name:                      pointer.Get("my-cluster"),
						Paused:                    pointer.Get(true),
						PitEnabled:                pointer.Get(true),
						RedactClientLogData:       pointer.Get(true),
						ReplicaSetScalingStrategy: pointer.Get("Auto"),
						ReplicationSpecs: &[]v1.ReplicationSpecs{
							{
								ZoneId:   pointer.Get("zone-id-1"),
								ZoneName: pointer.Get("zone-name-1"),
								RegionConfigs: &[]v1.RegionConfigs{
									{
										RegionName: pointer.Get("us-east-1"),
										AnalyticsSpecs: &v1.AnalyticsSpecs{
											DiskIOPS:      pointer.Get(1000),
											DiskSizeGB:    pointer.Get(10.0),
											EbsVolumeType: pointer.Get("gp2"),
											InstanceSize:  pointer.Get("M10"),
											NodeCount:     pointer.Get(3),
										},
										AutoScaling: &v1.AnalyticsAutoScaling{
											Compute: &v1.Compute{
												Enabled:           pointer.Get(true),
												ScaleDownEnabled:  pointer.Get(true),
												MaxInstanceSize:   pointer.Get("M20"),
												MinInstanceSize:   pointer.Get("M10"),
												PredictiveEnabled: pointer.Get(true),
											},
											DiskGB: &v1.DiskGB{
												Enabled: pointer.Get(true),
											},
										},
										AnalyticsAutoScaling: &v1.AnalyticsAutoScaling{
											Compute: &v1.Compute{
												Enabled:           pointer.Get(true),
												ScaleDownEnabled:  pointer.Get(true),
												MaxInstanceSize:   pointer.Get("M30"),
												MinInstanceSize:   pointer.Get("M10"),
												PredictiveEnabled: pointer.Get(true),
											},
											DiskGB: &v1.DiskGB{
												Enabled: pointer.Get(true),
											},
										},
										BackingProviderName: pointer.Get("AWS"),
										ElectableSpecs: &v1.ElectableSpecs{
											DiskIOPS:              pointer.Get(1000),
											DiskSizeGB:            pointer.Get(10.0),
											EbsVolumeType:         pointer.Get("gp2"),
											EffectiveInstanceSize: pointer.Get("M10"),
											InstanceSize:          pointer.Get("M10"),
											NodeCount:             pointer.Get(3),
										},
										Priority:     pointer.Get(1),
										ProviderName: pointer.Get("AWS"),
										ReadOnlySpecs: &v1.AnalyticsSpecs{
											DiskIOPS:      pointer.Get(1000),
											DiskSizeGB:    pointer.Get(10.0),
											EbsVolumeType: pointer.Get("gp2"),
											InstanceSize:  pointer.Get("M10"),
											NodeCount:     pointer.Get(3),
										},
									},
									{
										RegionName: pointer.Get("us-east-2"),
										AnalyticsSpecs: &v1.AnalyticsSpecs{
											DiskIOPS:      pointer.Get(2000),
											DiskSizeGB:    pointer.Get(10.0),
											EbsVolumeType: pointer.Get("gp3"),
											InstanceSize:  pointer.Get("M20"),
											NodeCount:     pointer.Get(3),
										},
										AutoScaling: &v1.AnalyticsAutoScaling{
											Compute: &v1.Compute{
												Enabled:           pointer.Get(true),
												ScaleDownEnabled:  pointer.Get(true),
												MaxInstanceSize:   pointer.Get("M50"),
												MinInstanceSize:   pointer.Get("M20"),
												PredictiveEnabled: pointer.Get(true),
											},
											DiskGB: &v1.DiskGB{
												Enabled: pointer.Get(true),
											},
										},
										AnalyticsAutoScaling: &v1.AnalyticsAutoScaling{
											Compute: &v1.Compute{
												Enabled:           pointer.Get(true),
												ScaleDownEnabled:  pointer.Get(true),
												MaxInstanceSize:   pointer.Get("M40"),
												MinInstanceSize:   pointer.Get("M10"),
												PredictiveEnabled: pointer.Get(true),
											},
											DiskGB: &v1.DiskGB{
												Enabled: pointer.Get(true),
											},
										},
										BackingProviderName: pointer.Get("AWS"),
										ElectableSpecs: &v1.ElectableSpecs{
											DiskIOPS:              pointer.Get(1000),
											DiskSizeGB:            pointer.Get(10.0),
											EbsVolumeType:         pointer.Get("gp2"),
											EffectiveInstanceSize: pointer.Get("M10"),
											InstanceSize:          pointer.Get("M10"),
											NodeCount:             pointer.Get(3),
										},
										Priority:     pointer.Get(1),
										ProviderName: pointer.Get("AWS"),
										ReadOnlySpecs: &v1.AnalyticsSpecs{
											DiskIOPS:      pointer.Get(1000),
											DiskSizeGB:    pointer.Get(10.0),
											EbsVolumeType: pointer.Get("gp2"),
											InstanceSize:  pointer.Get("M10"),
											NodeCount:     pointer.Get(3),
										},
									},
								},
							},
						},
						RootCertType: pointer.Get("X509"),
						Tags: &[]v1.Tags{
							{Key: "key1", Value: "value1"},
							{Key: "key2", Value: "value2"},
						},
						TerminationProtectionEnabled: pointer.Get(true),
						VersionReleaseSystem:         pointer.Get("Atlas"),
					},
					GroupId: pointer.Get("32b6e34b3d91647abb20e7b8"),
				},
			},
			target: &admin2025.ClusterDescription20240805{},
			want: &admin2025.ClusterDescription20240805{
				AcceptDataRisksAndForceReplicaSetReconfig: pointer.Get(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				AdvancedConfiguration: &admin2025.ApiAtlasClusterAdvancedConfiguration{
					CustomOpensslCipherConfigTls12: &[]string{
						"TLS_AES_256_GCM_SHA384", "TLS_CHACHA20_POLY1305_SHA256",
					},
					MinimumEnabledTlsProtocol: pointer.Get("TLS1.2"),
					TlsCipherConfigMode:       pointer.Get("Custom"),
				},
				BackupEnabled:                             pointer.Get(true),
				BiConnector:                               &admin2025.BiConnector{Enabled: pointer.Get(true)},
				ClusterType:                               pointer.Get("ReplicaSet"),
				ConfigServerManagementMode:                pointer.Get("Managed"),
				ConfigServerType:                          pointer.Get("ReplicaSet"),
				DiskWarmingMode:                           pointer.Get("Enabled"),
				EncryptionAtRestProvider:                  pointer.Get("AWS-KMS"),
				FeatureCompatibilityVersion:               pointer.Get("7.0"),
				FeatureCompatibilityVersionExpirationDate: pointer.Get(time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)),
				GlobalClusterSelfManagedSharding:          pointer.Get(true),
				Labels: &[]admin2025.ComponentLabel{
					{Key: pointer.Get("key1"), Value: pointer.Get("value1")},
					{Key: pointer.Get("key2"), Value: pointer.Get("value2")},
				},
				MongoDBEmployeeAccessGrant: &admin2025.EmployeeAccessGrant{
					ExpirationTime: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
					GrantType:      "Temporary",
				},
				MongoDBMajorVersion:       pointer.Get("8.0"),
				Name:                      pointer.Get("my-cluster"),
				Paused:                    pointer.Get(true),
				PitEnabled:                pointer.Get(true),
				RedactClientLogData:       pointer.Get(true),
				ReplicaSetScalingStrategy: pointer.Get("Auto"),
				ReplicationSpecs: &[]admin2025.ReplicationSpec20240805{
					{
						ZoneId:   pointer.Get("zone-id-1"),
						ZoneName: pointer.Get("zone-name-1"),
						RegionConfigs: &[]admin2025.CloudRegionConfig20240805{
							{
								RegionName: pointer.Get("us-east-1"),
								AnalyticsSpecs: &admin2025.DedicatedHardwareSpec20240805{
									DiskIOPS:      pointer.Get(1000),
									DiskSizeGB:    pointer.Get(10.0),
									EbsVolumeType: pointer.Get("gp2"),
									InstanceSize:  pointer.Get("M10"),
									NodeCount:     pointer.Get(3),
								},
								AutoScaling: &admin2025.AdvancedAutoScalingSettings{
									Compute: &admin2025.AdvancedComputeAutoScaling{
										Enabled:           pointer.Get(true),
										ScaleDownEnabled:  pointer.Get(true),
										MaxInstanceSize:   pointer.Get("M20"),
										MinInstanceSize:   pointer.Get("M10"),
										PredictiveEnabled: pointer.Get(true),
									},
									DiskGB: &admin2025.DiskGBAutoScaling{
										Enabled: pointer.Get(true),
									},
								},
								AnalyticsAutoScaling: &admin2025.AdvancedAutoScalingSettings{
									Compute: &admin2025.AdvancedComputeAutoScaling{
										Enabled:           pointer.Get(true),
										ScaleDownEnabled:  pointer.Get(true),
										MaxInstanceSize:   pointer.Get("M30"),
										MinInstanceSize:   pointer.Get("M10"),
										PredictiveEnabled: pointer.Get(true),
									},
									DiskGB: &admin2025.DiskGBAutoScaling{
										Enabled: pointer.Get(true),
									},
								},
								BackingProviderName: pointer.Get("AWS"),
								ElectableSpecs: &admin2025.HardwareSpec20240805{
									DiskIOPS:              pointer.Get(1000),
									DiskSizeGB:            pointer.Get(10.0),
									EbsVolumeType:         pointer.Get("gp2"),
									EffectiveInstanceSize: pointer.Get("M10"),
									InstanceSize:          pointer.Get("M10"),
									NodeCount:             pointer.Get(3),
								},
								Priority:     pointer.Get(1),
								ProviderName: pointer.Get("AWS"),
								ReadOnlySpecs: &admin2025.DedicatedHardwareSpec20240805{
									DiskIOPS:      pointer.Get(1000),
									DiskSizeGB:    pointer.Get(10.0),
									EbsVolumeType: pointer.Get("gp2"),
									InstanceSize:  pointer.Get("M10"),
									NodeCount:     pointer.Get(3),
								},
							},
							{
								RegionName: pointer.Get("us-east-2"),
								AnalyticsSpecs: &admin2025.DedicatedHardwareSpec20240805{
									DiskIOPS:      pointer.Get(2000),
									DiskSizeGB:    pointer.Get(10.0),
									EbsVolumeType: pointer.Get("gp3"),
									InstanceSize:  pointer.Get("M20"),
									NodeCount:     pointer.Get(3),
								},
								AutoScaling: &admin2025.AdvancedAutoScalingSettings{
									Compute: &admin2025.AdvancedComputeAutoScaling{
										Enabled:           pointer.Get(true),
										ScaleDownEnabled:  pointer.Get(true),
										MaxInstanceSize:   pointer.Get("M50"),
										MinInstanceSize:   pointer.Get("M20"),
										PredictiveEnabled: pointer.Get(true),
									},
									DiskGB: &admin2025.DiskGBAutoScaling{
										Enabled: pointer.Get(true),
									},
								},
								AnalyticsAutoScaling: &admin2025.AdvancedAutoScalingSettings{
									Compute: &admin2025.AdvancedComputeAutoScaling{
										Enabled:           pointer.Get(true),
										ScaleDownEnabled:  pointer.Get(true),
										MaxInstanceSize:   pointer.Get("M40"),
										MinInstanceSize:   pointer.Get("M10"),
										PredictiveEnabled: pointer.Get(true),
									},
									DiskGB: &admin2025.DiskGBAutoScaling{
										Enabled: pointer.Get(true),
									},
								},
								BackingProviderName: pointer.Get("AWS"),
								ElectableSpecs: &admin2025.HardwareSpec20240805{
									DiskIOPS:              pointer.Get(1000),
									DiskSizeGB:            pointer.Get(10.0),
									EbsVolumeType:         pointer.Get("gp2"),
									EffectiveInstanceSize: pointer.Get("M10"),
									InstanceSize:          pointer.Get("M10"),
									NodeCount:             pointer.Get(3),
								},
								Priority:     pointer.Get(1),
								ProviderName: pointer.Get("AWS"),
								ReadOnlySpecs: &admin2025.DedicatedHardwareSpec20240805{
									DiskIOPS:      pointer.Get(1000),
									DiskSizeGB:    pointer.Get(10.0),
									EbsVolumeType: pointer.Get("gp2"),
									InstanceSize:  pointer.Get("M10"),
									NodeCount:     pointer.Get(3),
								},
							},
						},
					},
				},
				RootCertType: pointer.Get("X509"),
				Tags: &[]admin2025.ResourceTag{
					{Key: "key1", Value: "value1"},
					{Key: "key2", Value: "value2"},
				},
				TerminationProtectionEnabled: pointer.Get(true),
				VersionReleaseSystem:         pointer.Get("Atlas"),
				GroupId:                      pointer.Get("32b6e34b3d91647abb20e7b8"),
			},
		},
		{
			name:       "searchindex definition fields",
			crd:        "SearchIndex",
			sdkVersion: "v20250312",
			spec: v1.SearchIndexSpec{
				V20250312: &v1.SearchIndexSpecV20250312{
					Entry: &v1.SearchIndexSpecV20250312Entry{
						Analyzer: pointer.Get("lucene.standard"),
						Database: "database-name",
						Analyzers: &[]v1.Analyzers{
							{
								Name: "custom-analyzer",
								CharFilters: &[]apiextensionsv1.JSON{
									{Raw: []byte(`{"key":"value"}`)},
									{Raw: []byte(`{"key2":"value2"}`)},
								},
								TokenFilters: &[]apiextensionsv1.JSON{
									{Raw: []byte(`{"key3":"value3"}`)},
									{Raw: []byte(`{"key4":"value4"}`)},
								},
								Tokenizer: v1.Tokenizer{
									Group:          pointer.Get(2),
									MaxGram:        pointer.Get(100),
									MaxTokenLength: pointer.Get(50),
									MinGram:        pointer.Get(1),
									Pattern:        pointer.Get("pattern"),
									Type:           pointer.Get("custom"),
								},
							},
						},
						CollectionName: "collection-name",
						Fields: &[]apiextensionsv1.JSON{
							{Raw: []byte(`{"field1":"value1"}`)},
							{Raw: []byte(`{"field2":"value2"}`)},
							{Raw: []byte(`{"field3":"value3"}`)},
						},
						Mappings: &v1.Mappings{
							Dynamic: pointer.Get(true),
							Fields: &apiextensionsv1.JSON{
								Raw: []byte(`{"field4":"value4"}`),
							},
						},
						Name:           "index-name",
						NumPartitions:  pointer.Get(3),
						SearchAnalyzer: pointer.Get("lucene.standard"),
						StoredSource: &apiextensionsv1.JSON{
							Raw: []byte(`{"enabled": true}`),
						},
						Synonyms: &[]v1.Synonyms{
							{
								Analyzer: "synonym-analyzer",
								Name:     "synonym-name",
								Source: v1.Source{
									Collection: "synonym-collection",
								},
							},
						},
						Type: pointer.Get("search-index-type"),
					},
					GroupId: pointer.Get("group-id-101"),
				},
			},
			target: &admin2025.BaseSearchIndexCreateRequestDefinition{},
			want: &admin2025.BaseSearchIndexCreateRequestDefinition{
				Analyzer: pointer.Get("lucene.standard"),
				Analyzers: &[]admin2025.AtlasSearchAnalyzer{
					{
						Name: "custom-analyzer",
						CharFilters: &[]any{
							map[string]any{"key": "value"},
							map[string]any{"key2": "value2"},
						},
						TokenFilters: &[]any{
							map[string]any{"key3": "value3"},
							map[string]any{"key4": "value4"},
						},
						Tokenizer: map[string]any{
							"group":          2.0,
							"maxGram":        100.0,
							"maxTokenLength": 50.0,
							"minGram":        1.0,
							"pattern":        "pattern",
							"type":           "custom",
						},
					},
				},
				Fields: &[]any{
					map[string]any{"field1": "value1"},
					map[string]any{"field2": "value2"},
					map[string]any{"field3": "value3"},
				},
				Mappings: &admin2025.SearchMappings{
					Dynamic: pointer.Get(true),
					Fields: &map[string]any{
						"field4": "value4",
					},
				},
				NumPartitions:  pointer.Get(3),
				SearchAnalyzer: pointer.Get("lucene.standard"),
				StoredSource:   any(map[string]any{"enabled": true}),
				Synonyms: &[]admin2025.SearchSynonymMappingDefinition{
					{
						Analyzer: "synonym-analyzer",
						Name:     "synonym-name",
						Source: admin2025.SynonymSource{
							Collection: "synonym-collection",
						},
					},
				},
			},
		},
		{
			name:       "searchindex create request fields",
			crd:        "SearchIndex",
			sdkVersion: "v20250312",
			spec: v1.SearchIndexSpec{
				V20250312: &v1.SearchIndexSpecV20250312{
					Entry: &v1.SearchIndexSpecV20250312Entry{
						Analyzer: pointer.Get("lucene.standard"),
						Database: "database-name",
						Analyzers: &[]v1.Analyzers{
							{
								Name: "custom-analyzer",
								CharFilters: &[]apiextensionsv1.JSON{
									{Raw: []byte(`{"key":"value"}`)},
									{Raw: []byte(`{"key2":"value2"}`)},
								},
								TokenFilters: &[]apiextensionsv1.JSON{
									{Raw: []byte(`{"key3":"value3"}`)},
									{Raw: []byte(`{"key4":"value4"}`)},
								},
								Tokenizer: v1.Tokenizer{
									Group:          pointer.Get(2),
									MaxGram:        pointer.Get(100),
									MaxTokenLength: pointer.Get(50),
									MinGram:        pointer.Get(1),
									Pattern:        pointer.Get("pattern"),
									Type:           pointer.Get("custom"),
								},
							},
						},
						CollectionName: "collection-name",
						Fields: &[]apiextensionsv1.JSON{
							{Raw: []byte(`{"field1":"value1"}`)},
							{Raw: []byte(`{"field2":"value2"}`)},
							{Raw: []byte(`{"field3":"value3"}`)},
						},
						Mappings: &v1.Mappings{
							Dynamic: pointer.Get(true),
							Fields: &apiextensionsv1.JSON{
								Raw: []byte(`{"field4":"value4"}`),
							},
						},
						Name:           "index-name",
						NumPartitions:  pointer.Get(3),
						SearchAnalyzer: pointer.Get("lucene.standard"),
						StoredSource: &apiextensionsv1.JSON{
							Raw: []byte(`{"enabled": true}`),
						},
						Synonyms: &[]v1.Synonyms{
							{
								Analyzer: "synonym-analyzer",
								Name:     "synonym-name",
								Source: v1.Source{
									Collection: "synonym-collection",
								},
							},
						},
						Type: pointer.Get("search-index-type"),
					},
					GroupId: pointer.Get("group-id-101"),
				},
			},
			target: &admin2025.SearchIndexCreateRequest{},
			want: &admin2025.SearchIndexCreateRequest{
				CollectionName: "collection-name",
				Database:       "database-name",
				Name:           "index-name",
				Type:           pointer.Get("search-index-type"),
			},
		},
		{
			name:       "searchindex update request fields",
			crd:        "SearchIndex",
			sdkVersion: "v20250312",
			spec: v1.SearchIndexSpec{
				V20250312: &v1.SearchIndexSpecV20250312{
					Entry: &v1.SearchIndexSpecV20250312Entry{
						Analyzer: pointer.Get("lucene.standard"),
						Database: "database-name",
						Analyzers: &[]v1.Analyzers{
							{
								Name: "custom-analyzer",
								CharFilters: &[]apiextensionsv1.JSON{
									{Raw: []byte(`{"key":"value"}`)},
									{Raw: []byte(`{"key2":"value2"}`)},
								},
								TokenFilters: &[]apiextensionsv1.JSON{
									{Raw: []byte(`{"key3":"value3"}`)},
									{Raw: []byte(`{"key4":"value4"}`)},
								},
								Tokenizer: v1.Tokenizer{
									Group:          pointer.Get(2),
									MaxGram:        pointer.Get(100),
									MaxTokenLength: pointer.Get(50),
									MinGram:        pointer.Get(1),
									Pattern:        pointer.Get("pattern"),
									Type:           pointer.Get("custom"),
								},
							},
						},
						CollectionName: "collection-name",
						Fields: &[]apiextensionsv1.JSON{
							{Raw: []byte(`{"field1":"value1"}`)},
							{Raw: []byte(`{"field2":"value2"}`)},
							{Raw: []byte(`{"field3":"value3"}`)},
						},
						Mappings: &v1.Mappings{
							Dynamic: pointer.Get(true),
							Fields: &apiextensionsv1.JSON{
								Raw: []byte(`{"field4":"value4"}`),
							},
						},
						Name:           "index-name",
						NumPartitions:  pointer.Get(3),
						SearchAnalyzer: pointer.Get("lucene.standard"),
						StoredSource: &apiextensionsv1.JSON{
							Raw: []byte(`{"enabled": true}`),
						},
						Synonyms: &[]v1.Synonyms{
							{
								Analyzer: "synonym-analyzer",
								Name:     "synonym-name",
								Source: v1.Source{
									Collection: "synonym-collection",
								},
							},
						},
						Type: pointer.Get("search-index-type"),
					},
					GroupId: pointer.Get("group-id-101"),
				},
			},
			target: &admin2025.SearchIndexUpdateRequestDefinition{},
			want: &admin2025.SearchIndexUpdateRequestDefinition{
				Analyzer: pointer.Get("lucene.standard"),
				Analyzers: &[]admin2025.AtlasSearchAnalyzer{
					{
						Name: "custom-analyzer",
						CharFilters: &[]any{
							map[string]any{"key": "value"},
							map[string]any{"key2": "value2"},
						},
						TokenFilters: &[]any{
							map[string]any{"key3": "value3"},
							map[string]any{"key4": "value4"},
						},
						Tokenizer: map[string]any{
							"group":          2.0,
							"maxGram":        100.0,
							"maxTokenLength": 50.0,
							"minGram":        1.0,
							"pattern":        "pattern",
							"type":           "custom",
						},
					},
				},
				Fields: &[]any{
					map[string]any{"field1": "value1"},
					map[string]any{"field2": "value2"},
					map[string]any{"field3": "value3"},
				},
				Mappings: &admin2025.SearchMappings{
					Dynamic: pointer.Get(true),
					Fields: &map[string]any{
						"field4": "value4",
					},
				},
				NumPartitions:  pointer.Get(3),
				SearchAnalyzer: pointer.Get("lucene.standard"),
				StoredSource:   any(map[string]any{"enabled": true}),
				Synonyms: &[]admin2025.SearchSynonymMappingDefinition{
					{
						Analyzer: "synonym-analyzer",
						Name:     "synonym-name",
						Source: admin2025.SynonymSource{
							Collection: "synonym-collection",
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			crdsYML, err := samples.Open("samples/crds.yaml")
			require.NoError(t, err)
			defer crdsYML.Close()
			crd, err := extractCRD(tc.crd, bufio.NewScanner(crdsYML))
			require.NoError(t, err)
			typeInfo := translate.TypeInfo{
				CRDVersion: version,
				SDKVersion: tc.sdkVersion,
				CRD:        crd,
			}
			require.NoError(t, translate.ToAPI(&typeInfo, tc.target, &tc.spec, tc.deps...))
			assert.Equal(t, tc.want, tc.target)
		})
	}
}

func extractCRD(kind string, scanner *bufio.Scanner) (*apiextensionsv1.CustomResourceDefinition, error) {
	for {
		crd, err := crd2go.ParseCRD(scanner)
		if err != nil {
			return nil, fmt.Errorf("failed to extract CRD schema for kind %q: %w", kind, err)
		}
		if crd.Spec.Names.Kind == kind {
			return crd, nil
		}
	}
}
