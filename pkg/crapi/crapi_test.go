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

	crd2gok8s "github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin2025 "go.mongodb.org/atlas-sdk/v20250312018/admin"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/crds"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/testdata"
	samplesv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/testdata/samples/v1"
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
					Id:           new("6127378123219"),
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
					WithDefaultAlertsSettings: new(true),
				}
				target := samplesv1.Group{
					Spec: samplesv1.GroupSpec{
						V20250312: &samplesv1.GroupSpecV20250312{
							ProjectOwnerId: new(""),
						},
					},
				}
				want := &samplesv1.Group{
					Spec: samplesv1.GroupSpec{
						V20250312: &samplesv1.GroupSpecV20250312{
							Entry: &samplesv1.GroupSpecV20250312Entry{
								Name:  "test-project",
								OrgId: testProjectID,
								Tags: &[]samplesv1.Tags{
									{
										Key:   "key0",
										Value: "value0",
									},
									{
										Key:   "key",
										Value: "value",
									},
								},
								WithDefaultAlertsSettings: new(true),
							},
							ProjectOwnerId: new(""),
						},
					},
					Status: samplesv1.GroupStatus{
						V20250312: &samplesv1.GroupStatusV20250312{
							Created: "2025-01-01T01:30:15Z",
							Id:      new("6127378123219"),
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
					AwsIAMType:      new("NONE AWS"),
					DatabaseName:    "dbname",
					DeleteAfterDate: new(time.Date(2025, 2, 1, 1, 30, 15, 0, time.UTC)),
					Description:     new("sample db user"),
					GroupId:         testProjectID,
					Labels: &[]admin2025.ComponentLabel{
						{
							Key:   new("key0"),
							Value: new("value0"),
						},
						{
							Key:   new("key1"),
							Value: new("value1"),
						},
					},
					LdapAuthType: new("NONE LDAP"),
					OidcAuthType: new("NONE OIDC"),
					// TODO: new crd should put this on a secret
					Password: new("fakepass"),
					Roles: []admin2025.DatabaseUserRole{
						{
							CollectionName: new("collection0"),
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
					X509Type: new("NONE X509"),
				}
				target := samplesv1.DatabaseUser{}
				want := &samplesv1.DatabaseUser{
					Spec: samplesv1.DatabaseUserSpec{
						V20250312: &samplesv1.DatabaseUserSpecV20250312{
							Entry: &samplesv1.DatabaseUserSpecV20250312Entry{
								AwsIAMType:      new("NONE AWS"),
								DatabaseName:    "dbname",
								DeleteAfterDate: new("2025-02-01T01:30:15Z"),
								Description:     new("sample db user"),
								Labels: &[]samplesv1.Tags{
									{
										Key:   "key0",
										Value: "value0",
									},
									{
										Key:   "key1",
										Value: "value1",
									},
								},
								LdapAuthType: new("NONE LDAP"),
								OidcAuthType: new("NONE OIDC"),
								PasswordSecretRef: &samplesv1.PasswordSecretRef{
									Key:  new("password"),
									Name: "-6cb55bffddcfffc5d4c",
								},
								Roles: &[]samplesv1.Roles{
									{
										CollectionName: new("collection0"),
										DatabaseName:   "mydb",
										RoleName:       "admin",
									},
								},
								Scopes: &[]samplesv1.Scopes{
									{
										Name: "scopeName",
										Type: "scopeType",
									},
								},
								Username: "dbuser",
								X509Type: new("NONE X509"),
							},
							GroupId: new(testProjectID),
						},
					},
				}
				wantSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: "-6cb55bffddcfffc5d4c",
					},
					Data: map[string][]byte{
						"password": []byte("fakepass"),
					},
				}
				testFromAPI(t, "DatabaseUser", &target, &input, want, wantSecret)
			},
		},

		{
			name: "GroupAlertConfigs",
			test: func(t *testing.T) {
				input := admin2025.GroupAlertsConfig{
					Enabled:       new(true),
					EventTypeName: new("OUTSIDE_STREAM_PROCESSOR_METRIC_THRESHOLD"),
					GroupId:       new(testProjectID),
					Id:            new("notification id"),
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
							DatadogApiKey: new("fake api key"),
							DatadogRegion: new("US"),
							DelayMin:      new(42),
							IntegrationId: new("32b6e34b3d91647abb20e7b8"),
							IntervalMin:   new(43),
							NotifierId:    new("32b6e34b3d91647abb20e7b8"),
							TypeName:      new("DATADOG"),
						},
					},
					SeverityOverride: new("CRITICIAL"),
					MetricThreshold: &admin2025.FlexClusterMetricThreshold{
						MetricName: "metric",
						Mode:       new("mode"),
						Operator:   new("op"),
						Threshold:  new(0.1),
						Units:      new("unit"),
					},
					Threshold: &admin2025.StreamProcessorMetricThreshold{
						MetricName: new("anotherMetric"),
						Mode:       new("a mode"),
						Operator:   new("an op"),
						Threshold:  new(0.2),
						Units:      new("a unit"),
					},
				}
				target := samplesv1.GroupAlertsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "groupalertscfg",
						Namespace: "ns",
					},
				}
				want := &samplesv1.GroupAlertsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "groupalertscfg",
						Namespace: "ns",
					},
					Spec: samplesv1.GroupAlertsConfigSpec{
						V20250312: &samplesv1.GroupAlertsConfigSpecV20250312{
							Entry: &samplesv1.GroupAlertsConfigSpecV20250312Entry{
								Enabled:       new(true),
								EventTypeName: new("OUTSIDE_STREAM_PROCESSOR_METRIC_THRESHOLD"),
								Matchers: &[]samplesv1.Matchers{
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
								MetricThreshold: &samplesv1.MetricThreshold{
									MetricName: "metric",
									Mode:       new("mode"),
									Operator:   new("op"),
									Threshold:  new(0.1),
									Units:      new("unit"),
								},
								Notifications: &[]samplesv1.Notifications{
									{
										DatadogApiKeySecretRef: &samplesv1.PasswordSecretRef{
											Key:  new("datadogApiKey"),
											Name: "groupalertscfg-f4f4b5f9c849fc4cbdc",
										},
										DatadogRegion: new("US"),
										DelayMin:      new(42),
										IntegrationId: new("32b6e34b3d91647abb20e7b8"),
										IntervalMin:   new(43),
										NotifierId:    new("32b6e34b3d91647abb20e7b8"),
										TypeName:      new("DATADOG"),
									},
								},
								SeverityOverride: new("CRITICIAL"),
								Threshold: &samplesv1.MetricThreshold{
									MetricName: "anotherMetric",
									Mode:       new("a mode"),
									Operator:   new("an op"),
									Threshold:  new(0.2),
									Units:      new("a unit"),
								},
							},
							GroupId: new(testProjectID),
						},
					},
					Status: samplesv1.GroupAlertsConfigStatus{
						V20250312: &samplesv1.GroupAlertsConfigStatusV20250312{
							GroupId: new(testProjectID),
							Id:      new("notification id"),
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
					Id:          new("SomeID"),
					Type:        new("SLACK"),
					ApiToken:    new("some fake api token"),
					ChannelName: new("alert-channel"),
					TeamName:    new("some-team"),
				}
				target := samplesv1.ThirdPartyIntegration{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "3rdparty-slack",
						Namespace: "ns",
					},
					Spec: samplesv1.ThirdPartyIntegrationSpec{
						V20250312: &samplesv1.ThirdPartyIntegrationSpecV20250312{
							// TODO: is this a valid trick?
							// This API struct, unlike others, does NOT include the Group ID
							// it is part of the parameters, but not the response
							GroupId: new(testProjectID),
							// TODO: similarly to the Group ID the IntegrationType would
							// be the aparameter thet corresponds with "type" in the response
							// but there is no indication of such semantics from the CRD
							IntegrationType: "SLACK",
						},
					},
				}
				want := &samplesv1.ThirdPartyIntegration{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "3rdparty-slack",
						Namespace: "ns",
					},
					Spec: samplesv1.ThirdPartyIntegrationSpec{
						V20250312: &samplesv1.ThirdPartyIntegrationSpecV20250312{
							Entry: &samplesv1.ThirdPartyIntegrationSpecV20250312Entry{
								Type: new("SLACK"),
								ApiTokenSecretRef: &samplesv1.PasswordSecretRef{
									Name: "3rdparty-slack-5798d555ff4dc66f7c99",
									Key:  new("apiToken"),
								},
								ChannelName: new("alert-channel"),
								TeamName:    new("some-team"),
							},
							// Pre-existing from input
							GroupId: new(string(testProjectID)),
							// Pre-existing from input
							IntegrationType: "SLACK",
						},
					},
					Status: samplesv1.ThirdPartyIntegrationStatus{
						V20250312: &samplesv1.ThirdPartyIntegrationStatusV20250312{
							Id:   new("SomeID"),
							Type: new("SLACK"),
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

func testFromAPI[S any](t *testing.T, kind string, target client.Object, input *S, want client.Object, wantDeps ...client.Object) {
	scheme := testScheme(t)
	crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
	crd, err := extractCRD(kind, bufio.NewScanner(crdsYML))
	require.NoError(t, err)
	tr, err := crapi.NewTranslator(scheme, crd, version, sdkVersion)
	require.NoError(t, err)
	results, err := tr.FromAPI(target, input)
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
			input: &samplesv1.GroupAlertsConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "GroupAlertsConfig",
					APIVersion: "atlas.generated.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-group-alerts-config",
					Namespace: "ns",
				},
				Spec: samplesv1.GroupAlertsConfigSpec{
					V20250312: &samplesv1.GroupAlertsConfigSpecV20250312{
						Entry: &samplesv1.GroupAlertsConfigSpecV20250312Entry{
							Enabled:       new(true),
							EventTypeName: new("some-event"),
							Matchers: &[]samplesv1.Matchers{
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
							MetricThreshold: &samplesv1.MetricThreshold{
								MetricName: "metric",
								Mode:       new("mode"),
								Operator:   new("operator"),
								Threshold:  new(1.0),
								Units:      new("unit"),
							},
							Notifications: &[]samplesv1.Notifications{
								{
									DatadogApiKeySecretRef: &samplesv1.PasswordSecretRef{
										Name: "alert-secrets-0",
										Key:  new("apiKey"),
									},
									DatadogRegion: new("US"),
								},
								{
									WebhookSecretSecretRef: &samplesv1.PasswordSecretRef{
										Name: "alert-secrets-0",
										Key:  new("webhookSecret"),
									},
									WebhookUrlSecretRef: &samplesv1.PasswordSecretRef{
										Name: "alert-secrets-1",
										Key:  new("webhookUrl"),
									},
								},
							},
							SeverityOverride: new("severe"),
							Threshold: &samplesv1.MetricThreshold{
								MetricName: "metric",
								Mode:       new("mode-t"),
								Operator:   new("op-t"),
								Threshold:  new(2.0),
								Units:      new("unit-t"),
							},
						},
						GroupRef: &crd2gok8s.LocalReference{
							Name: "my-project",
						},
					},
				},
			},
			deps: []client.Object{
				&samplesv1.Group{
					TypeMeta:   metav1.TypeMeta{Kind: "Group", APIVersion: "atlas.generated.mongodb.com/v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "my-project", Namespace: "ns"},
					Spec: samplesv1.GroupSpec{
						V20250312: &samplesv1.GroupSpecV20250312{
							Entry: &samplesv1.GroupSpecV20250312Entry{
								Name:  "some-project",
								OrgId: "621454123423x125235142",
							},
						},
					},
					Status: samplesv1.GroupStatus{
						V20250312: &samplesv1.GroupStatusV20250312{
							Id: new("62b6e34b3d91647abb20e7b8"),
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
				&corev1.Secret{ // works without type meta set as well
					ObjectMeta: metav1.ObjectMeta{Name: "alert-secrets-1", Namespace: "ns"},
					Data: map[string][]byte{
						"webhookUrl": ([]byte)("sample-webhook-url"),
					},
				},
			},
			// nolint:dupl
			want: admin2025.GroupAlertsConfig{
				Enabled:       new(true),
				EventTypeName: new("some-event"),
				GroupId:       new("62b6e34b3d91647abb20e7b8"),
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
						DatadogApiKey: new("sample-api-key"),
						DatadogRegion: new("US"),
					},
					{
						WebhookSecret: new("sample-webhook-secret"),
						WebhookUrl:    new("sample-webhook-url"),
					},
				},
				SeverityOverride: new("severe"),
				MetricThreshold: &admin2025.FlexClusterMetricThreshold{
					MetricName: "metric",
					Mode:       new("mode"),
					Operator:   new("operator"),
					Threshold:  new(1.0),
					Units:      new("unit"),
				},
				Threshold: &admin2025.StreamProcessorMetricThreshold{
					MetricName: new("metric"),
					Mode:       new("mode-t"),
					Operator:   new("op-t"),
					Threshold:  new(2.0),
					Units:      new("unit-t"),
				},
			},
			target: admin2025.GroupAlertsConfig{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scheme := testScheme(t)
			crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
			crd, err := extractCRD(tc.crd, bufio.NewScanner(crdsYML))
			require.NoError(t, err)
			tr, err := crapi.NewTranslator(scheme, crd, version, sdkVersion)
			require.NoError(t, err)
			require.NoError(t, tr.ToAPI(&tc.target, tc.input, tc.deps...))
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
				input := &samplesv1.BackupCompliancePolicy{
					Spec: samplesv1.BackupCompliancePolicySpec{
						V20250312: &samplesv1.BackupCompliancePolicySpecV20250312{
							Entry: &samplesv1.BackupCompliancePolicySpecV20250312Entry{
								AuthorizedEmail:         "user@example.com",
								CopyProtectionEnabled:   new(true),
								EncryptionAtRestEnabled: new(true),
								AuthorizedUserFirstName: "first-name",
								AuthorizedUserLastName:  "last-name",
								OnDemandPolicyItem: &samplesv1.OnDemandPolicyItem{
									FrequencyInterval: 1,
									FrequencyType:     "some-freq",
									RetentionUnit:     "some-unit",
									RetentionValue:    2,
								},
								PitEnabled:        new(true),
								ProjectId:         new("project-id"),
								RestoreWindowDays: new(3),
								ScheduledPolicyItems: &[]samplesv1.OnDemandPolicyItem{
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
							GroupId:                 new("32b6e34b3d91647abb20e7b8"),
							OverwriteBackupPolicies: new(true),
						},
					},
				}
				target := &admin2025.DataProtectionSettings20231001{}
				want := &admin2025.DataProtectionSettings20231001{
					AuthorizedEmail:         "user@example.com",
					CopyProtectionEnabled:   new(true),
					EncryptionAtRestEnabled: new(true),
					AuthorizedUserFirstName: "first-name",
					AuthorizedUserLastName:  "last-name",
					OnDemandPolicyItem: &admin2025.BackupComplianceOnDemandPolicyItem{
						FrequencyInterval: 1,
						FrequencyType:     "some-freq",
						RetentionUnit:     "some-unit",
						RetentionValue:    2,
					},
					PitEnabled:        new(true),
					ProjectId:         new("project-id"),
					RestoreWindowDays: new(3),
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
				input := &samplesv1.BackupSchedule{
					Spec: samplesv1.BackupScheduleSpec{
						V20250312: &samplesv1.BackupScheduleSpecV20250312{
							Entry: &samplesv1.BackupScheduleSpecV20250312Entry{
								ReferenceHourOfDay:    new(2),
								ReferenceMinuteOfHour: new(30),
								RestoreWindowDays:     new(7),
								UpdateSnapshots:       new(true),
								AutoExportEnabled:     new(true),
								CopySettings: &[]samplesv1.CopySettings{
									{
										CloudProvider:    new("AWS"),
										Frequencies:      &[]string{"freq-1", "freq-2"},
										RegionName:       new("us-east-1"),
										ShouldCopyOplogs: new(true),
										ZoneId:           "zone-id",
									},
									{
										CloudProvider:    new("GCE"),
										Frequencies:      &[]string{"freq-3", "freq-4"},
										RegionName:       new("us-east-3"),
										ShouldCopyOplogs: new(true),
										ZoneId:           "zone-id-0",
									},
								},
								DeleteCopiedBackups: &[]samplesv1.DeleteCopiedBackups{
									{
										CloudProvider: new("Azure"),
										RegionName:    new("us-west-2"),
										ZoneId:        new("zone-id"),
									},
								},
								Export: &samplesv1.Export{
									ExportBucketId: new("ExportBucketId"),
									FrequencyType:  new("FrequencyType"),
								},
								ExtraRetentionSettings: &[]samplesv1.ExtraRetentionSettings{
									{
										FrequencyType: new("FrequencyType0"),
										RetentionDays: new(1),
									},
									{
										FrequencyType: new("FrequencyType1"),
										RetentionDays: new(2),
									},
								},
								Policies: &[]samplesv1.Policies{
									{
										Id: new("id0"),
										PolicyItems: &[]samplesv1.OnDemandPolicyItem{
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
								UseOrgAndGroupNamesInExportPrefix: new(true),
							},
							GroupId:     new("group-id-101"),
							ClusterName: "cluster-name",
						},
					},
				}
				target := &admin2025.DiskBackupSnapshotSchedule20240805{}
				want := &admin2025.DiskBackupSnapshotSchedule20240805{
					ReferenceHourOfDay:    new(2),
					ReferenceMinuteOfHour: new(30),
					RestoreWindowDays:     new(7),
					UpdateSnapshots:       new(true),
					AutoExportEnabled:     new(true),
					CopySettings: &[]admin2025.DiskBackupCopySetting20240805{
						{
							CloudProvider:    new("AWS"),
							Frequencies:      &[]string{"freq-1", "freq-2"},
							RegionName:       new("us-east-1"),
							ShouldCopyOplogs: new(true),
							ZoneId:           "zone-id",
						},
						{
							CloudProvider:    new("GCE"),
							Frequencies:      &[]string{"freq-3", "freq-4"},
							RegionName:       new("us-east-3"),
							ShouldCopyOplogs: new(true),
							ZoneId:           "zone-id-0",
						},
					},
					DeleteCopiedBackups: &[]admin2025.DeleteCopiedBackups20240805{
						{
							CloudProvider: new("Azure"),
							RegionName:    new("us-west-2"),
							ZoneId:        new("zone-id"),
						},
					},
					Export: &admin2025.AutoExportPolicy{
						ExportBucketId: new("ExportBucketId"),
						FrequencyType:  new("FrequencyType"),
					},
					ExtraRetentionSettings: &[]admin2025.ExtraRetentionSetting{
						{
							FrequencyType: new("FrequencyType0"),
							RetentionDays: new(1),
						},
						{
							FrequencyType: new("FrequencyType1"),
							RetentionDays: new(2),
						},
					},
					Policies: []admin2025.AdvancedDiskBackupSnapshotSchedulePolicy{
						{
							Id: new("id0"),
							PolicyItems: []admin2025.DiskBackupApiPolicyItem{
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
					UseOrgAndGroupNamesInExportPrefix: new(true),
					ClusterName:                       new("cluster-name"),
				}
				testToAPI(t, "BackupSchedule", input, nil, target, want)
			},
		},

		//nolint:dupl
		{
			name: "cluster all fields",
			test: func(t *testing.T) {
				input := &samplesv1.Cluster{
					Spec: samplesv1.ClusterSpec{
						V20250312: &samplesv1.ClusterSpecV20250312{
							Entry: &samplesv1.ClusterSpecV20250312Entry{
								AcceptDataRisksAndForceReplicaSetReconfig: new("2025-01-01T00:00:00Z"),
								AdvancedConfiguration: &samplesv1.AdvancedConfiguration{
									CustomOpensslCipherConfigTls12: &[]string{
										"TLS_AES_256_GCM_SHA384", "TLS_CHACHA20_POLY1305_SHA256",
									},
									MinimumEnabledTlsProtocol: new("TLS1.2"),
									TlsCipherConfigMode:       new("Custom"),
								},
								BackupEnabled:                    new(true),
								BiConnector:                      &samplesv1.BiConnector{Enabled: new(true)},
								ClusterType:                      new("ReplicaSet"),
								ConfigServerManagementMode:       new("Managed"),
								ConfigServerType:                 new("ReplicaSet"),
								DiskWarmingMode:                  new("Enabled"),
								EncryptionAtRestProvider:         new("AWS-KMS"),
								GlobalClusterSelfManagedSharding: new(true),
								Labels: &[]samplesv1.Tags{
									{Key: "key1", Value: "value1"},
									{Key: "key2", Value: "value2"},
								},
								MongoDBEmployeeAccessGrant: &samplesv1.MongoDBEmployeeAccessGrant{
									ExpirationTime: "2025-12-31T00:00:00Z",
									GrantType:      "Temporary",
								},
								MongoDBMajorVersion:       new("8.0"),
								Name:                      new("my-cluster"),
								Paused:                    new(true),
								PitEnabled:                new(true),
								RedactClientLogData:       new(true),
								ReplicaSetScalingStrategy: new("Auto"),
								ReplicationSpecs: &[]samplesv1.ReplicationSpecs{
									{
										ZoneId:   new("zone-id-1"),
										ZoneName: new("zone-name-1"),
										RegionConfigs: &[]samplesv1.RegionConfigs{
											{
												RegionName: new("us-east-1"),
												AnalyticsSpecs: &samplesv1.AnalyticsSpecs{
													DiskIOPS:      new(1000),
													DiskSizeGB:    new(10.0),
													EbsVolumeType: new("gp2"),
													InstanceSize:  new("M10"),
													NodeCount:     new(3),
												},
												AutoScaling: &samplesv1.AnalyticsAutoScaling{
													Compute: &samplesv1.Compute{
														Enabled:          new(true),
														ScaleDownEnabled: new(true),
														MaxInstanceSize:  new("M20"),
														MinInstanceSize:  new("M10"),
													},
													DiskGB: &samplesv1.DiskGB{
														Enabled: new(true),
													},
												},
												AnalyticsAutoScaling: &samplesv1.AnalyticsAutoScaling{
													Compute: &samplesv1.Compute{
														Enabled:          new(true),
														ScaleDownEnabled: new(true),
														MaxInstanceSize:  new("M30"),
														MinInstanceSize:  new("M10"),
													},
													DiskGB: &samplesv1.DiskGB{
														Enabled: new(true),
													},
												},
												BackingProviderName: new("AWS"),
												ElectableSpecs: &samplesv1.ElectableSpecs{
													DiskIOPS:              new(1000),
													DiskSizeGB:            new(10.0),
													EbsVolumeType:         new("gp2"),
													EffectiveInstanceSize: new("M10"),
													InstanceSize:          new("M10"),
													NodeCount:             new(3),
												},
												Priority:     new(1),
												ProviderName: new("AWS"),
												ReadOnlySpecs: &samplesv1.AnalyticsSpecs{
													DiskIOPS:      new(1000),
													DiskSizeGB:    new(10.0),
													EbsVolumeType: new("gp2"),
													InstanceSize:  new("M10"),
													NodeCount:     new(3),
												},
											},
											{
												RegionName: new("us-east-2"),
												AnalyticsSpecs: &samplesv1.AnalyticsSpecs{
													DiskIOPS:      new(2000),
													DiskSizeGB:    new(10.0),
													EbsVolumeType: new("gp3"),
													InstanceSize:  new("M20"),
													NodeCount:     new(3),
												},
												AutoScaling: &samplesv1.AnalyticsAutoScaling{
													Compute: &samplesv1.Compute{
														Enabled:          new(true),
														ScaleDownEnabled: new(true),
														MaxInstanceSize:  new("M50"),
														MinInstanceSize:  new("M20"),
													},
													DiskGB: &samplesv1.DiskGB{
														Enabled: new(true),
													},
												},
												AnalyticsAutoScaling: &samplesv1.AnalyticsAutoScaling{
													Compute: &samplesv1.Compute{
														Enabled:          new(true),
														ScaleDownEnabled: new(true),
														MaxInstanceSize:  new("M40"),
														MinInstanceSize:  new("M10"),
													},
													DiskGB: &samplesv1.DiskGB{
														Enabled: new(true),
													},
												},
												BackingProviderName: new("AWS"),
												ElectableSpecs: &samplesv1.ElectableSpecs{
													DiskIOPS:              new(1000),
													DiskSizeGB:            new(10.0),
													EbsVolumeType:         new("gp2"),
													EffectiveInstanceSize: new("M10"),
													InstanceSize:          new("M10"),
													NodeCount:             new(3),
												},
												Priority:     new(1),
												ProviderName: new("AWS"),
												ReadOnlySpecs: &samplesv1.AnalyticsSpecs{
													DiskIOPS:      new(1000),
													DiskSizeGB:    new(10.0),
													EbsVolumeType: new("gp2"),
													InstanceSize:  new("M10"),
													NodeCount:     new(3),
												},
											},
										},
									},
								},
								RootCertType: new("X509"),
								Tags: &[]samplesv1.Tags{
									{Key: "key1", Value: "value1"},
									{Key: "key2", Value: "value2"},
								},
								TerminationProtectionEnabled: new(true),
								VersionReleaseSystem:         new("Atlas"),
							},
							GroupId: new("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.ClusterDescription20240805{}
				want := &admin2025.ClusterDescription20240805{
					AcceptDataRisksAndForceReplicaSetReconfig: new(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					AdvancedConfiguration: &admin2025.ApiAtlasClusterAdvancedConfiguration{
						CustomOpensslCipherConfigTls12: &[]string{
							"TLS_AES_256_GCM_SHA384", "TLS_CHACHA20_POLY1305_SHA256",
						},
						MinimumEnabledTlsProtocol: new("TLS1.2"),
						TlsCipherConfigMode:       new("Custom"),
					},
					BackupEnabled:                    new(true),
					BiConnector:                      &admin2025.BiConnector{Enabled: new(true)},
					ClusterType:                      new("ReplicaSet"),
					ConfigServerManagementMode:       new("Managed"),
					ConfigServerType:                 new("ReplicaSet"),
					DiskWarmingMode:                  new("Enabled"),
					EncryptionAtRestProvider:         new("AWS-KMS"),
					GlobalClusterSelfManagedSharding: new(true),
					Labels: &[]admin2025.ComponentLabel{
						{Key: new("key1"), Value: new("value1")},
						{Key: new("key2"), Value: new("value2")},
					},
					MongoDBEmployeeAccessGrant: &admin2025.EmployeeAccessGrant{
						ExpirationTime: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
						GrantType:      "Temporary",
					},
					MongoDBMajorVersion:       new("8.0"),
					Name:                      new("my-cluster"),
					Paused:                    new(true),
					PitEnabled:                new(true),
					RedactClientLogData:       new(true),
					ReplicaSetScalingStrategy: new("Auto"),
					ReplicationSpecs: &[]admin2025.ReplicationSpec20240805{
						{
							ZoneId:   new("zone-id-1"),
							ZoneName: new("zone-name-1"),
							RegionConfigs: &[]admin2025.CloudRegionConfig20240805{
								{
									RegionName: new("us-east-1"),
									AnalyticsSpecs: &admin2025.DedicatedHardwareSpec20240805{
										DiskIOPS:      new(1000),
										DiskSizeGB:    new(10.0),
										EbsVolumeType: new("gp2"),
										InstanceSize:  new("M10"),
										NodeCount:     new(3),
									},
									AutoScaling: &admin2025.AdvancedAutoScalingSettings{
										Compute: &admin2025.AdvancedComputeAutoScaling{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
											MaxInstanceSize:  new("M20"),
											MinInstanceSize:  new("M10"),
										},
										DiskGB: &admin2025.DiskGBAutoScaling{
											Enabled: new(true),
										},
									},
									AnalyticsAutoScaling: &admin2025.AdvancedAutoScalingSettings{
										Compute: &admin2025.AdvancedComputeAutoScaling{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
											MaxInstanceSize:  new("M30"),
											MinInstanceSize:  new("M10"),
										},
										DiskGB: &admin2025.DiskGBAutoScaling{
											Enabled: new(true),
										},
									},
									BackingProviderName: new("AWS"),
									ElectableSpecs: &admin2025.HardwareSpec20240805{
										DiskIOPS:              new(1000),
										DiskSizeGB:            new(10.0),
										EbsVolumeType:         new("gp2"),
										EffectiveInstanceSize: new("M10"),
										InstanceSize:          new("M10"),
										NodeCount:             new(3),
									},
									Priority:     new(1),
									ProviderName: new("AWS"),
									ReadOnlySpecs: &admin2025.DedicatedHardwareSpec20240805{
										DiskIOPS:      new(1000),
										DiskSizeGB:    new(10.0),
										EbsVolumeType: new("gp2"),
										InstanceSize:  new("M10"),
										NodeCount:     new(3),
									},
								},
								{
									RegionName: new("us-east-2"),
									AnalyticsSpecs: &admin2025.DedicatedHardwareSpec20240805{
										DiskIOPS:      new(2000),
										DiskSizeGB:    new(10.0),
										EbsVolumeType: new("gp3"),
										InstanceSize:  new("M20"),
										NodeCount:     new(3),
									},
									AutoScaling: &admin2025.AdvancedAutoScalingSettings{
										Compute: &admin2025.AdvancedComputeAutoScaling{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
											MaxInstanceSize:  new("M50"),
											MinInstanceSize:  new("M20"),
										},
										DiskGB: &admin2025.DiskGBAutoScaling{
											Enabled: new(true),
										},
									},
									AnalyticsAutoScaling: &admin2025.AdvancedAutoScalingSettings{
										Compute: &admin2025.AdvancedComputeAutoScaling{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
											MaxInstanceSize:  new("M40"),
											MinInstanceSize:  new("M10"),
										},
										DiskGB: &admin2025.DiskGBAutoScaling{
											Enabled: new(true),
										},
									},
									BackingProviderName: new("AWS"),
									ElectableSpecs: &admin2025.HardwareSpec20240805{
										DiskIOPS:              new(1000),
										DiskSizeGB:            new(10.0),
										EbsVolumeType:         new("gp2"),
										EffectiveInstanceSize: new("M10"),
										InstanceSize:          new("M10"),
										NodeCount:             new(3),
									},
									Priority:     new(1),
									ProviderName: new("AWS"),
									ReadOnlySpecs: &admin2025.DedicatedHardwareSpec20240805{
										DiskIOPS:      new(1000),
										DiskSizeGB:    new(10.0),
										EbsVolumeType: new("gp2"),
										InstanceSize:  new("M10"),
										NodeCount:     new(3),
									},
								},
							},
						},
					},
					RootCertType: new("X509"),
					Tags: &[]admin2025.ResourceTag{
						{Key: "key1", Value: "value1"},
						{Key: "key2", Value: "value2"},
					},
					TerminationProtectionEnabled: new(true),
					VersionReleaseSystem:         new("Atlas"),
					GroupId:                      new("32b6e34b3d91647abb20e7b8"),
				}
				testToAPI(t, "Cluster", input, nil, target, want)
			},
		},

		//nolint:dupl
		{
			name: "data federation all fields",
			test: func(t *testing.T) {
				input := &samplesv1.DataFederation{
					Spec: samplesv1.DataFederationSpec{
						V20250312: &samplesv1.DataFederationSpecV20250312{
							Entry: &samplesv1.DataFederationSpecV20250312Entry{
								CloudProviderConfig: &samplesv1.CloudProviderConfig{
									Aws: &samplesv1.Aws{
										RoleId:       "aws-role-id-123",
										TestS3Bucket: "my-s3-bucket",
									},
									Azure: &samplesv1.Azure{
										RoleId: "azure-role-id-456",
									},
									Gcp: &samplesv1.Azure{
										RoleId: "gcp-role-id-789",
									},
								},
								DataProcessRegion: &samplesv1.DataProcessRegion{
									CloudProvider: "GCE",
									Region:        "eu-north-2",
								},
								Name: new("some-name"),
								Storage: &samplesv1.Storage{
									Databases: &[]samplesv1.Databases{
										{
											Collections: &[]samplesv1.Collections{
												{
													DataSources: &[]samplesv1.DataSources{
														{
															AllowInsecure:       new(true),
															Collection:          new("some-name"),
															CollectionRegex:     new("collection-regex"),
															Database:            new("db"),
															DatabaseRegex:       new("db-regex"),
															DatasetName:         new("dataset-name"),
															DatasetPrefix:       new("dataset-prefix"),
															DefaultFormat:       new("default-format"),
															Path:                new("path"),
															ProvenanceFieldName: new("provenqance-field-name"),
															StoreName:           new("store-name"),
															TrimLevel:           new(1),
															Urls:                &[]string{"url1", "url2"},
														},
													},
													Name: new("collection0"),
												},
											},
											MaxWildcardCollections: new(3),
											Name:                   new("db0"),
											Views: &[]samplesv1.Views{
												{
													Name:     new("view0"),
													Pipeline: new("pipeline0"),
													Source:   new("source0"),
												},
											},
										},
									},
									Stores: &[]samplesv1.Stores{
										{
											AdditionalStorageClasses: &[]string{"stc1", "stc2"},
											AllowInsecure:            new(true),
											Bucket:                   new("bucket-name"),
											ClusterName:              new("cluster-name"),
											ContainerName:            new("container-name"),
											DefaultFormat:            new("default-format"),
											Delimiter:                new("delimiter"),
											IncludeTags:              new(true),
											Name:                     new("store-name"),
											Prefix:                   new("prefix"),
											Provider:                 "AWS",
											Public:                   new(true),
											ReadConcern: &samplesv1.ReadConcern{
												Level: new("local"),
											},
											ReadPreference: &samplesv1.ReadPreference{
												Mode: new("primary"),
											},
											Region:               new("us-east-1"),
											ReplacementDelimiter: new("replacement-delimiter"),
											ServiceURL:           new("https://service-url.com"),
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
					Name: new("some-name"),
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
												AllowInsecure:       new(true),
												Collection:          new("some-name"),
												CollectionRegex:     new("collection-regex"),
												Database:            new("db"),
												DatabaseRegex:       new("db-regex"),
												DatasetName:         new("dataset-name"),
												DatasetPrefix:       new("dataset-prefix"),
												DefaultFormat:       new("default-format"),
												Path:                new("path"),
												ProvenanceFieldName: new("provenqance-field-name"),
												StoreName:           new("store-name"),
												TrimLevel:           new(1),
												Urls:                &[]string{"url1", "url2"},
											},
										},
										Name: new("collection0"),
									},
								},
								MaxWildcardCollections: new(3),
								Name:                   new("db0"),
								Views: &[]admin2025.DataLakeApiBase{
									{
										Name:     new("view0"),
										Pipeline: new("pipeline0"),
										Source:   new("source0"),
									},
								},
							},
						},
						Stores: &[]admin2025.DataLakeStoreSettings{
							{
								AdditionalStorageClasses: &[]string{"stc1", "stc2"},
								AllowInsecure:            new(true),
								Bucket:                   new("bucket-name"),
								ClusterName:              new("cluster-name"),
								ContainerName:            new("container-name"),
								DefaultFormat:            new("default-format"),
								Delimiter:                new("delimiter"),
								IncludeTags:              new(true),
								Name:                     new("store-name"),
								Prefix:                   new("prefix"),
								Provider:                 "AWS",
								Public:                   new(true),
								ReadConcern: &admin2025.DataLakeAtlasStoreReadConcern{
									Level: new("local"),
								},
								ReadPreference: &admin2025.DataLakeAtlasStoreReadPreference{
									Mode: new("primary"),
								},
								Region:               new("us-east-1"),
								ReplacementDelimiter: new("replacement-delimiter"),
								ServiceURL:           new("https://service-url.com"),
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
				input := &samplesv1.DatabaseUser{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
					},
					Spec: samplesv1.DatabaseUserSpec{
						V20250312: &samplesv1.DatabaseUserSpecV20250312{
							Entry: &samplesv1.DatabaseUserSpecV20250312Entry{
								Username:     "test-user",
								DatabaseName: "admin",
								Roles: &[]samplesv1.Roles{
									{DatabaseName: "admin", RoleName: "readWrite"},
								},
								AwsIAMType:      new("aws-iam-type"),
								DeleteAfterDate: new("2025-07-01T00:00:00Z"),
								Description:     new("description"),
								Labels: &[]samplesv1.Tags{
									{Key: "key-1", Value: "value-1"},
									{Key: "key-2", Value: "value-2"},
								},
								LdapAuthType: new("ldap-auth-type"),
								OidcAuthType: new("oidc-auth-type"),
								PasswordSecretRef: &samplesv1.PasswordSecretRef{
									Name: "password-secret",
									Key:  new("password"),
								},
								Scopes: &[]samplesv1.Scopes{
									{Name: "scope-1", Type: "type-1"},
									{Name: "scope-2", Type: "type-2"},
								},
								X509Type: new("x509-type"),
							},
							GroupId: new("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.CloudDatabaseUser{}
				want := &admin2025.CloudDatabaseUser{
					Username:     "test-user",
					DatabaseName: "admin",
					GroupId:      "32b6e34b3d91647abb20e7b8",
					Roles: []admin2025.DatabaseUserRole{
						{DatabaseName: "admin", RoleName: "readWrite"},
					},
					AwsIAMType:      new("aws-iam-type"),
					DeleteAfterDate: new(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)),
					Description:     new("description"),
					Labels: &[]admin2025.ComponentLabel{
						{Key: new("key-1"), Value: new("value-1")},
						{Key: new("key-2"), Value: new("value-2")},
					},
					LdapAuthType: new("ldap-auth-type"),
					OidcAuthType: new("oidc-auth-type"),
					Password:     new("sample-password"),
					Scopes: &[]admin2025.UserScope{
						{Name: "scope-1", Type: "type-1"},
						{Name: "scope-2", Type: "type-2"},
					},
					X509Type: new("x509-type"),
				}
				objs := []client.Object{
					&corev1.Secret{
						TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
						ObjectMeta: metav1.ObjectMeta{Name: "password-secret", Namespace: "ns"},
						Data: map[string][]byte{
							"password": ([]byte)("sample-password"),
						},
					},
				}
				testToAPI(t, "DatabaseUser", input, objs, target, want)
			},
		},

		{
			name: "flex cluster with all fields",
			test: func(t *testing.T) {
				input := &samplesv1.FlexCluster{
					Spec: samplesv1.FlexClusterSpec{
						V20250312: &samplesv1.FlexClusterSpecV20250312{
							Entry: &samplesv1.FlexClusterSpecV20250312Entry{
								Name: "flex-cluster-name",
								ProviderSettings: samplesv1.ProviderSettings{
									BackingProviderName: "AWS",
									RegionName:          "us-east-1",
								},
								Tags: &[]samplesv1.Tags{
									{Key: "key1", Value: "value1"},
									{Key: "key2", Value: "value2"},
								},
								TerminationProtectionEnabled: new(true),
							},
							GroupId: new("32b6e34b3d91647abb20e7b8"),
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
					TerminationProtectionEnabled: new(true),
				}
				testToAPI(t, "FlexCluster", input, nil, target, want)
			},
		},

		{
			name: "simple group",
			test: func(t *testing.T) {
				input := &samplesv1.Group{
					Spec: samplesv1.GroupSpec{
						V20250312: &samplesv1.GroupSpecV20250312{
							Entry: &samplesv1.GroupSpecV20250312Entry{
								Name:                      "project-name",
								OrgId:                     "60987654321654321",
								RegionUsageRestrictions:   new("fake-restriction"),
								WithDefaultAlertsSettings: new(true),
								Tags: &[]samplesv1.Tags{
									{Key: "key", Value: "value"},
								},
							},
							// read only field, not translated back to the API
							ProjectOwnerId: new("61234567890123456"),
						},
					},
				}
				target := &admin2025.Group{}
				want := &admin2025.Group{
					Name:                      "project-name",
					OrgId:                     "60987654321654321",
					RegionUsageRestrictions:   new("fake-restriction"),
					WithDefaultAlertsSettings: new(true),
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
				input := &samplesv1.GroupAlertsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
					},
					Spec: samplesv1.GroupAlertsConfigSpec{
						V20250312: &samplesv1.GroupAlertsConfigSpecV20250312{
							Entry: &samplesv1.GroupAlertsConfigSpecV20250312Entry{
								Enabled:       new(true),
								EventTypeName: new("event-type"),
								Matchers: &[]samplesv1.Matchers{
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
								MetricThreshold: &samplesv1.MetricThreshold{
									MetricName: "metric-1",
									Mode:       new("mode"),
									Operator:   new("operator"),
									Threshold:  new(1.1),
									Units:      new("units"),
								},
								Threshold: &samplesv1.MetricThreshold{
									MetricName: "metric-t",
									Mode:       new("mode-t"),
									Operator:   new("operator-t"),
									Threshold:  new(2.2),
									Units:      new("units-t"),
								},
								Notifications: &[]samplesv1.Notifications{
									{
										DatadogApiKeySecretRef: &samplesv1.PasswordSecretRef{
											Name: "datadog-secret",
										},
										DatadogRegion: new("US"),
									},
								},
								SeverityOverride: new("some-severity-override"),
							},
							GroupId: new("60965432187654321"),
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
					Enabled:       new(true),
					EventTypeName: new("event-type"),
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
						Mode:       new("mode"),
						Operator:   new("operator"),
						Threshold:  new(1.1),
						Units:      new("units"),
					},
					Threshold: &admin2025.StreamProcessorMetricThreshold{
						MetricName: new("metric-t"),
						Mode:       new("mode-t"),
						Operator:   new("operator-t"),
						Threshold:  new(2.2),
						Units:      new("units-t"),
					},
					Notifications: &[]admin2025.AlertsNotificationRootForGroup{
						{
							DatadogApiKey: new("sample-password"),
							DatadogRegion: new("US"),
						},
					},
					GroupId:          new("60965432187654321"),
					SeverityOverride: new("some-severity-override"),
				}
				testToAPI(t, "GroupAlertsConfig", input, objs, target, want)
			},
		},

		{
			name: "sample network peering connection",
			test: func(t *testing.T) {
				input := &samplesv1.NetworkPeeringConnection{
					Spec: samplesv1.NetworkPeeringConnectionSpec{
						V20250312: &samplesv1.NetworkPeeringConnectionSpecV20250312{
							Entry: &samplesv1.NetworkPeeringConnectionSpecV20250312Entry{
								AccepterRegionName:  new("accepter-region-name"),
								AwsAccountId:        new("aws-account-id"),
								AzureDirectoryId:    new("azure-dir-id"),
								AzureSubscriptionId: new("azure-subcription-id"),
								ContainerId:         "container-id",
								GcpProjectId:        new("azure-subcription-id"),
								NetworkName:         new("net-name"),
								ProviderName:        new("provider-name"),
								ResourceGroupName:   new("resource-group-name"),
								RouteTableCidrBlock: new("cidr"),
								VnetName:            new("vnet-name"),
								VpcId:               new("vpc-id"),
							},
							GroupId: new("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.BaseNetworkPeeringConnectionSettings{}
				want := &admin2025.BaseNetworkPeeringConnectionSettings{
					ContainerId:         "container-id",
					ProviderName:        new("provider-name"),
					AccepterRegionName:  new("accepter-region-name"),
					AwsAccountId:        new("aws-account-id"),
					RouteTableCidrBlock: new("cidr"),
					VpcId:               new("vpc-id"),
					AzureDirectoryId:    new("azure-dir-id"),
					AzureSubscriptionId: new("azure-subcription-id"),
					ResourceGroupName:   new("resource-group-name"),
					VnetName:            new("vnet-name"),
					GcpProjectId:        new("azure-subcription-id"),
					NetworkName:         new("net-name"),
				}
				testToAPI(t, "NetworkPeeringConnection", input, nil, target, want)
			},
		},

		{
			name: "network permission entries all fields",
			test: func(t *testing.T) {
				input := &samplesv1.NetworkPermissionEntries{
					Spec: samplesv1.NetworkPermissionEntriesSpec{
						V20250312: &samplesv1.NetworkPermissionEntriesSpecV20250312{
							Entry: &[]samplesv1.NetworkPermissionEntriesSpecV20250312Entry{
								{
									AwsSecurityGroup: new("sg-12345678"),
									CidrBlock:        new("cird"),
									Comment:          new("comment"),
									DeleteAfterDate:  new("2025-07-01T00:00:00Z"),
									IpAddress:        new("1.1.1.1"),
								},
							},
							GroupId: new("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &NetworkPermissions{}
				want := &NetworkPermissions{
					Entry: []admin2025.NetworkPermissionEntry{
						{
							AwsSecurityGroup: new("sg-12345678"),
							CidrBlock:        new("cird"),
							Comment:          new("comment"),
							DeleteAfterDate:  new(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)),
							IpAddress:        new("1.1.1.1"),
						},
					},
				}
				testToAPI(t, "NetworkPermissionEntries", input, nil, target, want)
			},
		},

		{
			name: "sample organization",
			test: func(t *testing.T) {
				input := &samplesv1.Organization{
					Spec: samplesv1.OrganizationSpec{
						V20250312: &samplesv1.V20250312{
							Entry: &samplesv1.Entry{
								ApiKey: &samplesv1.ApiKey{
									Desc:  "description",
									Roles: []string{"role-1", "role-2"},
								},
								FederationSettingsId:      new("fed-id"),
								Name:                      "org-name",
								OrgOwnerId:                new("org-owner-id"),
								SkipDefaultAlertsSettings: new(true),
							},
						},
					},
				}
				target := &admin2025.AtlasOrganization{}
				want := &admin2025.AtlasOrganization{
					Name:                      "org-name",
					SkipDefaultAlertsSettings: new(true),
				}
				testToAPI(t, "Organization", input, nil, target, want)
			},
		},

		{
			name: "Organization setting with all fields",
			test: func(t *testing.T) {
				input := &samplesv1.OrganizationSetting{
					Spec: samplesv1.OrganizationSettingSpec{
						V20250312: &samplesv1.OrganizationSettingSpecV20250312{
							Entry: &samplesv1.V20250312Entry{
								ApiAccessListRequired:                  new(true),
								GenAIFeaturesEnabled:                   new(true),
								MaxServiceAccountSecretValidityInHours: new(24),
								MultiFactorAuthRequired:                new(true),
								RestrictEmployeeAccess:                 new(true),
								SecurityContact:                        new("contact-info"),
								StreamsCrossGroupEnabled:               new(true),
							},
							OrgId: "org-id",
						},
					},
				}
				target := &admin2025.OrganizationSettings{}
				want := &admin2025.OrganizationSettings{
					ApiAccessListRequired:                  new(true),
					GenAIFeaturesEnabled:                   new(true),
					MaxServiceAccountSecretValidityInHours: new(24),
					MultiFactorAuthRequired:                new(true),
					RestrictEmployeeAccess:                 new(true),
					SecurityContact:                        new("contact-info"),
					StreamsCrossGroupEnabled:               new(true),
				}
				testToAPI(t, "OrganizationSetting", input, nil, target, want)
			},
		},

		{
			name: "customrole with all fields",
			test: func(t *testing.T) {
				input :=
					&samplesv1.CustomRole{
						Spec: samplesv1.CustomRoleSpec{
							V20250312: &samplesv1.CustomRoleSpecV20250312{
								Entry: &samplesv1.CustomRoleSpecV20250312Entry{
									RoleName: "custom-role-name",
									Actions: &[]samplesv1.Actions{
										{
											Action: "action1",
											Resources: &[]samplesv1.Resources{
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
									InheritedRoles: &[]samplesv1.InheritedRoles{
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
							Resources: []admin2025.DatabasePermittedNamespaceResource{
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
				input := &samplesv1.SearchIndex{
					Spec: samplesv1.SearchIndexSpec{
						V20250312: &samplesv1.SearchIndexSpecV20250312{
							Entry: &samplesv1.SearchIndexSpecV20250312Entry{
								Database:       "database-name",
								CollectionName: "collection-name",
								Name:           "index-name",
								Type:           new("search-index-type"),
								Definition: &samplesv1.Definition{
									Analyzer: new("lucene.standard"),
									Analyzers: &[]samplesv1.Analyzers{
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
									Mappings: &samplesv1.Mappings{
										Dynamic: &apiextensionsv1.JSON{Raw: []byte(`true`)},
										Fields: &map[string]apiextensionsv1.JSON{
											"field1": {Raw: []byte(`{"key4":"value4"}`)},
										},
									},
									NumPartitions:  new(3),
									SearchAnalyzer: new("lucene.standard"),
									StoredSource: &apiextensionsv1.JSON{
										Raw: []byte(`{"enabled": true}`),
									},
									Synonyms: &[]samplesv1.Synonyms{
										{
											Analyzer: "synonym-analyzer",
											Name:     "synonym-name",
											Source: samplesv1.Source{
												Collection: "synonym-collection",
											},
										},
									},
								},
							},
							GroupId: new("group-id-101"),
						},
					},
				}
				target := &admin2025.SearchIndexCreateRequest{}
				want := &admin2025.SearchIndexCreateRequest{
					CollectionName: "collection-name",
					Database:       "database-name",
					Name:           "index-name",
					Type:           new("search-index-type"),
					Definition: &admin2025.BaseSearchIndexCreateRequestDefinition{
						Analyzer: new("lucene.standard"),
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
							Dynamic: true,
							Fields: &map[string]any{
								"field1": map[string]any{"key4": "value4"},
							},
						},
						NumPartitions:  new(3),
						SearchAnalyzer: new("lucene.standard"),
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
				input := &samplesv1.Team{
					Spec: samplesv1.TeamSpec{
						V20250312: &samplesv1.TeamSpecV20250312{
							Entry: &samplesv1.TeamSpecV20250312Entry{
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
				input := &samplesv1.ThirdPartyIntegration{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
					},
					Spec: samplesv1.ThirdPartyIntegrationSpec{
						V20250312: &samplesv1.ThirdPartyIntegrationSpecV20250312{
							IntegrationType: "ANY",
							Entry: &samplesv1.ThirdPartyIntegrationSpecV20250312Entry{
								AccountId: new("account-id"),
								ApiKeySecretRef: &samplesv1.PasswordSecretRef{
									Key:  new("apiKey"),
									Name: "multi-secret0",
								},
								ApiTokenSecretRef: &samplesv1.PasswordSecretRef{
									Key:  new("apiToken"),
									Name: "multi-secret0",
								},
								ChannelName: new("channel-name"),
								Enabled:     new(true),
								LicenseKeySecretRef: &samplesv1.PasswordSecretRef{
									Key:  new("licenseKey"),
									Name: "multi-secret1",
								},
								Region:                       new("some-region"),
								SendCollectionLatencyMetrics: new(true),
								SendDatabaseMetrics:          new(true),
								SendUserProvidedResourceTags: new(true),
								ServiceDiscovery:             new("service-discovery"),
								TeamName:                     new("some-team"),
								Type:                         new("some-type"),
								Username:                     new("username"),
							},
							GroupId: new("32b6e34b3d91647abb20e7b8"),
						},
					},
				}
				target := &admin2025.ThirdPartyIntegration{}
				want := &admin2025.ThirdPartyIntegration{
					Type:                         new("some-type"),
					ApiKey:                       new("sample-api-key"),
					Region:                       new("some-region"),
					SendCollectionLatencyMetrics: new(true),
					SendDatabaseMetrics:          new(true),
					SendUserProvidedResourceTags: new(true),
					AccountId:                    new("account-id"),
					LicenseKey:                   new("sample-license-key"),
					Enabled:                      new(true),
					ServiceDiscovery:             new("service-discovery"),
					Username:                     new("username"),
					ApiToken:                     new("sample-api-token"),
					ChannelName:                  new("channel-name"),
					TeamName:                     new("some-team"),
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
	trs, err := crapi.NewPerVersionTranslators(testScheme(t), crd, version, sdkVersion)
	require.NoError(t, err)
	tr := trs[sdkVersion]
	require.NotNil(t, tr)
	require.NoError(t, tr.ToAPI(target, input, objs...))
	assert.Equal(t, want, target)
}

// TestFromAPIRefMapping tests that when translating from the API
// with a groupId, the translator properly maps the reference based on provided dependencies.
func TestFromAPIRefMapping(t *testing.T) {
	const (
		groupID   = "62b6e34b3d91647abb20e7b8"
		groupName = "my-project"
	)

	// The Group object that represents the project referenced by groupId
	group := &samplesv1.Group{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Group",
			APIVersion: "atlas.generated.mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      groupName,
			Namespace: "default",
		},
		Spec: samplesv1.GroupSpec{
			V20250312: &samplesv1.GroupSpecV20250312{
				Entry: &samplesv1.GroupSpecV20250312Entry{
					Name:  "My Project",
					OrgId: "org123456789",
				},
			},
		},
		Status: samplesv1.GroupStatus{
			V20250312: &samplesv1.GroupStatusV20250312{
				Id: pointer.MakePtr(groupID),
			},
		},
	}

	for _, tc := range []struct {
		name                 string
		apiCluster           admin2025.ClusterDescription20240805
		targetCluster        *samplesv1.Cluster
		referencedObjects    []client.Object
		wantGroupRefName     string
		wantGroupId          *string
		wantExtraObjectCount int
	}{
		{
			name: "with Group reference - should convert groupId to groupRef",
			apiCluster: admin2025.ClusterDescription20240805{
				Name:        new("my-cluster"),
				ClusterType: new("REPLICASET"),
				GroupId:     pointer.MakePtr(groupID),
				ReplicationSpecs: &[]admin2025.ReplicationSpec20240805{
					{
						ZoneName: new("Zone 1"),
						RegionConfigs: &[]admin2025.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ElectableSpecs: &admin2025.HardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(3),
								},
							},
						},
					},
				},
			},
			targetCluster: &samplesv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-cluster",
					Namespace: "default",
				},
				Spec: samplesv1.ClusterSpec{
					V20250312: &samplesv1.ClusterSpecV20250312{
						GroupRef: &crd2gok8s.LocalReference{},
						Entry:    &samplesv1.ClusterSpecV20250312Entry{},
					},
				},
			},
			referencedObjects:    []client.Object{group},
			wantGroupRefName:     groupName,
			wantGroupId:          nil,
			wantExtraObjectCount: 0,
		},
		{
			name: "without Group reference - should keep groupId",
			apiCluster: admin2025.ClusterDescription20240805{
				Name:        new("my-cluster"),
				ClusterType: new("REPLICASET"),
				GroupId:     pointer.MakePtr(groupID),
			},
			targetCluster: &samplesv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-cluster",
					Namespace: "default",
				},
			},
			referencedObjects:    nil,
			wantGroupRefName:     "",
			wantGroupId:          pointer.MakePtr(groupID),
			wantExtraObjectCount: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scheme := testScheme(t)
			crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
			crd, err := extractCRD("Cluster", bufio.NewScanner(crdsYML))
			require.NoError(t, err)

			translator, err := crapi.NewTranslator(scheme, crd, version, sdkVersion)
			require.NoError(t, err)

			extraObjects, err := translator.FromAPI(tc.targetCluster, &tc.apiCluster, tc.referencedObjects...)
			require.NoError(t, err)

			assert.NotNil(t, tc.targetCluster.Spec.V20250312, "spec.v20250312 should not be nil")

			if tc.wantGroupRefName != "" {
				assert.NotNil(t, tc.targetCluster.Spec.V20250312.GroupRef,
					"spec.v20250312.groupRef should not be nil")
				assert.Equal(t, tc.wantGroupRefName, tc.targetCluster.Spec.V20250312.GroupRef.Name,
					"groupRef.name should match expected value")
			}

			if tc.wantGroupId == nil {
				assert.Nil(t, tc.targetCluster.Spec.V20250312.GroupId,
					"groupId should be nil when groupRef is set")
			} else {
				assert.Equal(t, tc.wantGroupId, tc.targetCluster.Spec.V20250312.GroupId,
					"groupId should match expected value")
			}

			assert.Len(t, extraObjects, tc.wantExtraObjectCount,
				"extra objects count should match expected value")
		})
	}
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

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, k8sscheme.AddToScheme(scheme))
	require.NoError(t, samplesv1.AddToScheme(scheme))
	return scheme
}
