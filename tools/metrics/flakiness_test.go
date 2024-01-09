package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlakyRank(t *testing.T) {
	testCase := struct {
		inputs   []jobID
		expected flakyRank
	}{
		inputs: []jobID{
			{Name: "B", RunID: 2, JobID: 1},
			{Name: "C", RunID: 3, JobID: 1},
			{Name: "B", RunID: 2, JobID: 2},
			{Name: "D", RunID: 4, JobID: 1},
			{Name: "B", RunID: 2, JobID: 3},
			{Name: "E", RunID: 5, JobID: 1},
			{Name: "B", RunID: 2, JobID: 4},
			{Name: "E", RunID: 5, JobID: 2},
			{Name: "A", RunID: 1, JobID: 1},
			{Name: "C", RunID: 3, JobID: 2},
			{Name: "D", RunID: 4, JobID: 2},
			{Name: "C", RunID: 3, JobID: 3},
			{Name: "D", RunID: 4, JobID: 3},
			{Name: "A", RunID: 1, JobID: 2},
			{Name: "A", RunID: 1, JobID: 3},
			{Name: "A", RunID: 1, JobID: 4},
			{Name: "A", RunID: 1, JobID: 5},
		},
		expected: flakyRank{
			rank: []*testFlakiness{
				{
					testIdentifier: testIdentifier{Name: "A", testType: E2E},
					tests: []jobID{
						{Name: "A", RunID: 1, JobID: 1},
						{Name: "A", RunID: 1, JobID: 2},
						{Name: "A", RunID: 1, JobID: 3},
						{Name: "A", RunID: 1, JobID: 4},
						{Name: "A", RunID: 1, JobID: 5},
					},
				},
				{
					testIdentifier: testIdentifier{Name: "B", testType: E2E},
					tests: []jobID{
						{Name: "B", RunID: 2, JobID: 1},
						{Name: "B", RunID: 2, JobID: 2},
						{Name: "B", RunID: 2, JobID: 3},
						{Name: "B", RunID: 2, JobID: 4},
					},
				},
				{
					testIdentifier: testIdentifier{Name: "C", testType: E2E},
					tests: []jobID{
						{Name: "C", RunID: 3, JobID: 1},
						{Name: "C", RunID: 3, JobID: 2},
						{Name: "C", RunID: 3, JobID: 3},
					},
				},
				{
					testIdentifier: testIdentifier{Name: "D", testType: E2E},
					tests: []jobID{
						{Name: "D", RunID: 4, JobID: 1},
						{Name: "D", RunID: 4, JobID: 2},
						{Name: "D", RunID: 4, JobID: 3},
					},
				},
				{
					testIdentifier: testIdentifier{Name: "E", testType: E2E},
					tests: []jobID{
						{Name: "E", RunID: 5, JobID: 1},
						{Name: "E", RunID: 5, JobID: 2},
					},
				},
			},
		},
	}
	rank := flakyRank{}
	for _, jid := range testCase.inputs {
		rank.add(identify(jid.Name), jid)
	}
	require.Equal(t, testCase.expected, rank)
}

func TestQueryFlakiness(t *testing.T) {
	srs, err := QueryFlakiness(newTestClient(), lastRecordingTime, Weekly, 2)
	assert.NoError(t, err)
	require.NotNil(t, srs)
}
