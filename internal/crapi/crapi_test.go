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
//

package crapi_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin2025 "go.mongodb.org/atlas-sdk/v20250312006/admin"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/crds"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/refs"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/testdata"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/testdata/samples/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/k8s"
)

const (
	version = "v1"

	sdkVersion = "v20250312"

	testProjectID = "6098765432109876"
)

func TestFromAPI(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "simple group",
			test: func(t *testing.T) {
				input := admin2025.Group{
					Created:      time.Date(2025, 1, 1, 1, 30, 15, 0, time.UTC),
					ClusterCount: 0,
					Id:           pointer.MakePtr("6127378123219"),
					Name:         "test-project",
					OrgId:        testProjectID,
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
					WithDefaultAlertsSettings: pointer.MakePtr(true),
				}
				target := v1.Group{
					Spec: v1.GroupSpec{
						V20250312: &v1.GroupSpecV20250312{
							ProjectOwnerId: "",
						},
					},
				}
				want := &v1.Group{
					Spec: v1.GroupSpec{
						V20250312: &v1.GroupSpecV20250312{
							Entry: &v1.V20250312Entry{
								Name:  "test-project",
								OrgId: testProjectID,
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
								WithDefaultAlertsSettings: pointer.MakePtr(true),
							},
						},
					},
					Status: v1.GroupStatus{
						V20250312: &v1.GroupStatusV20250312{
							Created: "2025-01-01T01:30:15Z",
							Id:      pointer.MakePtr("6127378123219"),
						},
					},
				}
				testFromAPI(t, "Group", &target, &input, want)
			},
		},

		{
			name: "dbuser with secret and group refs",
			test: func(t *testing.T) {
				input := admin2025.CloudDatabaseUser{
					AwsIAMType:      pointer.MakePtr("NONE AWS"),
					DatabaseName:    "dbname",
					DeleteAfterDate: pointer.MakePtr(time.Date(2025, 2, 1, 1, 30, 15, 0, time.UTC)),
					Description:     pointer.MakePtr("sample db user"),
					GroupId:         testProjectID,
					Labels: &[]admin2025.ComponentLabel{
						{
							Key:   pointer.MakePtr("key0"),
							Value: pointer.MakePtr("value0"),
						},
						{
							Key:   pointer.MakePtr("key1"),
							Value: pointer.MakePtr("value1"),
						},
					},
					LdapAuthType: pointer.MakePtr("NONE LDAP"),
					OidcAuthType: pointer.MakePtr("NONE OIDC"),
					// TODO: new crd should put this on a secret
					Password: pointer.MakePtr("fakepass"),
					Roles: &[]admin2025.DatabaseUserRole{
						{
							CollectionName: pointer.MakePtr("collection0"),
							DatabaseName:   "mydb",
							RoleName:       "admin",
						},
					},
					Scopes: &[]admin2025.UserScope{
						{
							Name: "scopeName",
							Type: "scopeType",
						},
					},
					Username: "dbuser",
					X509Type: pointer.MakePtr("NONE X509"),
				}
				target := v1.DatabaseUser{}
				want := &v1.DatabaseUser{
					Spec: v1.DatabaseUserSpec{
						V20250312: &v1.DatabaseUserSpecV20250312{
							Entry: &v1.DatabaseUserSpecV20250312Entry{
								AwsIAMType:      pointer.MakePtr("NONE AWS"),
								DatabaseName:    "dbname",
								DeleteAfterDate: pointer.MakePtr("2025-02-01T01:30:15Z"),
								Description:     pointer.MakePtr("sample db user"),
								GroupId:         testProjectID, // ID ref by default
								Labels: &[]v1.Tags{
									{
										Key:   "key0",
										Value: "value0",
									},
									{
										Key:   "key1",
										Value: "value1",
									},
								},
								LdapAuthType: pointer.MakePtr("NONE LDAP"),
								OidcAuthType: pointer.MakePtr("NONE OIDC"),
								Roles: &[]v1.Roles{
									{
										CollectionName: pointer.MakePtr("collection0"),
										DatabaseName:   "mydb",
										RoleName:       "admin",
									},
								},
								Scopes: &[]v1.Scopes{
									{
										Name: "scopeName",
										Type: "scopeType",
									},
								},
								Username: "dbuser",
								X509Type: pointer.MakePtr("NONE X509"),
								Password: pointer.MakePtr("fakepass"), // TODO: this should go away to a secret
							},
							GroupId: pointer.MakePtr(testProjectID),
						},
					},
				}
				testFromAPI(t, "DatabaseUser", &target, &input, want)
			},
		},

		{
			name: "GroupAlertConfigs",
			test: func(t *testing.T) {
				input := admin2025.GroupAlertsConfig{
					Enabled:       pointer.MakePtr(true),
					EventTypeName: pointer.MakePtr("OUTSIDE_STREAM_PROCESSOR_METRIC_THRESHOLD"),
					GroupId:       pointer.MakePtr(testProjectID),
					Id:            pointer.MakePtr("notification id"),
					Matchers: &[]admin2025.StreamsMatcher{
						{
							FieldName: "field0",
							Operator:  "EQUALS",
							Value:     "value0",
						},
						{
							FieldName: "field1",
							Operator:  "GREATER",
							Value:     "value1",
						},
					},
					Notifications: &[]admin2025.AlertsNotificationRootForGroup{
						{
							DatadogApiKey: pointer.MakePtr("fake api key"),
							DatadogRegion: pointer.MakePtr("US"),
							DelayMin:      pointer.MakePtr(42),
							IntegrationId: pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
							IntervalMin:   pointer.MakePtr(43),
							NotifierId:    pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
							TypeName:      pointer.MakePtr("DATADOG"),
						},
					},
					SeverityOverride: pointer.MakePtr("CRITICIAL"),
					MetricThreshold: &admin2025.FlexClusterMetricThreshold{
						MetricName: "metric",
						Mode:       pointer.MakePtr("mode"),
						Operator:   pointer.MakePtr("op"),
						Threshold:  pointer.MakePtr(0.1),
						Units:      pointer.MakePtr("unit"),
					},
					Threshold: &admin2025.StreamProcessorMetricThreshold{
						MetricName: pointer.MakePtr("anotherMetric"),
						Mode:       pointer.MakePtr("a mode"),
						Operator:   pointer.MakePtr("an op"),
						Threshold:  pointer.MakePtr(0.2),
						Units:      pointer.MakePtr("a unit"),
					},
				}
				target := v1.GroupAlertsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "groupalertscfg",
						Namespace: "ns",
					},
				}
				want := &v1.GroupAlertsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "groupalertscfg",
						Namespace: "ns",
					},
					Spec: v1.GroupAlertsConfigSpec{
						V20250312: &v1.GroupAlertsConfigSpecV20250312{
							Entry: &v1.GroupAlertsConfigSpecV20250312Entry{
								Enabled:       pointer.MakePtr(true),
								EventTypeName: pointer.MakePtr("OUTSIDE_STREAM_PROCESSOR_METRIC_THRESHOLD"),
								Matchers: &[]v1.Matchers{
									{
										FieldName: "field0",
										Operator:  "EQUALS",
										Value:     "value0",
									},
									{
										FieldName: "field1",
										Operator:  "GREATER",
										Value:     "value1",
									},
								},
								MetricThreshold: &v1.MetricThreshold{
									MetricName: "metric",
									Mode:       pointer.MakePtr("mode"),
									Operator:   pointer.MakePtr("op"),
									Threshold:  pointer.MakePtr(0.1),
									Units:      pointer.MakePtr("unit"),
								},
								Notifications: &[]v1.Notifications{
									{
										DatadogApiKeySecretRef: &v1.ApiTokenSecretRef{
											Key:  pointer.MakePtr("datadogApiKey"),
											Name: pointer.MakePtr("groupalertscfg-f4f4b5f9c849fc4cbdc"),
										},
										DatadogRegion: pointer.MakePtr("US"),
										DelayMin:      pointer.MakePtr(42),
										IntegrationId: pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
										IntervalMin:   pointer.MakePtr(43),
										NotifierId:    pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
										TypeName:      pointer.MakePtr("DATADOG"),
									},
								},
								SeverityOverride: pointer.MakePtr("CRITICIAL"),
								Threshold: &v1.MetricThreshold{
									MetricName: "anotherMetric",
									Mode:       pointer.MakePtr("a mode"),
									Operator:   pointer.MakePtr("an op"),
									Threshold:  pointer.MakePtr(0.2),
									Units:      pointer.MakePtr("a unit"),
								},
							},
							GroupId: pointer.MakePtr(testProjectID),
						},
					},
					Status: v1.GroupAlertsConfigStatus{
						V20250312: &v1.GroupAlertsConfigStatusV20250312{
							GroupId: pointer.MakePtr(testProjectID),
							Id:      pointer.MakePtr("notification id"),
						},
					},
				}
				wantSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "groupalertscfg-f4f4b5f9c849fc4cbdc",
						Namespace: "ns",
					},
					Data: map[string][]byte{
						"datadogApiKey": ([]byte)("fake api key"),
					},
				}
				testFromAPI(t, "GroupAlertsConfig", &target, &input, want, wantSecret)
			},
		},

		{
			name: "ThirdPartyIntegration",
			test: func(t *testing.T) {
				input := admin2025.ThirdPartyIntegration{
					Id:          pointer.MakePtr("SomeID"),
					Type:        pointer.MakePtr("SLACK"),
					ApiToken:    pointer.MakePtr("some fake api token"),
					ChannelName: pointer.MakePtr("alert-channel"),
					TeamName:    pointer.MakePtr("some-team"),
				}
				target := v1.ThirdPartyIntegration{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "3rdparty-slack",
						Namespace: "ns",
					},
					Spec: v1.ThirdPartyIntegrationSpec{
						V20250312: &v1.ThirdPartyIntegrationSpecV20250312{
							// TODO: is this a valid trick?
							// This API struct, unlike others, does NOT include the Group ID
							// it is part of the parameters, but not the response
							GroupId: pointer.MakePtr(testProjectID),
							// TODO: similarly to the Group ID the IntegrationType would
							// be the aparameter thet corresponds with "type" in the response
							// but there is no indication of such semantics from the CRD
							IntegrationType: "SLACK",
						},
					},
				}
				want := &v1.ThirdPartyIntegration{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "3rdparty-slack",
						Namespace: "ns",
					},
					Spec: v1.ThirdPartyIntegrationSpec{
						V20250312: &v1.ThirdPartyIntegrationSpecV20250312{
							Entry: &v1.ThirdPartyIntegrationSpecV20250312Entry{
								Type: pointer.MakePtr("SLACK"),
								ApiTokenSecretRef: &v1.ApiTokenSecretRef{
									Name: pointer.MakePtr("3rdparty-slack-5798d555ff4dc66f7c99"),
									Key:  pointer.MakePtr("apiToken"),
								},
								ChannelName: pointer.MakePtr("alert-channel"),
								TeamName:    pointer.MakePtr("some-team"),
							},
							// Pre-existing from input
							GroupId: pointer.MakePtr(string(testProjectID)),
							// Pre-existing from input
							IntegrationType: "SLACK",
						},
					},
					Status: v1.ThirdPartyIntegrationStatus{
						V20250312: &v1.ThirdPartyIntegrationStatusV20250312{
							Id:   pointer.MakePtr("SomeID"),
							Type: pointer.MakePtr("SLACK"),
						},
					},
				}
				wantSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "3rdparty-slack-5798d555ff4dc66f7c99",
						Namespace: "ns",
					},
					Data: map[string][]byte{
						"apiToken": ([]byte)("some fake api token"),
					},
				}
				testFromAPI(t, "ThirdPartyIntegration", &target, &input, want, wantSecret)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func testFromAPI[S any, T any, P refs.PtrClientObj[T]](t *testing.T, kind string, target P, input *S, want client.Object, wantDeps ...client.Object) {
	crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
	crd, err := extractCRD(kind, bufio.NewScanner(crdsYML))
	require.NoError(t, err)
	tr, err := crapi.NewTranslator(crd, version, sdkVersion)
	require.NoError(t, err)
	r := crapi.Request{Translator: tr}
	results, err := crapi.FromAPI(&r, target, input)
	require.NoError(t, err)
	assert.Equal(t, want, target)
	assert.Equal(t, wantDeps, results)
}

func TestToAPIAllRefs(t *testing.T) {
	for _, tc := range []struct {
		name   string
		crd    string
		input  client.Object
		deps   []client.Object
		target admin2025.GroupAlertsConfig
		want   admin2025.GroupAlertsConfig
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
							Enabled:       pointer.MakePtr(true),
							EventTypeName: pointer.MakePtr("some-event"),
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
								Mode:       pointer.MakePtr("mode"),
								Operator:   pointer.MakePtr("operator"),
								Threshold:  pointer.MakePtr(1.0),
								Units:      pointer.MakePtr("unit"),
							},
							Notifications: &[]v1.Notifications{
								{
									DatadogApiKeySecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.MakePtr("alert-secrets-0"),
										Key:  pointer.MakePtr("apiKey"),
									},
									DatadogRegion: pointer.MakePtr("US"),
								},
								{
									WebhookSecretSecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.MakePtr("alert-secrets-0"),
										Key:  pointer.MakePtr("webhookSecret"),
									},
									WebhookUrlSecretRef: &v1.ApiTokenSecretRef{
										Name: pointer.MakePtr("alert-secrets-1"),
										Key:  pointer.MakePtr("webhookUrl"),
									},
								},
							},
							SeverityOverride: pointer.MakePtr("severe"),
							Threshold: &v1.MetricThreshold{
								MetricName: "metric",
								Mode:       pointer.MakePtr("mode-t"),
								Operator:   pointer.MakePtr("op-t"),
								Threshold:  pointer.MakePtr(2.0),
								Units:      pointer.MakePtr("unit-t"),
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
							Id: pointer.MakePtr("62b6e34b3d91647abb20e7b8"),
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
			want: admin2025.GroupAlertsConfig{
				Enabled:       pointer.MakePtr(true),
				EventTypeName: pointer.MakePtr("some-event"),
				GroupId:       pointer.MakePtr("62b6e34b3d91647abb20e7b8"),
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
						DatadogApiKey: pointer.MakePtr("sample-api-key"),
						DatadogRegion: pointer.MakePtr("US"),
					},
					{
						WebhookSecret: pointer.MakePtr("sample-webhook-secret"),
						WebhookUrl:    pointer.MakePtr("sample-webhook-url"),
					},
				},
				SeverityOverride: pointer.MakePtr("severe"),
				MetricThreshold: &admin2025.FlexClusterMetricThreshold{
					MetricName: "metric",
					Mode:       pointer.MakePtr("mode"),
					Operator:   pointer.MakePtr("operator"),
					Threshold:  pointer.MakePtr(1.0),
					Units:      pointer.MakePtr("unit"),
				},
				Threshold: &admin2025.StreamProcessorMetricThreshold{
					MetricName: pointer.MakePtr("metric"),
					Mode:       pointer.MakePtr("mode-t"),
					Operator:   pointer.MakePtr("op-t"),
					Threshold:  pointer.MakePtr(2.0),
					Units:      pointer.MakePtr("unit-t"),
				},
			},
			target: admin2025.GroupAlertsConfig{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
			crd, err := extractCRD(tc.crd, bufio.NewScanner(crdsYML))
			require.NoError(t, err)
			tr, err := crapi.NewTranslator(crd, version, sdkVersion)
			require.NoError(t, err)
			r := crapi.Request{Translator: tr, Dependencies: tc.deps}
			require.NoError(t, crapi.ToAPI(&r, &tc.target, tc.input))
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

func TestToAPI(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "sample backup compliance policy",
			test: func(t *testing.T) {
				input := &v1.BackupCompliancePolicy{
					Spec: v1.BackupCompliancePolicySpec{
						V20250312: &v1.BackupCompliancePolicySpecV20250312{
							Entry: &v1.BackupCompliancePolicySpecV20250312Entry{
								AuthorizedEmail:         "user@example.com",
								CopyProtectionEnabled:   pointer.MakePtr(true),
								EncryptionAtRestEnabled: pointer.MakePtr(true),
								AuthorizedUserFirstName: "first-name",
								AuthorizedUserLastName:  "last-name",
								OnDemandPolicyItem: &v1.OnDemandPolicyItem{
									FrequencyInterval: 1,
									FrequencyType:     "some-freq",
									RetentionUnit:     "some-unit",
									RetentionValue:    2,
								},
								PitEnabled:        pointer.MakePtr(true),
								ProjectId:         pointer.MakePtr("project-id"),
								RestoreWindowDays: pointer.MakePtr(3),
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
							GroupId:                 pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
							OverwriteBackupPolicies: true,
						},
					},
				}
				target := &admin2025.DataProtectionSettings20231001{}
				want := &admin2025.DataProtectionSettings20231001{
					AuthorizedEmail:         "user@example.com",
					CopyProtectionEnabled:   pointer.MakePtr(true),
					EncryptionAtRestEnabled: pointer.MakePtr(true),
					AuthorizedUserFirstName: "first-name",
					AuthorizedUserLastName:  "last-name",
					OnDemandPolicyItem: &admin2025.BackupComplianceOnDemandPolicyItem{
						FrequencyInterval: 1,
						FrequencyType:     "some-freq",
						RetentionUnit:     "some-unit",
						RetentionValue:    2,
					},
					PitEnabled:        pointer.MakePtr(true),
					ProjectId:         pointer.MakePtr("project-id"),
					RestoreWindowDays: pointer.MakePtr(3),
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
				}
				testToAPI(t, "BackupCompliancePolicy", input, nil, target, want)
			},
		},

		{
			name: "backup schedule all fields",
			test: func(t *testing.T) {
				input := &v1.BackupSchedule{
					Spec: v1.BackupScheduleSpec{
						V20250312: &v1.BackupScheduleSpecV20250312{
							Entry: &v1.BackupScheduleSpecV20250312Entry{
								ReferenceHourOfDay:    pointer.MakePtr(2),
								ReferenceMinuteOfHour: pointer.MakePtr(30),
								RestoreWindowDays:     pointer.MakePtr(7),
								UpdateSnapshots:       pointer.MakePtr(true),
								AutoExportEnabled:     pointer.MakePtr(true),
								CopySettings: &[]v1.CopySettings{
									{
										CloudProvider:    pointer.MakePtr("AWS"),
										Frequencies:      &[]string{"freq-1", "freq-2"},
										RegionName:       pointer.MakePtr("us-east-1"),
										ShouldCopyOplogs: pointer.MakePtr(true),
										ZoneId:           "zone-id",
									},
									{
										CloudProvider:    pointer.MakePtr("GCE"),
										Frequencies:      &[]string{"freq-3", "freq-4"},
										RegionName:       pointer.MakePtr("us-east-3"),
										ShouldCopyOplogs: pointer.MakePtr(true),
										ZoneId:           "zone-id-0",
									},
								},
								DeleteCopiedBackups: &[]v1.DeleteCopiedBackups{
									{
										CloudProvider: pointer.MakePtr("Azure"),
										RegionName:    pointer.MakePtr("us-west-2"),
										ZoneId:        pointer.MakePtr("zone-id"),
									},
								},
								Export: &v1.Export{
									ExportBucketId: pointer.MakePtr("ExportBucketId"),
									FrequencyType:  pointer.MakePtr("FrequencyType"),
								},
								ExtraRetentionSettings: &[]v1.ExtraRetentionSettings{
									{
										FrequencyType: pointer.MakePtr("FrequencyType0"),
										RetentionDays: pointer.MakePtr(1),
									},
									{
										FrequencyType: pointer.MakePtr("FrequencyType1"),
										RetentionDays: pointer.MakePtr(2),
									},
								},
								Policies: &[]v1.Policies{
									{
										Id: pointer.MakePtr("id0"),
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
								UseOrgAndGroupNamesInExportPrefix: pointer.MakePtr(true),
							},
							GroupId:     pointer.MakePtr("group-id-101"),
							ClusterName: "cluster-name",
						},
					},
				}
				target := &admin2025.DiskBackupSnapshotSchedule20240805{}
				want := &admin2025.DiskBackupSnapshotSchedule20240805{
					ReferenceHourOfDay:    pointer.MakePtr(2),
					ReferenceMinuteOfHour: pointer.MakePtr(30),
					RestoreWindowDays:     pointer.MakePtr(7),
					UpdateSnapshots:       pointer.MakePtr(true),
					AutoExportEnabled:     pointer.MakePtr(true),
					CopySettings: &[]admin2025.DiskBackupCopySetting20240805{
						{
							CloudProvider:    pointer.MakePtr("AWS"),
							Frequencies:      &[]string{"freq-1", "freq-2"},
							RegionName:       pointer.MakePtr("us-east-1"),
							ShouldCopyOplogs: pointer.MakePtr(true),
							ZoneId:           "zone-id",
						},
						{
							CloudProvider:    pointer.MakePtr("GCE"),
							Frequencies:      &[]string{"freq-3", "freq-4"},
							RegionName:       pointer.MakePtr("us-east-3"),
							ShouldCopyOplogs: pointer.MakePtr(true),
							ZoneId:           "zone-id-0",
						},
					},
					DeleteCopiedBackups: &[]admin2025.DeleteCopiedBackups20240805{
						{
							CloudProvider: pointer.MakePtr("Azure"),
							RegionName:    pointer.MakePtr("us-west-2"),
							ZoneId:        pointer.MakePtr("zone-id"),
						},
					},
					Export: &admin2025.AutoExportPolicy{
						ExportBucketId: pointer.MakePtr("ExportBucketId"),
						FrequencyType:  pointer.MakePtr("FrequencyType"),
					},
					ExtraRetentionSettings: &[]admin2025.ExtraRetentionSetting{
						{
							FrequencyType: pointer.MakePtr("FrequencyType0"),
							RetentionDays: pointer.MakePtr(1),
						},
						{
							FrequencyType: pointer.MakePtr("FrequencyType1"),
							RetentionDays: pointer.MakePtr(2),
						},
					},
					Policies: &[]admin2025.AdvancedDiskBackupSnapshotSchedulePolicy{
						{
							Id: pointer.MakePtr("id0"),
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
					UseOrgAndGroupNamesInExportPrefix: pointer.MakePtr(true),
					ClusterName:                       pointer.MakePtr("cluster-name"),
				}
				testToAPI(t, "BackupSchedule", input, nil, target, want)
			},
		},

		//nolint:dupl
		{
			name: "cluster all fields",
			test: func(t *testing.T) {
				input := &v1.Cluster{
					Spec: v1.ClusterSpec{
						V20250312: &v1.ClusterSpecV20250312{
							Entry: &v1.ClusterSpecV20250312Entry{
								AcceptDataRisksAndForceReplicaSetReconfig: pointer.MakePtr("2025-01-01T00:00:00Z"),
								AdvancedConfiguration: &v1.AdvancedConfiguration{
									CustomOpensslCipherConfigTls12: &[]string{
										"TLS_AES_256_GCM_SHA384", "TLS_CHACHA20_POLY1305_SHA256",
									},
									MinimumEnabledTlsProtocol: pointer.MakePtr("TLS1.2"),
									TlsCipherConfigMode:       pointer.MakePtr("Custom"),
								},
								BackupEnabled:                             pointer.MakePtr(true),
								BiConnector:                               &v1.BiConnector{Enabled: pointer.MakePtr(true)},
								ClusterType:                               pointer.MakePtr("ReplicaSet"),
								ConfigServerManagementMode:                pointer.MakePtr("Managed"),
								ConfigServerType:                          pointer.MakePtr("ReplicaSet"),
								DiskWarmingMode:                           pointer.MakePtr("Enabled"),
								EncryptionAtRestProvider:                  pointer.MakePtr("AWS-KMS"),
								FeatureCompatibilityVersion:               pointer.MakePtr("7.0"),
								FeatureCompatibilityVersionExpirationDate: pointer.MakePtr("2025-12-31T00:00:00Z"),
								GlobalClusterSelfManagedSharding:          pointer.MakePtr(true),
								Labels: &[]v1.Tags{
									{Key: "key1", Value: "value1"},
									{Key: "key2", Value: "value2"},
								},
								MongoDBEmployeeAccessGrant: &v1.MongoDBEmployeeAccessGrant{
									ExpirationTime: "2025-12-31T00:00:00Z",
									GrantType:      "Temporary",
								},
								MongoDBMajorVersion:       pointer.MakePtr("8.0"),
								Name:                      pointer.MakePtr("my-cluster"),
								Paused:                    pointer.MakePtr(true),
								PitEnabled:                pointer.MakePtr(true),
								RedactClientLogData:       pointer.MakePtr(true),
								ReplicaSetScalingStrategy: pointer.MakePtr("Auto"),
								ReplicationSpecs: &[]v1.ReplicationSpecs{
									{
										ZoneId:   pointer.MakePtr("zone-id-1"),
										ZoneName: pointer.MakePtr("zone-name-1"),
										RegionConfigs: &[]v1.RegionConfigs{
											{
												RegionName: pointer.MakePtr("us-east-1"),
												AnalyticsSpecs: &v1.AnalyticsSpecs{
													DiskIOPS:      pointer.MakePtr(1000),
													DiskSizeGB:    pointer.MakePtr(10.0),
													EbsVolumeType: pointer.MakePtr("gp2"),
													InstanceSize:  pointer.MakePtr("M10"),
													NodeCount:     pointer.MakePtr(3),
												},
												AutoScaling: &v1.AnalyticsAutoScaling{
													Compute: &v1.Compute{
														Enabled:           pointer.MakePtr(true),
														ScaleDownEnabled:  pointer.MakePtr(true),
														MaxInstanceSize:   pointer.MakePtr("M20"),
														MinInstanceSize:   pointer.MakePtr("M10"),
														PredictiveEnabled: pointer.MakePtr(true),
													},
													DiskGB: &v1.DiskGB{
														Enabled: pointer.MakePtr(true),
													},
												},
												AnalyticsAutoScaling: &v1.AnalyticsAutoScaling{
													Compute: &v1.Compute{
														Enabled:           pointer.MakePtr(true),
														ScaleDownEnabled:  pointer.MakePtr(true),
														MaxInstanceSize:   pointer.MakePtr("M30"),
														MinInstanceSize:   pointer.MakePtr("M10"),
														PredictiveEnabled: pointer.MakePtr(true),
													},
													DiskGB: &v1.DiskGB{
														Enabled: pointer.MakePtr(true),
													},
												},
												BackingProviderName: pointer.MakePtr("AWS"),
												ElectableSpecs: &v1.ElectableSpecs{
													DiskIOPS:              pointer.MakePtr(1000),
													DiskSizeGB:            pointer.MakePtr(10.0),
													EbsVolumeType:         pointer.MakePtr("gp2"),
													EffectiveInstanceSize: pointer.MakePtr("M10"),
													InstanceSize:          pointer.MakePtr("M10"),
													NodeCount:             pointer.MakePtr(3),
												},
												Priority:     pointer.MakePtr(1),
												ProviderName: pointer.MakePtr("AWS"),
												ReadOnlySpecs: &v1.AnalyticsSpecs{
													DiskIOPS:      pointer.MakePtr(1000),
													DiskSizeGB:    pointer.MakePtr(10.0),
													EbsVolumeType: pointer.MakePtr("gp2"),
													InstanceSize:  pointer.MakePtr("M10"),
													NodeCount:     pointer.MakePtr(3),
												},
											},
											{
												RegionName: pointer.MakePtr("us-east-2"),
												AnalyticsSpecs: &v1.AnalyticsSpecs{
													DiskIOPS:      pointer.MakePtr(2000),
													DiskSizeGB:    pointer.MakePtr(10.0),
													EbsVolumeType: pointer.MakePtr("gp3"),
													InstanceSize:  pointer.MakePtr("M20"),
													NodeCount:     pointer.MakePtr(3),
												},
												AutoScaling: &v1.AnalyticsAutoScaling{
													Compute: &v1.Compute{
														Enabled:           pointer.MakePtr(true),
														ScaleDownEnabled:  pointer.MakePtr(true),
														MaxInstanceSize:   pointer.MakePtr("M50"),
														MinInstanceSize:   pointer.MakePtr("M20"),
														PredictiveEnabled: pointer.MakePtr(true),
													},
													DiskGB: &v1.DiskGB{
														Enabled: pointer.MakePtr(true),
													},
												},
												AnalyticsAutoScaling: &v1.AnalyticsAutoScaling{
													Compute: &v1.Compute{
														Enabled:           pointer.MakePtr(true),
														ScaleDownEnabled:  pointer.MakePtr(true),
														MaxInstanceSize:   pointer.MakePtr("M40"),
														MinInstanceSize:   pointer.MakePtr("M10"),
														PredictiveEnabled: pointer.MakePtr(true),
													},
													DiskGB: &v1.DiskGB{
														Enabled: pointer.MakePtr(true),
													},
												},
												BackingProviderName: pointer.MakePtr("AWS"),
												ElectableSpecs: &v1.ElectableSpecs{
													DiskIOPS:              pointer.MakePtr(1000),
													DiskSizeGB:            pointer.MakePtr(10.0),
													EbsVolumeType:         pointer.MakePtr("gp2"),
													EffectiveInstanceSize: pointer.MakePtr("M10"),
													InstanceSize:          pointer.MakePtr("M10"),
													NodeCount:             pointer.MakePtr(3),
												},
												Priority:     pointer.MakePtr(1),
												ProviderName: pointer.MakePtr("AWS"),
												ReadOnlySpecs: &v1.AnalyticsSpecs{
													DiskIOPS:      pointer.MakePtr(1000),
													DiskSizeGB:    pointer.MakePtr(10.0),
													EbsVolumeType: pointer.MakePtr("gp2"),
													InstanceSize:  pointer.MakePtr("M10"),
													NodeCount:     pointer.MakePtr(3),
												},
											},
										},
									},
								},
								RootCertType: pointer.MakePtr("X509"),
								Tags: &[]v1.Tags{
									{Key: "key1", Value: "value1"},
									{Key: "key2", Value: "value2"},
								},
								TerminationProtectionEnabled: pointer.MakePtr(true),
								VersionReleaseSystem:         pointer.MakePtr("Atlas"),
							},
							GroupId: pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.ClusterDescription20240805{}
				want := &admin2025.ClusterDescription20240805{
					AcceptDataRisksAndForceReplicaSetReconfig: pointer.MakePtr(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					AdvancedConfiguration: &admin2025.ApiAtlasClusterAdvancedConfiguration{
						CustomOpensslCipherConfigTls12: &[]string{
							"TLS_AES_256_GCM_SHA384", "TLS_CHACHA20_POLY1305_SHA256",
						},
						MinimumEnabledTlsProtocol: pointer.MakePtr("TLS1.2"),
						TlsCipherConfigMode:       pointer.MakePtr("Custom"),
					},
					BackupEnabled:                             pointer.MakePtr(true),
					BiConnector:                               &admin2025.BiConnector{Enabled: pointer.MakePtr(true)},
					ClusterType:                               pointer.MakePtr("ReplicaSet"),
					ConfigServerManagementMode:                pointer.MakePtr("Managed"),
					ConfigServerType:                          pointer.MakePtr("ReplicaSet"),
					DiskWarmingMode:                           pointer.MakePtr("Enabled"),
					EncryptionAtRestProvider:                  pointer.MakePtr("AWS-KMS"),
					FeatureCompatibilityVersion:               pointer.MakePtr("7.0"),
					FeatureCompatibilityVersionExpirationDate: pointer.MakePtr(time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)),
					GlobalClusterSelfManagedSharding:          pointer.MakePtr(true),
					Labels: &[]admin2025.ComponentLabel{
						{Key: pointer.MakePtr("key1"), Value: pointer.MakePtr("value1")},
						{Key: pointer.MakePtr("key2"), Value: pointer.MakePtr("value2")},
					},
					MongoDBEmployeeAccessGrant: &admin2025.EmployeeAccessGrant{
						ExpirationTime: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
						GrantType:      "Temporary",
					},
					MongoDBMajorVersion:       pointer.MakePtr("8.0"),
					Name:                      pointer.MakePtr("my-cluster"),
					Paused:                    pointer.MakePtr(true),
					PitEnabled:                pointer.MakePtr(true),
					RedactClientLogData:       pointer.MakePtr(true),
					ReplicaSetScalingStrategy: pointer.MakePtr("Auto"),
					ReplicationSpecs: &[]admin2025.ReplicationSpec20240805{
						{
							ZoneId:   pointer.MakePtr("zone-id-1"),
							ZoneName: pointer.MakePtr("zone-name-1"),
							RegionConfigs: &[]admin2025.CloudRegionConfig20240805{
								{
									RegionName: pointer.MakePtr("us-east-1"),
									AnalyticsSpecs: &admin2025.DedicatedHardwareSpec20240805{
										DiskIOPS:      pointer.MakePtr(1000),
										DiskSizeGB:    pointer.MakePtr(10.0),
										EbsVolumeType: pointer.MakePtr("gp2"),
										InstanceSize:  pointer.MakePtr("M10"),
										NodeCount:     pointer.MakePtr(3),
									},
									AutoScaling: &admin2025.AdvancedAutoScalingSettings{
										Compute: &admin2025.AdvancedComputeAutoScaling{
											Enabled:           pointer.MakePtr(true),
											ScaleDownEnabled:  pointer.MakePtr(true),
											MaxInstanceSize:   pointer.MakePtr("M20"),
											MinInstanceSize:   pointer.MakePtr("M10"),
											PredictiveEnabled: pointer.MakePtr(true),
										},
										DiskGB: &admin2025.DiskGBAutoScaling{
											Enabled: pointer.MakePtr(true),
										},
									},
									AnalyticsAutoScaling: &admin2025.AdvancedAutoScalingSettings{
										Compute: &admin2025.AdvancedComputeAutoScaling{
											Enabled:           pointer.MakePtr(true),
											ScaleDownEnabled:  pointer.MakePtr(true),
											MaxInstanceSize:   pointer.MakePtr("M30"),
											MinInstanceSize:   pointer.MakePtr("M10"),
											PredictiveEnabled: pointer.MakePtr(true),
										},
										DiskGB: &admin2025.DiskGBAutoScaling{
											Enabled: pointer.MakePtr(true),
										},
									},
									BackingProviderName: pointer.MakePtr("AWS"),
									ElectableSpecs: &admin2025.HardwareSpec20240805{
										DiskIOPS:              pointer.MakePtr(1000),
										DiskSizeGB:            pointer.MakePtr(10.0),
										EbsVolumeType:         pointer.MakePtr("gp2"),
										EffectiveInstanceSize: pointer.MakePtr("M10"),
										InstanceSize:          pointer.MakePtr("M10"),
										NodeCount:             pointer.MakePtr(3),
									},
									Priority:     pointer.MakePtr(1),
									ProviderName: pointer.MakePtr("AWS"),
									ReadOnlySpecs: &admin2025.DedicatedHardwareSpec20240805{
										DiskIOPS:      pointer.MakePtr(1000),
										DiskSizeGB:    pointer.MakePtr(10.0),
										EbsVolumeType: pointer.MakePtr("gp2"),
										InstanceSize:  pointer.MakePtr("M10"),
										NodeCount:     pointer.MakePtr(3),
									},
								},
								{
									RegionName: pointer.MakePtr("us-east-2"),
									AnalyticsSpecs: &admin2025.DedicatedHardwareSpec20240805{
										DiskIOPS:      pointer.MakePtr(2000),
										DiskSizeGB:    pointer.MakePtr(10.0),
										EbsVolumeType: pointer.MakePtr("gp3"),
										InstanceSize:  pointer.MakePtr("M20"),
										NodeCount:     pointer.MakePtr(3),
									},
									AutoScaling: &admin2025.AdvancedAutoScalingSettings{
										Compute: &admin2025.AdvancedComputeAutoScaling{
											Enabled:           pointer.MakePtr(true),
											ScaleDownEnabled:  pointer.MakePtr(true),
											MaxInstanceSize:   pointer.MakePtr("M50"),
											MinInstanceSize:   pointer.MakePtr("M20"),
											PredictiveEnabled: pointer.MakePtr(true),
										},
										DiskGB: &admin2025.DiskGBAutoScaling{
											Enabled: pointer.MakePtr(true),
										},
									},
									AnalyticsAutoScaling: &admin2025.AdvancedAutoScalingSettings{
										Compute: &admin2025.AdvancedComputeAutoScaling{
											Enabled:           pointer.MakePtr(true),
											ScaleDownEnabled:  pointer.MakePtr(true),
											MaxInstanceSize:   pointer.MakePtr("M40"),
											MinInstanceSize:   pointer.MakePtr("M10"),
											PredictiveEnabled: pointer.MakePtr(true),
										},
										DiskGB: &admin2025.DiskGBAutoScaling{
											Enabled: pointer.MakePtr(true),
										},
									},
									BackingProviderName: pointer.MakePtr("AWS"),
									ElectableSpecs: &admin2025.HardwareSpec20240805{
										DiskIOPS:              pointer.MakePtr(1000),
										DiskSizeGB:            pointer.MakePtr(10.0),
										EbsVolumeType:         pointer.MakePtr("gp2"),
										EffectiveInstanceSize: pointer.MakePtr("M10"),
										InstanceSize:          pointer.MakePtr("M10"),
										NodeCount:             pointer.MakePtr(3),
									},
									Priority:     pointer.MakePtr(1),
									ProviderName: pointer.MakePtr("AWS"),
									ReadOnlySpecs: &admin2025.DedicatedHardwareSpec20240805{
										DiskIOPS:      pointer.MakePtr(1000),
										DiskSizeGB:    pointer.MakePtr(10.0),
										EbsVolumeType: pointer.MakePtr("gp2"),
										InstanceSize:  pointer.MakePtr("M10"),
										NodeCount:     pointer.MakePtr(3),
									},
								},
							},
						},
					},
					RootCertType: pointer.MakePtr("X509"),
					Tags: &[]admin2025.ResourceTag{
						{Key: "key1", Value: "value1"},
						{Key: "key2", Value: "value2"},
					},
					TerminationProtectionEnabled: pointer.MakePtr(true),
					VersionReleaseSystem:         pointer.MakePtr("Atlas"),
					GroupId:                      pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
				}
				testToAPI(t, "Cluster", input, nil, target, want)
			},
		},

		//nolint:dupl
		{
			name: "data federation all fields",
			test: func(t *testing.T) {
				input := &v1.DataFederation{
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
								Name: pointer.MakePtr("some-name"),
								Storage: &v1.Storage{
									Databases: &[]v1.Databases{
										{
											Collections: &[]v1.Collections{
												{
													DataSources: &[]v1.DataSources{
														{
															AllowInsecure:       pointer.MakePtr(true),
															Collection:          pointer.MakePtr("some-name"),
															CollectionRegex:     pointer.MakePtr("collection-regex"),
															Database:            pointer.MakePtr("db"),
															DatabaseRegex:       pointer.MakePtr("db-regex"),
															DatasetName:         pointer.MakePtr("dataset-name"),
															DatasetPrefix:       pointer.MakePtr("dataset-prefix"),
															DefaultFormat:       pointer.MakePtr("default-format"),
															Path:                pointer.MakePtr("path"),
															ProvenanceFieldName: pointer.MakePtr("provenqance-field-name"),
															StoreName:           pointer.MakePtr("store-name"),
															TrimLevel:           pointer.MakePtr(1),
															Urls:                &[]string{"url1", "url2"},
														},
													},
													Name: pointer.MakePtr("collection0"),
												},
											},
											MaxWildcardCollections: pointer.MakePtr(3),
											Name:                   pointer.MakePtr("db0"),
											Views: &[]v1.Views{
												{
													Name:     pointer.MakePtr("view0"),
													Pipeline: pointer.MakePtr("pipeline0"),
													Source:   pointer.MakePtr("source0"),
												},
											},
										},
									},
									Stores: &[]v1.Stores{
										{
											AdditionalStorageClasses: &[]string{"stc1", "stc2"},
											AllowInsecure:            pointer.MakePtr(true),
											Bucket:                   pointer.MakePtr("bucket-name"),
											ClusterName:              pointer.MakePtr("cluster-name"),
											ContainerName:            pointer.MakePtr("container-name"),
											DefaultFormat:            pointer.MakePtr("default-format"),
											Delimiter:                pointer.MakePtr("delimiter"),
											IncludeTags:              pointer.MakePtr(true),
											Name:                     pointer.MakePtr("store-name"),
											Prefix:                   pointer.MakePtr("prefix"),
											Provider:                 "AWS",
											Public:                   pointer.MakePtr(true),
											ReadConcern: &v1.ReadConcern{
												Level: pointer.MakePtr("local"),
											},
											ReadPreference: &v1.ReadPreference{
												Mode: pointer.MakePtr("primary"),
											},
											Region:               pointer.MakePtr("us-east-1"),
											ReplacementDelimiter: pointer.MakePtr("replacement-delimiter"),
											ServiceURL:           pointer.MakePtr("https://service-url.com"),
											Urls:                 &[]string{"url1", "url2"},
										},
									},
								},
							},
						},
					},
				}
				target := &admin2025.DataLakeTenant{}
				want := &admin2025.DataLakeTenant{
					Name: pointer.MakePtr("some-name"),
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
												AllowInsecure:       pointer.MakePtr(true),
												Collection:          pointer.MakePtr("some-name"),
												CollectionRegex:     pointer.MakePtr("collection-regex"),
												Database:            pointer.MakePtr("db"),
												DatabaseRegex:       pointer.MakePtr("db-regex"),
												DatasetName:         pointer.MakePtr("dataset-name"),
												DatasetPrefix:       pointer.MakePtr("dataset-prefix"),
												DefaultFormat:       pointer.MakePtr("default-format"),
												Path:                pointer.MakePtr("path"),
												ProvenanceFieldName: pointer.MakePtr("provenqance-field-name"),
												StoreName:           pointer.MakePtr("store-name"),
												TrimLevel:           pointer.MakePtr(1),
												Urls:                &[]string{"url1", "url2"},
											},
										},
										Name: pointer.MakePtr("collection0"),
									},
								},
								MaxWildcardCollections: pointer.MakePtr(3),
								Name:                   pointer.MakePtr("db0"),
								Views: &[]admin2025.DataLakeApiBase{
									{
										Name:     pointer.MakePtr("view0"),
										Pipeline: pointer.MakePtr("pipeline0"),
										Source:   pointer.MakePtr("source0"),
									},
								},
							},
						},
						Stores: &[]admin2025.DataLakeStoreSettings{
							{
								AdditionalStorageClasses: &[]string{"stc1", "stc2"},
								AllowInsecure:            pointer.MakePtr(true),
								Bucket:                   pointer.MakePtr("bucket-name"),
								ClusterName:              pointer.MakePtr("cluster-name"),
								ContainerName:            pointer.MakePtr("container-name"),
								DefaultFormat:            pointer.MakePtr("default-format"),
								Delimiter:                pointer.MakePtr("delimiter"),
								IncludeTags:              pointer.MakePtr(true),
								Name:                     pointer.MakePtr("store-name"),
								Prefix:                   pointer.MakePtr("prefix"),
								Provider:                 "AWS",
								Public:                   pointer.MakePtr(true),
								ReadConcern: &admin2025.DataLakeAtlasStoreReadConcern{
									Level: pointer.MakePtr("local"),
								},
								ReadPreference: &admin2025.DataLakeAtlasStoreReadPreference{
									Mode: pointer.MakePtr("primary"),
								},
								Region:               pointer.MakePtr("us-east-1"),
								ReplacementDelimiter: pointer.MakePtr("replacement-delimiter"),
								ServiceURL:           pointer.MakePtr("https://service-url.com"),
								Urls:                 &[]string{"url1", "url2"},
							},
						},
					},
				}
				testToAPI(t, "DataFederation", input, nil, target, want)
			},
		},

		{
			name: "sample database user",
			test: func(t *testing.T) {
				input := &v1.DatabaseUser{
					Spec: v1.DatabaseUserSpec{
						V20250312: &v1.DatabaseUserSpecV20250312{
							Entry: &v1.DatabaseUserSpecV20250312Entry{
								Username:     "test-user",
								DatabaseName: "admin",
								GroupId:      "32b6e34b3d91647abb20e7b8",
								Roles: &[]v1.Roles{
									{DatabaseName: "admin", RoleName: "readWrite"},
								},
								AwsIAMType:      pointer.MakePtr("aws-iam-type"),
								DeleteAfterDate: pointer.MakePtr("2025-07-01T00:00:00Z"),
								Description:     pointer.MakePtr("description"),
								Labels: &[]v1.Tags{
									{Key: "key-1", Value: "value-1"},
									{Key: "key-2", Value: "value-2"},
								},
								LdapAuthType: pointer.MakePtr("ldap-auth-type"),
								OidcAuthType: pointer.MakePtr("oidc-auth-type"),
								Password:     pointer.MakePtr("password"),
								Scopes: &[]v1.Scopes{
									{Name: "scope-1", Type: "type-1"},
									{Name: "scope-2", Type: "type-2"},
								},
								X509Type: pointer.MakePtr("x509-type"),
							},
							GroupId: pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.CloudDatabaseUser{}
				want := &admin2025.CloudDatabaseUser{
					Username:     "test-user",
					DatabaseName: "admin",
					GroupId:      "32b6e34b3d91647abb20e7b8",
					Roles: &[]admin2025.DatabaseUserRole{
						{DatabaseName: "admin", RoleName: "readWrite"},
					},
					AwsIAMType:      pointer.MakePtr("aws-iam-type"),
					DeleteAfterDate: pointer.MakePtr(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)),
					Description:     pointer.MakePtr("description"),
					Labels: &[]admin2025.ComponentLabel{
						{Key: pointer.MakePtr("key-1"), Value: pointer.MakePtr("value-1")},
						{Key: pointer.MakePtr("key-2"), Value: pointer.MakePtr("value-2")},
					},
					LdapAuthType: pointer.MakePtr("ldap-auth-type"),
					OidcAuthType: pointer.MakePtr("oidc-auth-type"),
					Password:     pointer.MakePtr("password"),
					Scopes: &[]admin2025.UserScope{
						{Name: "scope-1", Type: "type-1"},
						{Name: "scope-2", Type: "type-2"},
					},
					X509Type: pointer.MakePtr("x509-type"),
				}
				testToAPI(t, "DatabaseUser", input, nil, target, want)
			},
		},

		{
			name: "flex cluster with all fields",
			test: func(t *testing.T) {
				input := &v1.FlexCluster{
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
								TerminationProtectionEnabled: pointer.MakePtr(true),
							},
							GroupId: pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.FlexClusterDescriptionCreate20241113{}
				want := &admin2025.FlexClusterDescriptionCreate20241113{
					Name: "flex-cluster-name",
					ProviderSettings: admin2025.FlexProviderSettingsCreate20241113{
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: &[]admin2025.ResourceTag{
						{Key: "key1", Value: "value1"},
						{Key: "key2", Value: "value2"},
					},
					TerminationProtectionEnabled: pointer.MakePtr(true),
				}
				testToAPI(t, "FlexCluster", input, nil, target, want)
			},
		},

		{
			name: "simple group",
			test: func(t *testing.T) {
				input := &v1.Group{
					Spec: v1.GroupSpec{
						V20250312: &v1.GroupSpecV20250312{
							Entry: &v1.V20250312Entry{
								Name:                      "project-name",
								OrgId:                     "60987654321654321",
								RegionUsageRestrictions:   pointer.MakePtr("fake-restriction"),
								WithDefaultAlertsSettings: pointer.MakePtr(true),
								Tags: &[]v1.Tags{
									{Key: "key", Value: "value"},
								},
							},
							// read only field, not translated back to the API
							ProjectOwnerId: "61234567890123456",
						},
					},
				}
				target := &admin2025.Group{}
				want := &admin2025.Group{
					Name:                      "project-name",
					OrgId:                     "60987654321654321",
					RegionUsageRestrictions:   pointer.MakePtr("fake-restriction"),
					WithDefaultAlertsSettings: pointer.MakePtr(true),
					Tags: &[]admin2025.ResourceTag{
						{Key: "key", Value: "value"},
					},
				}
				testToAPI(t, "Group", input, nil, target, want)
			},
		},

		{
			name: "group alert config with project and credential references",
			test: func(t *testing.T) {
				input := &v1.GroupAlertsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
					},
					Spec: v1.GroupAlertsConfigSpec{
						V20250312: &v1.GroupAlertsConfigSpecV20250312{
							Entry: &v1.GroupAlertsConfigSpecV20250312Entry{
								Enabled:       pointer.MakePtr(true),
								EventTypeName: pointer.MakePtr("event-type"),
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
									Mode:       pointer.MakePtr("mode"),
									Operator:   pointer.MakePtr("operator"),
									Threshold:  pointer.MakePtr(1.1),
									Units:      pointer.MakePtr("units"),
								},
								Threshold: &v1.MetricThreshold{
									MetricName: "metric-t",
									Mode:       pointer.MakePtr("mode-t"),
									Operator:   pointer.MakePtr("operator-t"),
									Threshold:  pointer.MakePtr(2.2),
									Units:      pointer.MakePtr("units-t"),
								},
								Notifications: &[]v1.Notifications{
									{
										DatadogApiKeySecretRef: &v1.ApiTokenSecretRef{
											Name: pointer.MakePtr("datadog-secret"),
										},
										DatadogRegion: pointer.MakePtr("US"),
									},
								},
								SeverityOverride: pointer.MakePtr("some-severity-override"),
							},
							GroupId: pointer.MakePtr("60965432187654321"),
						},
					},
				}
				objs := []client.Object{
					&corev1.Secret{
						TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
						ObjectMeta: metav1.ObjectMeta{Name: "datadog-secret", Namespace: "ns"},
						Data: map[string][]byte{
							"datadogApiKey": ([]byte)("sample-password"),
						},
					},
				}
				target := &admin2025.GroupAlertsConfig{}
				want := &admin2025.GroupAlertsConfig{
					Enabled:       pointer.MakePtr(true),
					EventTypeName: pointer.MakePtr("event-type"),
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
						Mode:       pointer.MakePtr("mode"),
						Operator:   pointer.MakePtr("operator"),
						Threshold:  pointer.MakePtr(1.1),
						Units:      pointer.MakePtr("units"),
					},
					Threshold: &admin2025.StreamProcessorMetricThreshold{
						MetricName: pointer.MakePtr("metric-t"),
						Mode:       pointer.MakePtr("mode-t"),
						Operator:   pointer.MakePtr("operator-t"),
						Threshold:  pointer.MakePtr(2.2),
						Units:      pointer.MakePtr("units-t"),
					},
					Notifications: &[]admin2025.AlertsNotificationRootForGroup{
						{
							DatadogApiKey: pointer.MakePtr("sample-password"),
							DatadogRegion: pointer.MakePtr("US"),
						},
					},
					GroupId:          pointer.MakePtr("60965432187654321"),
					SeverityOverride: pointer.MakePtr("some-severity-override"),
				}
				testToAPI(t, "GroupAlertsConfig", input, objs, target, want)
			},
		},

		{
			name: "sample network peering connection",
			test: func(t *testing.T) {
				input := &v1.NetworkPeeringConnection{
					Spec: v1.NetworkPeeringConnectionSpec{
						V20250312: &v1.NetworkPeeringConnectionSpecV20250312{
							Entry: &v1.NetworkPeeringConnectionSpecV20250312Entry{
								AccepterRegionName:  pointer.MakePtr("accepter-region-name"),
								AwsAccountId:        pointer.MakePtr("aws-account-id"),
								AzureDirectoryId:    pointer.MakePtr("azure-dir-id"),
								AzureSubscriptionId: pointer.MakePtr("azure-subcription-id"),
								ContainerId:         "container-id",
								GcpProjectId:        pointer.MakePtr("azure-subcription-id"),
								NetworkName:         pointer.MakePtr("net-name"),
								ProviderName:        pointer.MakePtr("provider-name"),
								ResourceGroupName:   pointer.MakePtr("resource-group-name"),
								RouteTableCidrBlock: pointer.MakePtr("cidr"),
								VnetName:            pointer.MakePtr("vnet-name"),
								VpcId:               pointer.MakePtr("vpc-id"),
							},
							GroupId: pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.BaseNetworkPeeringConnectionSettings{}
				want := &admin2025.BaseNetworkPeeringConnectionSettings{
					ContainerId:         "container-id",
					ProviderName:        pointer.MakePtr("provider-name"),
					AccepterRegionName:  pointer.MakePtr("accepter-region-name"),
					AwsAccountId:        pointer.MakePtr("aws-account-id"),
					RouteTableCidrBlock: pointer.MakePtr("cidr"),
					VpcId:               pointer.MakePtr("vpc-id"),
					AzureDirectoryId:    pointer.MakePtr("azure-dir-id"),
					AzureSubscriptionId: pointer.MakePtr("azure-subcription-id"),
					ResourceGroupName:   pointer.MakePtr("resource-group-name"),
					VnetName:            pointer.MakePtr("vnet-name"),
					GcpProjectId:        pointer.MakePtr("azure-subcription-id"),
					NetworkName:         pointer.MakePtr("net-name"),
				}
				testToAPI(t, "NetworkPeeringConnection", input, nil, target, want)
			},
		},

		{
			name: "network permission entries all fields",
			test: func(t *testing.T) {
				input := &v1.NetworkPermissionEntries{
					Spec: v1.NetworkPermissionEntriesSpec{
						V20250312: &v1.NetworkPermissionEntriesSpecV20250312{
							Entry: &[]v1.NetworkPermissionEntriesSpecV20250312Entry{
								{
									AwsSecurityGroup: pointer.MakePtr("sg-12345678"),
									CidrBlock:        pointer.MakePtr("cird"),
									Comment:          pointer.MakePtr("comment"),
									DeleteAfterDate:  pointer.MakePtr("2025-07-01T00:00:00Z"),
									IpAddress:        pointer.MakePtr("1.1.1.1"),
								},
							},
							GroupId: pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &NetworkPermissions{}
				want := &NetworkPermissions{
					Entry: []admin2025.NetworkPermissionEntry{
						{
							AwsSecurityGroup: pointer.MakePtr("sg-12345678"),
							CidrBlock:        pointer.MakePtr("cird"),
							Comment:          pointer.MakePtr("comment"),
							DeleteAfterDate:  pointer.MakePtr(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)),
							IpAddress:        pointer.MakePtr("1.1.1.1"),
						},
					},
				}
				testToAPI(t, "NetworkPermissionEntries", input, nil, target, want)
			},
		},

		{
			name: "sample organization",
			test: func(t *testing.T) {
				input := &v1.Organization{
					Spec: v1.OrganizationSpec{
						V20250312: &v1.V20250312{
							Entry: &v1.Entry{
								ApiKey: &v1.ApiKey{
									Desc:  "description",
									Roles: []string{"role-1", "role-2"},
								},
								FederationSettingsId:      pointer.MakePtr("fed-id"),
								Name:                      "org-name",
								OrgOwnerId:                pointer.MakePtr("org-owner-id"),
								SkipDefaultAlertsSettings: pointer.MakePtr(true),
							},
						},
					},
				}
				target := &admin2025.AtlasOrganization{}
				want := &admin2025.AtlasOrganization{
					Name:                      "org-name",
					SkipDefaultAlertsSettings: pointer.MakePtr(true),
				}
				testToAPI(t, "Organization", input, nil, target, want)
			},
		},

		{
			name: "Organization setting with all fields",
			test: func(t *testing.T) {
				input := &v1.OrganizationSetting{
					Spec: v1.OrganizationSettingSpec{
						V20250312: &v1.OrganizationSettingSpecV20250312{
							Entry: &v1.OrganizationSettingSpecV20250312Entry{
								ApiAccessListRequired:                  pointer.MakePtr(true),
								GenAIFeaturesEnabled:                   pointer.MakePtr(true),
								MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
								MultiFactorAuthRequired:                pointer.MakePtr(true),
								RestrictEmployeeAccess:                 pointer.MakePtr(true),
								SecurityContact:                        pointer.MakePtr("contact-info"),
								StreamsCrossGroupEnabled:               pointer.MakePtr(true),
							},
							OrgId: "org-id",
						},
					},
				}
				target := &admin2025.OrganizationSettings{}
				want := &admin2025.OrganizationSettings{
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(true),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(true),
					SecurityContact:                        pointer.MakePtr("contact-info"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(true),
				}
				testToAPI(t, "OrganizationSetting", input, nil, target, want)
			},
		},

		{
			name: "customrole with all fields",
			test: func(t *testing.T) {
				input :=
					&v1.CustomRole{
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
					}
				target := &admin2025.UserCustomDBRole{}
				want := &admin2025.UserCustomDBRole{
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
				}
				testToAPI(t, "CustomRole", input, nil, target, want)
			},
		},

		{
			name: "searchindex create request fields",
			test: func(t *testing.T) {
				input := &v1.SearchIndex{
					Spec: v1.SearchIndexSpec{
						V20250312: &v1.SearchIndexSpecV20250312{
							Entry: &v1.SearchIndexSpecV20250312Entry{
								Database:       "database-name",
								CollectionName: "collection-name",
								Name:           "index-name",
								Type:           pointer.MakePtr("search-index-type"),
								Definition: &v1.Definition{
									Analyzer: pointer.MakePtr("lucene.standard"),
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
										Dynamic: pointer.MakePtr(true),
										Fields: &map[string]apiextensionsv1.JSON{
											"field1": {Raw: []byte(`{"key4":"value4"}`)},
										},
									},
									NumPartitions:  pointer.MakePtr(3),
									SearchAnalyzer: pointer.MakePtr("lucene.standard"),
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
							GroupId: pointer.MakePtr("group-id-101"),
						},
					},
				}
				target := &admin2025.SearchIndexCreateRequest{}
				want := &admin2025.SearchIndexCreateRequest{
					CollectionName: "collection-name",
					Database:       "database-name",
					Name:           "index-name",
					Type:           pointer.MakePtr("search-index-type"),
					Definition: &admin2025.BaseSearchIndexCreateRequestDefinition{
						Analyzer: pointer.MakePtr("lucene.standard"),
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
							Dynamic: pointer.MakePtr(true),
							Fields: &map[string]any{
								"field1": map[string]any{"key4": "value4"},
							},
						},
						NumPartitions:  pointer.MakePtr(3),
						SearchAnalyzer: pointer.MakePtr("lucene.standard"),
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
				}
				testToAPI(t, "SearchIndex", input, nil, target, want)
			},
		},

		{
			name: "team all fields",
			test: func(t *testing.T) {
				input := &v1.Team{
					Spec: v1.TeamSpec{
						V20250312: &v1.TeamSpecV20250312{
							Entry: &v1.TeamSpecV20250312Entry{
								Name:      "team-name",
								Usernames: []string{"user1", "user2"},
							},
							OrgId: "org-id",
						},
					},
				}
				target := &admin2025.Team{}
				want := &admin2025.Team{
					Name: "team-name",
					Usernames: []string{
						"user1", "user2",
					},
				}
				testToAPI(t, "Team", input, nil, target, want)
			},
		},

		{
			name: "third party integration all fields",
			test: func(t *testing.T) {
				input := &v1.ThirdPartyIntegration{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
					},
					Spec: v1.ThirdPartyIntegrationSpec{
						V20250312: &v1.ThirdPartyIntegrationSpecV20250312{
							IntegrationType: "ANY",
							Entry: &v1.ThirdPartyIntegrationSpecV20250312Entry{
								AccountId: pointer.MakePtr("account-id"),
								ApiKeySecretRef: &v1.ApiTokenSecretRef{
									Key:  pointer.MakePtr("apiKey"),
									Name: pointer.MakePtr("multi-secret0"),
								},
								ApiTokenSecretRef: &v1.ApiTokenSecretRef{
									Key:  pointer.MakePtr("apiToken"),
									Name: pointer.MakePtr("multi-secret0"),
								},
								ChannelName: pointer.MakePtr("channel-name"),
								Enabled:     pointer.MakePtr(true),
								LicenseKeySecretRef: &v1.ApiTokenSecretRef{
									Key:  pointer.MakePtr("licenseKey"),
									Name: pointer.MakePtr("multi-secret1"),
								},
								Region:                       pointer.MakePtr("some-region"),
								SendCollectionLatencyMetrics: pointer.MakePtr(true),
								SendDatabaseMetrics:          pointer.MakePtr(true),
								SendUserProvidedResourceTags: pointer.MakePtr(true),
								ServiceDiscovery:             pointer.MakePtr("service-discovery"),
								TeamName:                     pointer.MakePtr("some-team"),
								Type:                         pointer.MakePtr("some-type"),
								Username:                     pointer.MakePtr("username"),
							},
							GroupId: pointer.MakePtr("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.ThirdPartyIntegration{}
				want := &admin2025.ThirdPartyIntegration{
					Type:                         pointer.MakePtr("some-type"),
					ApiKey:                       pointer.MakePtr("sample-api-key"),
					Region:                       pointer.MakePtr("some-region"),
					SendCollectionLatencyMetrics: pointer.MakePtr(true),
					SendDatabaseMetrics:          pointer.MakePtr(true),
					SendUserProvidedResourceTags: pointer.MakePtr(true),
					AccountId:                    pointer.MakePtr("account-id"),
					LicenseKey:                   pointer.MakePtr("sample-license-key"),
					Enabled:                      pointer.MakePtr(true),
					ServiceDiscovery:             pointer.MakePtr("service-discovery"),
					Username:                     pointer.MakePtr("username"),
					ApiToken:                     pointer.MakePtr("sample-api-token"),
					ChannelName:                  pointer.MakePtr("channel-name"),
					TeamName:                     pointer.MakePtr("some-team"),
				}
				objs := []client.Object{
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
				}
				testToAPI(t, "ThirdPartyIntegration", input, objs, target, want)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func testToAPI[T any](t *testing.T, kind string, input client.Object, objs []client.Object, target, want *T) {
	crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
	crd, err := extractCRD(kind, bufio.NewScanner(crdsYML))
	require.NoError(t, err)
	trs, err := crapi.NewPerVersionTranslators(crd, version, sdkVersion)
	require.NoError(t, err)
	tr := trs[sdkVersion]
	require.NotNil(t, tr)
	r := crapi.Request{Translator: tr, Dependencies: objs}
	require.NoError(t, crapi.ToAPI(&r, target, input))
	assert.Equal(t, want, target)
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
