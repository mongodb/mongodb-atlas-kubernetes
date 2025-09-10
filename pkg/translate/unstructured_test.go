package translate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin2025 "go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/josvazg/akotranslate/internal/pointer"
)

type testStruct struct {
	ID        int
	YesNo     bool
	Text      string
	Data      float32
	Timestamp time.Time
	SubStruct []subStruct
}

type subStruct struct {
	SubID   int
	YesNo   bool
	Text    string
	Data    float32
	TheTime time.Time
}

func TestParamsFill(t *testing.T) {
	unstructuredSample := map[string]any{
		"groupid": "62b6e34b3d91647abb20e7b8",
		"groupalertsconfig": map[string]any{
			"enabled":       true,
			"eventtypename": "some-event",
		},
	}
	result := admin2025.CreateAlertConfigurationApiParams{}
	require.NoError(t, fromUnstructured(&result, unstructuredSample))
	assert.Equal(t, admin2025.CreateAlertConfigurationApiParams{
		GroupId: "62b6e34b3d91647abb20e7b8",
		GroupAlertsConfig: &admin2025.GroupAlertsConfig{
			Enabled:       pointer.Get(true),
			EventTypeName: pointer.Get("some-event"),
		},
	}, result)
}

func TestToAndFromUnstructured(t *testing.T) {
	sample := testStruct{
		ID:        12,
		YesNo:     true,
		Text:      "some text",
		Data:      2.4,
		Timestamp: time.Now().UTC().Round(time.Nanosecond),
		SubStruct: []subStruct{
			{
				SubID:   15,
				YesNo:   true,
				Text:    "more text",
				Data:    0.67,
				TheTime: time.Now().Add(-24 * 7 * time.Hour).UTC().Round(time.Nanosecond),
			},
			{
				SubID:   14,
				YesNo:   false,
				Text:    "and even more text",
				Data:    1.67,
				TheTime: time.Now().Add(-24 * 15 * time.Hour).UTC().Round(time.Nanosecond),
			},
		},
	}
	sampleMap, err := toUnstructured(sample)
	require.NoError(t, err)
	result := testStruct{}
	require.NoError(t, fromUnstructured(&result, sampleMap))
	assert.Equal(t, sample, result)
}
