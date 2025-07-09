package translate_test

import (
	"bufio"
	"embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin2023 "go.mongodb.org/atlas-sdk/v20231115014/admin"
	admin2024 "go.mongodb.org/atlas-sdk/v20241113005/admin"
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
			sdkVersion: "v20231115",
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
			target: &admin2023.Group{},
			want: &admin2023.Group{
				Name:                      "project-name",
				OrgId:                     "60987654321654321",
				RegionUsageRestrictions:   pointer.Get("fake-restriction"),
				WithDefaultAlertsSettings: pointer.Get(true),
				Tags: &[]admin2023.ResourceTag{
					{Key: "key", Value: "value"},
				},
			},
		},
		{
			name:       "group alert config with project and credential references",
			crd:        "GroupAlertsConfig",
			sdkVersion: "v20241113",
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
							Operator:  pointer.Get("op"),
							Threshold: pointer.Get(1),
							Units:     pointer.Get("units"),
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
					GroupId: "60965432187654321",
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
						"key": ([]byte)("sample-password"), // should be apiKey, not key
					},
				},
			},
			target: &admin2024.GroupAlertsConfig{},
			want: &admin2024.GroupAlertsConfig{
				Enabled:       pointer.Get(true),
				EventTypeName: pointer.Get("event-type"),
				Matchers: &[]admin2024.StreamsMatcher{
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
				MetricThreshold: &admin2024.FlexClusterMetricThreshold{
					MetricName: "metric-1",
					Mode:       pointer.Get("mode"),
					Operator:   pointer.Get("operator"),
					Threshold:  pointer.Get(1.1),
					Units:      pointer.Get("units"),
				},
				Threshold: &admin2024.GreaterThanRawThreshold{
					Operator:  pointer.Get("op"),
					Threshold: pointer.Get(1),
					Units:     pointer.Get("units"),
				},
				Notifications: &[]admin2024.AlertsNotificationRootForGroup{
					{
						DatadogApiKey: pointer.Get("sample-password"),
						DatadogRegion: pointer.Get("US"),
					},
				},
				GroupId: pointer.Get("60965432187654321"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			crdsYML, err := samples.Open("samples/crds.yml")
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
