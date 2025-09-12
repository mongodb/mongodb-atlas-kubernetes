package main

import (
	"context"
	"fmt"
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
		fmt.Println("GITHUB_TOKEN environment variable not set")
		os.Exit(1)
	}

	w, ok := os.LookupEnv("GH_RUN_ID")
	if !ok {
		fmt.Println("GH_RUN_ID environment variable not set")
		os.Exit(1)
	}

	workflowId, err := strconv.ParseInt(w, 10, 64)
	if err != nil {
		panic(err)
	}

	// auth with github
	client := github.NewClient(nil).WithAuthToken(token)

	// get first page of jobs for id
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	jobList := make([]int64, 0)
	page := 1

	for {
		fmt.Println("getting page" + strconv.Itoa(page))
		opt := github.ListWorkflowJobsOptions{ListOptions: github.ListOptions{Page: page}}
		jobs, _, err := client.Actions.ListWorkflowJobs(ctx, "mongodb", "mongodb-atlas-kubernetes", workflowId, &opt)
		if err != nil {
			panic(err)
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
