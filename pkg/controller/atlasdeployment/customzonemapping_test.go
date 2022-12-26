package atlasdeployment

import (
	"testing"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

type CMZTestData struct {
	desired             []v1.CustomZoneMapping
	existing            map[string]string
	czmMap              map[string]string
	expectedToCreate    bool
	expectedToBeDeleted bool
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
		t.Errorf("expected to shouldCreate %v, got %v", testData.expectedToCreate, shouldCreate)
	}
	if shouldDelete != testData.expectedToBeDeleted {
		t.Errorf("expected to shouldDelete %v, got %v", testData.expectedToBeDeleted, shouldDelete)
	}
}

func TestCompareZoneMappingStates_ShouldDoNothing(t *testing.T) {
	data := &CMZTestData{
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
	}
	runCMZTest(t, data)
}

func TestCompareZoneMappingStates_WrongZone(t *testing.T) {
	data := &CMZTestData{
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
	}
	runCMZTest(t, data)
}

func TestCompareZoneMappingStates_Recreate(t *testing.T) {
	data := &CMZTestData{
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
	}
	runCMZTest(t, data)
}

func TestCompareZoneMappingStates_DeleteOnly(t *testing.T) {
	data := &CMZTestData{
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
	}
	runCMZTest(t, data)
}

func TestCompareZoneMappingStates_AddOnly(t *testing.T) {
	data := &CMZTestData{
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
	}
	runCMZTest(t, data)
}
