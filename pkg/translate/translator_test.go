package translate_test

import (
	"bufio"
	"bytes"
	_ "embed"
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

	"github.com/josvazg/akotranslate/internal/crds"
	"github.com/josvazg/akotranslate/internal/pointer"
	"github.com/josvazg/akotranslate/pkg/translate"
	v1 "github.com/josvazg/akotranslate/pkg/translate/samples/v1"
	"github.com/josvazg/crd2go/k8s"
)

const (
	version = "v1"

	sdkVersion = "v20250312"
)

//go:embed samples/crds.yaml
var crdsYAMLBytes []byte

func TestFromAPI(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "empty",
			test: func(t *testing.T) {
				input := admin2025.Group{
					ClusterCount: 0,
					Id:           pointer.Get("6127378123219"),
					Name:         "test-project",
					OrgId:        "6129312312334",
					Tags: &[]admin2025.ResourceTag{
						{
							Key:   "key0",
							Value: "value0",
						},
						{
							Key:   "key",
							Value: "value",
						},
					},
					WithDefaultAlertsSettings: pointer.Get(true),
				}
				target := v1.Group{}
				want := []client.Object{
					&v1.Group{
						Spec: v1.GroupSpec{
							V20250312: &v1.GroupSpecV20250312{
								Entry: &v1.V20250312Entry{
									Name:  "test-project",
									OrgId: "6129312312334",
									Tags: &[]v1.Tags{
										{
											Key:   "key0",
											Value: "value0",
										},
										{
											Key:   "key",
											Value: "value",
										},
									},
									WithDefaultAlertsSettings: pointer.Get(true),
								},
							},
						},
						Status: v1.GroupStatus{
							V20250312: &v1.GroupStatusV20250312{
								Created: "0001-01-01T00:00:00Z", // TODO: how to remove this?
								Id:      pointer.Get("6127378123219"),
							},
						},
					},
				}
				testFromAPI(t, "Group", &target, &input, want)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func testFromAPI[S any, T any, P interface {
	*T
	client.Object
}](t *testing.T, kind string, target P, input *S, want []client.Object) {
	crdsYML := bytes.NewBuffer(crdsYAMLBytes)
	crd, err := extractCRD(kind, bufio.NewScanner(crdsYML))
	require.NoError(t, err)
	deps := translate.NewStaticDependencies("ns")
	translator := translate.NewTranslator(crd, version, sdkVersion, deps)
	results, err := translate.FromAPI(translator, target, input)
	require.NoError(t, err)
	assert.Equal(t, want, results)
}

func TestToAPIAllRefs(t *testing.T) {
	for _, tc := range []struct {
		name   string
		crd    string
		input  client.Object
		deps   []client.Object
		target admin2025.CreateAlertConfigurationApiParams
		want   admin2025.CreateAlertConfigurationApiParams
	}{
		{
			name: "group alert config with a groupRef and secrets",
			crd:  "GroupAlertsConfig",
			input: &v1.GroupAlertsConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "GroupAlertsConfig",
					APIVersion: "atlas.generated.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-group-alerts-config",
					Namespace: "ns",
				},
				Spec: v1.GroupAlertsConfigSpec{
					V20250312: &v1.GroupAlertsConfigSpecV20250312{
						Entry: &v1.GroupAlertsConfigSpecV20250312Entry{
							Enabled:       pointer.Get(true),
							EventTypeName: pointer.Get("some-event"),
							Matchers: &[]v1.Matchers{
								{
									FieldName: "field1",
									Operator:  "op1",
									Value:     "value1",
								},
								{
									FieldName: "field2",
									Operator:  "op2",
									Value:     "value2",
								},
							},
							MetricThreshold: &v1.MetricThreshold{
								MetricName: "metric",
								Mode:       pointer.Get("mode"),
								Operator:   pointer.Get("operator"),
								Threshold:  pointer.Get(1.0),
								Units:      pointer.Get("unit"),
							},
							Notifications: &[]v1.Notifications{
								{
									DatadogApiKeySecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.Get("alert-secrets-0"),
										Key:  pointer.Get("apiKey"),
									},
									DatadogRegion: pointer.Get("US"),
								},
								{
									WebhookSecretSecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.Get("alert-secrets-0"),
										Key:  pointer.Get("webhookSecret"),
									},
									WebhookUrlSecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.Get("alert-secrets-1"),
										Key:  pointer.Get("webhookUrl"),
									},
								},
							},
							SeverityOverride: pointer.Get("severe"),
							Threshold: &v1.MetricThreshold{
								MetricName: "metric",
								Mode:       pointer.Get("mode-t"),
								Operator:   pointer.Get("op-t"),
								Threshold:  pointer.Get(2.0),
								Units:      pointer.Get("unit-t"),
							},
						},
						GroupRef: &k8s.LocalReference{
							Name: "my-project",
						},
					},
				},
			},
			deps: []client.Object{
				&v1.Group{
					TypeMeta:   metav1.TypeMeta{Kind: "Group", APIVersion: "atlas.generated.mongodb.com/v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "my-project", Namespace: "ns"},
					Spec: v1.GroupSpec{
						V20250312: &v1.GroupSpecV20250312{
							Entry: &v1.V20250312Entry{
								Name:  "some-project",
								OrgId: "621454123423x125235142",
							},
						},
					},
					Status: v1.GroupStatus{
						V20250312: &v1.GroupStatusV20250312{
							Id: pointer.Get("62b6e34b3d91647abb20e7b8"),
						},
					},
				},
				&corev1.Secret{
					TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "alert-secrets-0", Namespace: "ns"},
					Data: map[string][]byte{
						"apiKey":        ([]byte)("sample-api-key"),
						"webhookSecret": ([]byte)("sample-webhook-secret"),
					},
				},
				&corev1.Secret{
					TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "alert-secrets-1", Namespace: "ns"},
					Data: map[string][]byte{
						"webhookUrl": ([]byte)("sample-webhook-url"),
					},
				},
			},
			want: admin2025.CreateAlertConfigurationApiParams{
				GroupId: "62b6e34b3d91647abb20e7b8",
				GroupAlertsConfig: &admin2025.GroupAlertsConfig{
					Enabled:       pointer.Get(true),
					EventTypeName: pointer.Get("some-event"),
					Matchers: &[]admin2025.StreamsMatcher{
						{
							FieldName: "field1",
							Operator:  "op1",
							Value:     "value1",
						},
						{
							FieldName: "field2",
							Operator:  "op2",
							Value:     "value2",
						},
					},
					Notifications: &[]admin2025.AlertsNotificationRootForGroup{
						{
							DatadogApiKey: pointer.Get("sample-api-key"),
							DatadogRegion: pointer.Get("US"),
						},
						{
							WebhookSecret: pointer.Get("sample-webhook-secret"),
							WebhookUrl:    pointer.Get("sample-webhook-url"),
						},
					},
					SeverityOverride: pointer.Get("severe"),
					MetricThreshold: &admin2025.FlexClusterMetricThreshold{
						MetricName: "metric",
						Mode:       pointer.Get("mode"),
						Operator:   pointer.Get("operator"),
						Threshold:  pointer.Get(1.0),
						Units:      pointer.Get("unit"),
					},
					Threshold: &admin2025.StreamProcessorMetricThreshold{
						MetricName: pointer.Get("metric"),
						Mode:       pointer.Get("mode-t"),
						Operator:   pointer.Get("op-t"),
						Threshold:  pointer.Get(2.0),
						Units:      pointer.Get("unit-t"),
					},
				},
			},
			target: admin2025.CreateAlertConfigurationApiParams{},
		},

		{
			name: "group alert config with secrets but a direct groupId",
			crd:  "GroupAlertsConfig",
			input: &v1.GroupAlertsConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "GroupAlertsConfig",
					APIVersion: "atlas.generated.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-group-alerts-config",
					Namespace: "ns",
				},
				Spec: v1.GroupAlertsConfigSpec{
					V20250312: &v1.GroupAlertsConfigSpecV20250312{
						Entry: &v1.GroupAlertsConfigSpecV20250312Entry{
							Enabled:       pointer.Get(true),
							EventTypeName: pointer.Get("some-event"),
							Matchers: &[]v1.Matchers{
								{
									FieldName: "field1",
									Operator:  "op1",
									Value:     "value1",
								},
								{
									FieldName: "field2",
									Operator:  "op2",
									Value:     "value2",
								},
							},
							MetricThreshold: &v1.MetricThreshold{
								MetricName: "metric",
								Mode:       pointer.Get("mode"),
								Operator:   pointer.Get("operator"),
								Threshold:  pointer.Get(1.0),
								Units:      pointer.Get("unit"),
							},
							Notifications: &[]v1.Notifications{
								{
									DatadogApiKeySecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.Get("alert-secrets-0"),
										Key:  pointer.Get("apiKey"),
									},
									DatadogRegion: pointer.Get("US"),
								},
								{
									WebhookSecretSecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.Get("alert-secrets-0"),
										Key:  pointer.Get("webhookSecret"),
									},
									WebhookUrlSecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.Get("alert-secrets-1"),
										Key:  pointer.Get("webhookUrl"),
									},
								},
							},
							SeverityOverride: pointer.Get("severe"),
							Threshold: &v1.MetricThreshold{
								MetricName: "metric",
								Mode:       pointer.Get("mode-t"),
								Operator:   pointer.Get("op-t"),
								Threshold:  pointer.Get(2.0),
								Units:      pointer.Get("unit-t"),
							},
						},
						GroupId: pointer.Get("62b6e34b3d91647abb20e7b8"),
					},
				},
			},
			deps: []client.Object{
				&corev1.Secret{
					TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "alert-secrets-0", Namespace: "ns"},
					Data: map[string][]byte{
						"apiKey":        ([]byte)("sample-api-key"),
						"webhookSecret": ([]byte)("sample-webhook-secret"),
					},
				},
				&corev1.Secret{
					TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "alert-secrets-1", Namespace: "ns"},
					Data: map[string][]byte{
						"webhookUrl": ([]byte)("sample-webhook-url"),
					},
				},
			},
			want: admin2025.CreateAlertConfigurationApiParams{
				GroupId: "62b6e34b3d91647abb20e7b8",
				GroupAlertsConfig: &admin2025.GroupAlertsConfig{
					Enabled:       pointer.Get(true),
					EventTypeName: pointer.Get("some-event"),
					Matchers: &[]admin2025.StreamsMatcher{
						{
							FieldName: "field1",
							Operator:  "op1",
							Value:     "value1",
						},
						{
							FieldName: "field2",
							Operator:  "op2",
							Value:     "value2",
						},
					},
					Notifications: &[]admin2025.AlertsNotificationRootForGroup{
						{
							DatadogApiKey: pointer.Get("sample-api-key"),
							DatadogRegion: pointer.Get("US"),
						},
						{
							WebhookSecret: pointer.Get("sample-webhook-secret"),
							WebhookUrl:    pointer.Get("sample-webhook-url"),
						},
					},
					SeverityOverride: pointer.Get("severe"),
					MetricThreshold: &admin2025.FlexClusterMetricThreshold{
						MetricName: "metric",
						Mode:       pointer.Get("mode"),
						Operator:   pointer.Get("operator"),
						Threshold:  pointer.Get(1.0),
						Units:      pointer.Get("unit"),
					},
					Threshold: &admin2025.StreamProcessorMetricThreshold{
						MetricName: pointer.Get("metric"),
						Mode:       pointer.Get("mode-t"),
						Operator:   pointer.Get("op-t"),
						Threshold:  pointer.Get(2.0),
						Units:      pointer.Get("unit-t"),
					},
				},
			},
			target: admin2025.CreateAlertConfigurationApiParams{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			crdsYML := bytes.NewBuffer(crdsYAMLBytes)
			crd, err := extractCRD(tc.crd, bufio.NewScanner(crdsYML))
			require.NoError(t, err)
			deps := translate.NewStaticDependencies("ns", tc.deps...)
			translator := translate.NewTranslator(crd, version, sdkVersion, deps)
			// , reflect.TypeOf(tc.target)
			require.NoError(t, translate.ToAPI(translator, &tc.target, tc.input))
			assert.Equal(t, tc.want, tc.target)
		})
	}
}

// NetworkPermissions is a required struct wrapper to match the API structure
// TODO: do we need a mapping option? for this case a rename would suffice to
// load the entry array field as results in a PaginatedNetworkAccess.
// On the other hand, is extracting the whole list the proper way interact with the API?
type NetworkPermissions struct {
	Entry []admin2025.NetworkPermissionEntry `json:"entry"`
}

type testToAPICase[T any] struct {
	name   string
	crd    string
	input  client.Object
	deps   []client.Object
	target *T
	want   *T
}

func TestToAPI(t *testing.T) {
	for _, gtc := range []any{
		testToAPICase[admin2025.DataProtectionSettings20231001]{
			name: "sample backup compliance policy",
			crd:  "BackupCompliancePolicy",
			input: &v1.BackupCompliancePolicy{
				Spec: v1.BackupCompliancePolicySpec{
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
		testToAPICase[admin2025.DiskBackupSnapshotSchedule20240805]{
			name: "backup schedule all fields",
			crd:  "BackupSchedule",
			input: &v1.BackupSchedule{
				Spec: v1.BackupScheduleSpec{
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
		testToAPICase[admin2025.ClusterDescription20240805]{
			name: "cluster all fields",
			crd:  "Cluster",
			input: &v1.Cluster{
				Spec: v1.ClusterSpec{
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
		testToAPICase[admin2025.DataLakeTenant]{
			name: "data federation all fields",
			crd:  "DataFederation",
			input: &v1.DataFederation{
				Spec: v1.DataFederationSpec{
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
		testToAPICase[admin2025.CloudDatabaseUser]{
			name: "sample database user",
			crd:  "DatabaseUser",
			input: &v1.DatabaseUser{
				Spec: v1.DatabaseUserSpec{
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
		testToAPICase[admin2025.FlexClusterDescriptionCreate20241113]{
			name: "flex cluster with all fields",
			crd:  "FlexCluster",
			input: &v1.FlexCluster{
				Spec: v1.FlexClusterSpec{
					V20250312: &v1.FlexClusterSpecV20250312{
						Entry: &v1.FlexClusterSpecV20250312Entry{
							Name: "flex-cluster-name",
							ProviderSettings: v1.ProviderSettings{
								BackingProviderName: "AWS",
								RegionName:          "us-east-1",
							},
							Tags: &[]v1.Tags{
								{Key: "key1", Value: "value1"},
								{Key: "key2", Value: "value2"},
							},
							TerminationProtectionEnabled: pointer.Get(true),
						},
						GroupId: pointer.Get("32b6e34b3d91647abb20e7b8"),
					},
				},
			},
			target: &admin2025.FlexClusterDescriptionCreate20241113{},
			want: &admin2025.FlexClusterDescriptionCreate20241113{
				Name: "flex-cluster-name",
				ProviderSettings: admin2025.FlexProviderSettingsCreate20241113{
					BackingProviderName: "AWS",
					RegionName:          "us-east-1",
				},
				Tags: &[]admin2025.ResourceTag{
					{Key: "key1", Value: "value1"},
					{Key: "key2", Value: "value2"},
				},
				TerminationProtectionEnabled: pointer.Get(true),
			},
		},
		testToAPICase[admin2025.Group]{
			name: "simple group",
			crd:  "Group",
			input: &v1.Group{
				Spec: v1.GroupSpec{
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
		testToAPICase[admin2025.GroupAlertsConfig]{
			name: "group alert config with project and credential references",
			crd:  "GroupAlertsConfig",
			input: &v1.GroupAlertsConfig{
				Spec: v1.GroupAlertsConfigSpec{
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
			},
			deps: []client.Object{
				&corev1.Secret{
					TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "datadog-secret", Namespace: "ns"},
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
		testToAPICase[admin2025.BaseNetworkPeeringConnectionSettings]{
			name: "sample network peering connection",
			crd:  "NetworkPeeringConnection",
			input: &v1.NetworkPeeringConnection{
				Spec: v1.NetworkPeeringConnectionSpec{
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
		testToAPICase[NetworkPermissions]{
			name: "network permission entries all fields",
			crd:  "NetworkPermissionEntries",
			input: &v1.NetworkPermissionEntries{
				Spec: v1.NetworkPermissionEntriesSpec{
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
		testToAPICase[admin2025.AtlasOrganization]{
			name: "sample organization",
			crd:  "Organization",
			input: &v1.Organization{
				Spec: v1.OrganizationSpec{
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
			},
			target: &admin2025.AtlasOrganization{},
			want: &admin2025.AtlasOrganization{
				Name:                      "org-name",
				SkipDefaultAlertsSettings: pointer.Get(true),
			},
		},
		testToAPICase[admin2025.OrganizationSettings]{
			name: "Organization setting with all fields",
			crd:  "OrganizationSetting",
			input: &v1.OrganizationSetting{
				Spec: v1.OrganizationSettingSpec{
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
		testToAPICase[admin2025.UserCustomDBRole]{
			name: "customrole with all fields",
			crd:  "CustomRole",
			input: &v1.CustomRole{
				Spec: v1.CustomRoleSpec{
					V20250312: &v1.CustomRoleSpecV20250312{
						Entry: &v1.CustomRoleSpecV20250312Entry{
							RoleName: "custom-role-name",
							Actions: &[]v1.Actions{
								{
									Action: "action1",
									Resources: &[]v1.Resources{
										{
											Collection: "collection0",
											Cluster:    true,
											Db:         "db0",
										},
										{
											Collection: "collection1",
											Cluster:    true,
											Db:         "db1",
										},
									},
								},
							},
							InheritedRoles: &[]v1.InheritedRoles{
								{
									Db:   "inherited-db-name1",
									Role: "inherited-role-name1",
								},
								{
									Db:   "inherited-db-name2",
									Role: "inherited-role-name2",
								},
							},
						},
					},
				},
			},
			target: &admin2025.UserCustomDBRole{},
			want: &admin2025.UserCustomDBRole{
				RoleName: "custom-role-name",
				Actions: &[]admin2025.DatabasePrivilegeAction{
					{
						Action: "action1",
						Resources: &[]admin2025.DatabasePermittedNamespaceResource{
							{
								Collection: "collection0",
								Cluster:    true,
								Db:         "db0",
							},
							{
								Collection: "collection1",
								Cluster:    true,
								Db:         "db1",
							},
						},
					},
				},
				InheritedRoles: &[]admin2025.DatabaseInheritedRole{
					{
						Db:   "inherited-db-name1",
						Role: "inherited-role-name1",
					},
					{
						Db:   "inherited-db-name2",
						Role: "inherited-role-name2",
					},
				},
			},
		},
		// "SampleDataset" only holds a name, there is no SDK API struc for the request
		// {
		// 	name:       "sample dataset all fields",
		// 	crd:        "SampleDataset",
		// 	input: &v1.SampleDataset {
		//      Spec: v1.SampleDatasetSpec{
		// 		    V20250312: &v1.SampleDatasetSpecV20250312{
		// 			    Name:     "sample-dataset",
		// 			    GroupId:  pointer.Get("32b6e34b3d91647abb20e7b8"),
		// 	 	    },
		// 	    },
		//  },
		// 	target: admin2025.SampleDatasetStatus{},
		// 	want:   admin2025.SampleDatasetStatus{},
		// },
		testToAPICase[admin2025.SearchIndexCreateRequest]{
			name: "searchindex create request fields",
			crd:  "SearchIndex",
			input: &v1.SearchIndex{
				Spec: v1.SearchIndexSpec{
					V20250312: &v1.SearchIndexSpecV20250312{
						Entry: &v1.SearchIndexSpecV20250312Entry{
							Database:       "database-name",
							CollectionName: "collection-name",
							Name:           "index-name",
							Type:           pointer.Get("search-index-type"),
							Definition: &v1.Definition{
								Analyzer: pointer.Get("lucene.standard"),
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
										Tokenizer: apiextensionsv1.JSON{
											Raw: []byte(`{"group":2,"maxGram":100,"maxTokenLength":50,"minGram":1,"pattern":"pattern","type":"custom"}`),
										},
									},
								},
								Fields: &[]apiextensionsv1.JSON{
									{Raw: []byte(`{"field1":"value1"}`)},
									{Raw: []byte(`{"field2":"value2"}`)},
									{Raw: []byte(`{"field3":"value3"}`)},
								},
								Mappings: &v1.Mappings{
									Dynamic: pointer.Get(true),
									Fields: &map[string]apiextensionsv1.JSON{
										"field1": {Raw: []byte(`{"key4":"value4"}`)},
									},
								},
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
							},
						},
						GroupId: pointer.Get("group-id-101"),
					},
				},
			},
			target: &admin2025.SearchIndexCreateRequest{},
			want: &admin2025.SearchIndexCreateRequest{
				CollectionName: "collection-name",
				Database:       "database-name",
				Name:           "index-name",
				Type:           pointer.Get("search-index-type"),
				Definition: &admin2025.BaseSearchIndexCreateRequestDefinition{
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
							"field1": map[string]any{"key4": "value4"},
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
		},
		testToAPICase[admin2025.Team]{
			name: "team all fields",
			crd:  "Team",
			input: &v1.Team{
				Spec: v1.TeamSpec{
					V20250312: &v1.TeamSpecV20250312{
						Entry: &v1.TeamSpecV20250312Entry{
							Name:      "team-name",
							Usernames: []string{"user1", "user2"},
						},
						OrgId: "org-id",
					},
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
		testToAPICase[admin2025.ThirdPartyIntegration]{
			name: "third part integration all fields",
			crd:  "ThirdPartyIntegration",
			input: &v1.ThirdPartyIntegration{
				Spec: v1.ThirdPartyIntegrationSpec{
					V20250312: &v1.ThirdPartyIntegrationSpecV20250312{
						IntegrationType: "ANY",
						Entry: &v1.ThirdPartyIntegrationSpecV20250312Entry{
							AccountId: pointer.Get("account-id"),
							ApiKeySecretRef: &v1.ApiTokenSecretRef{
								Key:  pointer.Get("apiKey"),
								Name: pointer.Get("multi-secret0"),
							},
							ApiTokenSecretRef: &v1.ApiTokenSecretRef{
								Key:  pointer.Get("apiToken"),
								Name: pointer.Get("multi-secret0"),
							},
							ChannelName: pointer.Get("channel-name"),
							Enabled:     pointer.Get(true),
							LicenseKeySecretRef: &v1.ApiTokenSecretRef{
								Key:  pointer.Get("licenseKey"),
								Name: pointer.Get("multi-secret1"),
							},
							Region:                       pointer.Get("some-region"),
							SendCollectionLatencyMetrics: pointer.Get(true),
							SendDatabaseMetrics:          pointer.Get(true),
							SendUserProvidedResourceTags: pointer.Get(true),
							ServiceDiscovery:             pointer.Get("service-discovery"),
							TeamName:                     pointer.Get("some-team"),
							Type:                         pointer.Get("some-type"),
							Username:                     pointer.Get("username"),
						},
						GroupId: pointer.Get("32b6e34b3d91647abb20e7b8"),
					},
				},
			},
			target: &admin2025.ThirdPartyIntegration{},
			want: &admin2025.ThirdPartyIntegration{
				Type:                         pointer.Get("some-type"),
				ApiKey:                       pointer.Get("sample-api-key"),
				Region:                       pointer.Get("some-region"),
				SendCollectionLatencyMetrics: pointer.Get(true),
				SendDatabaseMetrics:          pointer.Get(true),
				SendUserProvidedResourceTags: pointer.Get(true),
				AccountId:                    pointer.Get("account-id"),
				LicenseKey:                   pointer.Get("sample-license-key"),
				Enabled:                      pointer.Get(true),
				ServiceDiscovery:             pointer.Get("service-discovery"),
				Username:                     pointer.Get("username"),
				ApiToken:                     pointer.Get("sample-api-token"),
				ChannelName:                  pointer.Get("channel-name"),
				TeamName:                     pointer.Get("some-team"),
			},
			deps: []client.Object{
				&corev1.Secret{
					TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "multi-secret0", Namespace: "ns"},
					Data: map[string][]byte{
						"apiKey":   ([]byte)("sample-api-key"),
						"apiToken": ([]byte)("sample-api-token"),
					},
				},
				&corev1.Secret{
					TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "multi-secret1", Namespace: "ns"},
					Data: map[string][]byte{
						"licenseKey": ([]byte)("sample-license-key"),
					},
				},
			},
		},
	} {
		switch tc := gtc.(type) {
		case testToAPICase[admin2025.DataProtectionSettings20231001]:
		case testToAPICase[admin2025.DiskBackupSnapshotSchedule20240805]:
		case testToAPICase[admin2025.ClusterDescription20240805]:
		case testToAPICase[admin2025.DataLakeTenant]:
		case testToAPICase[admin2025.CloudDatabaseUser]:
		case testToAPICase[admin2025.FlexClusterDescriptionCreate20241113]:
		case testToAPICase[admin2025.Group]:
		case testToAPICase[admin2025.GroupAlertsConfig]:
		case testToAPICase[admin2025.BaseNetworkPeeringConnectionSettings]:
		case testToAPICase[NetworkPermissions]:
		case testToAPICase[admin2025.AtlasOrganization]:
		case testToAPICase[admin2025.OrganizationSettings]:
		case testToAPICase[admin2025.UserCustomDBRole]:
		case testToAPICase[admin2025.SearchIndexCreateRequest]:
		case testToAPICase[admin2025.Team]:
		case testToAPICase[admin2025.ThirdPartyIntegration]:
			runTestToAPICase(t, tc)
		default:
			t.Fatalf("unsupported want type %T", tc)
		}
	}
}

func runTestToAPICase[T any](t *testing.T, tc testToAPICase[T]) {
	t.Run(tc.name, func(t *testing.T) {
		crdsYML := bytes.NewBuffer(crdsYAMLBytes)
		crd, err := extractCRD(tc.crd, bufio.NewScanner(crdsYML))
		require.NoError(t, err)
		deps := translate.NewStaticDependencies("ns", tc.deps...)
		translator := translate.NewTranslator(crd, version, sdkVersion, deps)
		require.NoError(t, translate.ToAPI(translator, tc.target, tc.input))
		assert.Equal(t, tc.want, tc.target)
	})
}

func extractCRD(kind string, scanner *bufio.Scanner) (*apiextensionsv1.CustomResourceDefinition, error) {
	for {
		crd, err := crds.Parse(scanner)
		if err != nil {
			return nil, fmt.Errorf("failed to extract CRD schema for kind %q: %w", kind, err)
		}
		if crd.Spec.Names.Kind == kind {
			return crd, nil
		}
	}
}
