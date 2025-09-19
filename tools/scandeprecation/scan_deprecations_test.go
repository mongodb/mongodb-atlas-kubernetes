package main

import "testing"

const testLogLine = `2025-07-21T14:41:06.118Z	WARN	controllers.AtlasDeployment.deprecated	sunset	***"type": "sunset", "date": "Thu, 12 Mar 2026 00:00:00 GMT", "javaMethod": "ApiAtlasSearchDeploymentResource::getSearchDeployment", "path": "/api/atlas/v2/groups/687e4d839dc4c25a9d96f39d/clusters/search-nodes-test/search/deployment", "method": "GET"***`

func TestParseLogLine(t *testing.T) {
	out, err := parseLogLine(testLogLine)
	if err != nil {
		t.Fatal(err)
	}
	if out.Type != "sunset" && out.Date != "Thu, 12 Mar 2026 00:00:00 GMT" && out.JavaMethod != "ApiAtlasSearchDeploymentResource::getSearchDeployment" {
		t.Errorf("parseLogLine did not output expected struct")
	}
}

func TestParseLogLineErrors(t *testing.T) {
	_, err := parseLogLine("abd123")
	if err == nil {
		t.Errorf("parseLogLine() did not return an error")
	}
}
