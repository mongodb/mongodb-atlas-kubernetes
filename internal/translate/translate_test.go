package translate_test

import (
	"bufio"
	"embed"
	"fmt"
	"testing"

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
			sdkVersion: "V20231115",
			spec: v1.GroupSpec{
				V20231115: &v1.GroupSpecV20231115{
					Entry: &v1.GroupSpecV20231115Entry{
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
			sdkVersion: "V20241113",
			spec: v1.GroupAlertsConfigSpec{
				V20241113: &v1.GroupAlertsConfigSpecV20241113{
					Entry: &v1.GroupAlertsConfigSpecV20241113Entry{
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
						Threshold: &v1.Threshold{
							Operator:  pointer.Get("operator0"),
							Units:     pointer.Get("units0"),
							Threshold: pointer.Get(2),
						},
						Notifications: &[]v1.Notifications{
							{
								DatadogApiKeySecretRef: &v1.ApiTokenSecretRef{
									Name: pointer.Get("datadog-secret"),
								},
								DatadogRegion: pointer.Get("US"),
							},
						},
					},
					//GroupId: "60965432187654321",
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
					Operator:   pointer.Get("op"),
					Units:      pointer.Get("units"),
					MetricName: pointer.Get("metric"),
					Mode:       pointer.Get("mode"),
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
