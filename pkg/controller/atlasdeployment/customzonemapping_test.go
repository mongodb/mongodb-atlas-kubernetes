package atlasdeployment

import (
	"testing"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type CMZTestData struct {
	desired             []v1.CustomZoneMapping
	existing            map[string]string
	czmMap              map[string]string
	expectedToCreate    bool
	expectedToBeDeleted bool
	name                string
}

const (
	zone1     = "Zone 1"
	zone2     = "Zone 2"
	location1 = "CA"
	location2 = "US"
)

func runCMZTest(t *testing.T, testData *CMZTestData) {
	shouldCreate, shouldDelete := compareZoneMappingStates(testData.existing, testData.desired, testData.czmMap)
	if shouldCreate != testData.expectedToCreate {
		t.Errorf("Test: %s. expected to shouldCreate %v, got %v", testData.name, testData.expectedToCreate, shouldCreate)
	}
	if shouldDelete != testData.expectedToBeDeleted {
		t.Errorf("Test: %s. expected to shouldDelete %v, got %v", testData.name, testData.expectedToBeDeleted, shouldDelete)
	}
}

func TestCompareZoneMappingStates(t *testing.T) {
	tests := []*CMZTestData{
		{
			name: "All synced. No changes needed",
			desired: []v1.CustomZoneMapping{
				{
					Zone:     zone1,
					Location: location1,
				},
				{
					Zone:     zone2,
					Location: location2,
				},
			},
			existing: map[string]string{
				location1: "1",
				location2: "2",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    false,
			expectedToBeDeleted: false,
		},
		{
			name: "Wrong zone. Should be recreated",
			desired: []v1.CustomZoneMapping{
				{
					Zone:     zone1,
					Location: location1,
				},
				{
					Zone:     zone2,
					Location: location2,
				},
			},
			existing: map[string]string{
				location1: "1",
				location2: "1",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    true,
			expectedToBeDeleted: true,
		},
		{
			name: "Exist more than needed. Should be recreated",
			desired: []v1.CustomZoneMapping{
				{
					Zone:     zone1,
					Location: location1,
				},
			},
			existing: map[string]string{
				location1: "1",
				location2: "2",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    true,
			expectedToBeDeleted: true,
		},
		{
			name:    "Empty desired. Should be deleted",
			desired: []v1.CustomZoneMapping{},
			existing: map[string]string{
				location1: "1",
				location2: "2",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    false,
			expectedToBeDeleted: true,
		},
		{
			name: "Exist less than needed. Should be created",
			desired: []v1.CustomZoneMapping{
				{
					Zone:     zone1,
					Location: location1,
				},
				{
					Zone:     zone2,
					Location: location2,
				},
			},
			existing: map[string]string{
				location2: "2",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    true,
			expectedToBeDeleted: false,
		},
	}
	for _, test := range tests {
		runCMZTest(t, test)
	}
}
