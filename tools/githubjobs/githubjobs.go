package main

import (
	"context"
	"github.com/google/go-github/v57/github"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log := l.Sugar()

	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	w, ok := os.LookupEnv("GH_RUN_ID")
	if !ok {
		log.Fatal("GH_RUN_ID environment variable not set")
	}

	workflowId, err := strconv.ParseInt(w, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// auth with github
	client := github.NewClient(nil).WithAuthToken(token)

	// get first page of jobs for id
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	jobList := make([]int64, 0)
	page := 1

	for {
		opt := github.ListWorkflowJobsOptions{ListOptions: github.ListOptions{Page: page}}
		jobs, _, err := client.Actions.ListWorkflowJobs(ctx, "mongodb", "mongodb-atlas-kubernetes", workflowId, &opt)
		if err != nil {
			log.Fatal(err)
		}

		for _, job := range jobs.Jobs {
			jobList = append(jobList, *job.ID)
		}
		if len(jobList) == jobs.GetTotalCount() {
			break
		}

		page++
	}

	for _, job := range jobList {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		url, _, err := client.Actions.GetWorkflowJobLogs(ctx, "mongodb", "mongodb-atlas-kubernetes", job, 1)
		if err != nil {
			log.Errorf("failed to get job logs for %v: %v", job, err)
			continue
		}

		r, err := http.Get(url.String())
		if err != nil {
			log.Fatalf("failed to download job logs for %v: %v", job, err)
		}

		_, err = io.Copy(os.Stdout, r.Body)
		if err != nil {
			log.Fatalf("failed to print job logs to stdout for %v: %v", job, err)
		}
	}
}
