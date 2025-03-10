package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	quayUrl = "https://quay.io/v1/api/repository"
)

type Tag struct {
	Name         string `json:"name"`
	LastModified string `json:"last_modified"`
}

type TagResponse struct {
	Tags []Tag `json:"tags"`
}

func main() {
	removeOlderThanDays := flag.Uint64("remove-older-than-days", 30, "Remove images older than this many days")
	token := flag.String("token", "", "Quay API token")
	flag.Parse()

	if *token == "" {
		fmt.Println("Quay API token is required")
		os.Exit(-1)
	}
	fmt.Println("Starting cleaner")
	fmt.Printf("Removing images older than %d days\r\n", *removeOlderThanDays)

	client := &http.Client{Timeout: time.Second * 10}

	var expiredTags []string
	tags, err := fetchTags(client, *token)
	if err != nil {
		fmt.Println("Failed to fetch tags:", err)
		os.Exit(-1)
	}

	currentTime := time.Now()
	var errs []error
	for i := range tags {
		tag := &tags[i]
		tagTime, err := time.Parse(time.RFC3339, tag.LastModified)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse time for tag %s: %w", tag.Name, err))
			continue
		}

		delta := currentTime.Sub(tagTime).Hours() / 24
		if uint64(delta) > *removeOlderThanDays {
			expiredTags = append(expiredTags, tag.Name)
			fmt.Println("Tag", tag.Name, "is older than", *removeOlderThanDays, "days")
		}
	}

	for _, tag := range expiredTags {
		err := deleteTag(client, *token, tag)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to delete tag %s: %w", tag, err))
		} else {
			fmt.Println("Deleted tag", tag)
		}
	}

	if len(errs) > 0 {
		fmt.Println("Errors:")
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(-1)
	}

	fmt.Println("Done")
}

func fetchTags(client *http.Client, token string) ([]Tag, error) {
	req, err := http.NewRequest("GET", quayUrl+"/mongodb/mongodb-atlas-kubernetes/tag/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch tags, status: %d", resp.StatusCode)
	}

	var tagResp TagResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagResp); err != nil {
		return nil, err
	}

	return tagResp.Tags, nil
}

func deleteTag(client *http.Client, token, tag string) error {
	url := fmt.Sprintf("%s/mongodb/mongodb-atlas-kubernetes/tag/%s", quayUrl, tag)

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete tag %s, status: %d, response: %s", tag, resp.StatusCode, string(body))
	}

	return nil
}
