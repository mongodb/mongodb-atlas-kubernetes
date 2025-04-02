package main

import (
	"testing"
)

func TestMatchWildcards(t *testing.T) {
	tests := []struct {
		name           string
		prLabels       []string
		testLabels     []string
		testType       string
		expectedResult []string
	}{
		{
			name:           "Exact Match Test",
			prLabels:       []string{"test/int/users"},
			testLabels:     []string{"users", "books", "orders"},
			testType:       "int",
			expectedResult: []string{"users"},
		},
		{
			name:           "Wildcard match Test",
			prLabels:       []string{"test/int/user*"},
			testLabels:     []string{"users", "user-profiles", "books", "orders"},
			testType:       "int",
			expectedResult: []string{"users", "user-profiles"},
		},
		{
			name:           "Wildcard match all tests",
			prLabels:       []string{"test/e2e/*"},
			testLabels:     []string{"login", "signup", "users", "books", "orders"},
			testType:       "e2e",
			expectedResult: []string{"login", "signup", "users", "books", "orders"},
		},
		{
			name:           "Wildcards with prefix",
			prLabels:       []string{"test/e2e/login*"},
			testLabels:     []string{"login", "login-new", "signup", "users", "books", "orders"},
			testType:       "e2e",
			expectedResult: []string{"login", "login-new"},
		},
		{
			name:           "No tests should match",
			prLabels:       []string{"test/int/unknown"},
			testLabels:     []string{"users", "books", "orders"},
			testType:       "int",
			expectedResult: []string{},
		},
		{
			name:           "Wildcard on integration and e2e tests but for e2e",
			prLabels:       []string{"test/int/user*", "test/e2e/login*"},
			testLabels:     []string{"users", "user-profiles", "login", "signup"},
			testType:       "e2e",
			expectedResult: []string{"login"},
		},
		{
			name:           "Wildcard on integration and e2e tests but for integration",
			prLabels:       []string{"test/int/user*", "test/e2e/login*"},
			testLabels:     []string{"users", "user-profiles", "login", "signup"},
			testType:       "int",
			expectedResult: []string{"users", "user-profiles"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchWildcards(tt.prLabels, tt.testLabels, tt.testType)

			gotMap := make(map[string]struct{})
			for _, label := range result {
				gotMap[label] = struct{}{}
			}

			expectedMap := make(map[string]struct{})
			for _, label := range tt.expectedResult {
				expectedMap[label] = struct{}{}
			}

			if len(gotMap) != len(expectedMap) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expectedResult, result)
			}

			for label := range expectedMap {
				if _, found := gotMap[label]; !found {
					t.Errorf("Test %s failed: expected %v to be in the result.\r\n\tExpected: %s\r\n\tGot: %s",
						tt.name, label, jsonDump(expectedMap), jsonDump(gotMap))
				}
			}
		})
	}
}
