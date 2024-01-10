package main

import (
	"context"
	"os"

	"github.com/google/go-github/v57/github"
)

type WorkflowsQuerier interface {
	// TestWorkflowRuns are all the test run by AKO at page X (descendent order)
	TestWorkflowRuns(branch, event string, page int) (*github.WorkflowRuns, error)

	// TestWorkflowRunJobs are all the jobs at a given Workflow Run at page X (descendent order)
	TestWorkflowRunJobs(runID int64, filter string, page int) (*github.Jobs, error)
}

type QueryClient interface {
	WorkflowsQuerier
}

type ghQueryClient struct {
	author string
	repo   string
	client *github.Client
}

func NewDefaultQueryClient() *ghQueryClient {
	return newGHQueryClient(akoAuthor, ako, os.Getenv("GITHUB_TOKEN"))
}

func newGHQueryClient(author, repo, token string) *ghQueryClient {
	client := newGoGitHubClient(token)
	return &ghQueryClient{author: author, repo: repo, client: client}
}

// TestWorkflowRuns implements QueryCLient.
func (ghc *ghQueryClient) TestWorkflowRuns(branch, event string, page int) (*github.WorkflowRuns, error) {
	wfRuns, _, err := ghc.client.Actions.ListWorkflowRunsByFileName(
		context.Background(),
		ghc.author,
		ghc.repo,
		akoWorkflowFilename,
		&github.ListWorkflowRunsOptions{
			Branch: branch,
			Event:  event,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: PerPage,
			},
		},
	)
	return wfRuns, err
}

// TestWorkflowRunJobs implements QueryCLient.
func (ghc *ghQueryClient) TestWorkflowRunJobs(runID int64, filter string, page int) (*github.Jobs, error) {
	jobs, _, err := ghc.client.Actions.ListWorkflowJobs(
		context.Background(),
		ghc.author,
		ghc.repo,
		runID,
		&github.ListWorkflowJobsOptions{
			Filter: filter,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: PerPage,
			},
		})
	return jobs, err
}

func newGoGitHubClient(token string) *github.Client {
	client := github.NewClient(nil)
	if token != "" {
		client = client.WithAuthToken(token)
	}
	return client
}
