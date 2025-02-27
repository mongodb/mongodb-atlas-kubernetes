package httputil

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type DiffClient struct {
	origin http.RoundTripper
}

func NewDiffClient(origin http.RoundTripper) *DiffClient {
	fmt.Println("DIFF: Creating DiffClient")
	return &DiffClient{
		origin: origin,
	}
}

func (c *DiffClient) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println("DIFF: RoundTrip called with", req.Method, req.URL.String())
	switch req.Method {
	case http.MethodPatch, http.MethodPut:
		if !strings.Contains(req.URL.String(), "cloud-qa.mongodb.com") {
			fmt.Println("DIFF: Skipping diff for non-qa request:", req.URL.String())
			break
		}
		d, err := c.computeDiff(req)
		if err != nil {
			fmt.Println("DIFF: Error computing diff:", err)
			break
		}
		fmt.Println(`DIFF: Computed diff:`, d)
	}
	return c.origin.RoundTrip(req)
}

func (c *DiffClient) computeDiff(req *http.Request) (string, error) {
	fmt.Println("DIFF: Computing diff for request:", req.Method, req.URL.String())
	// Clone the original request
	originBody, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("DIFF: Error reading request body:", err)
		return "", err
	}

	// Make a get request to the same URL
	fmt.Println("DIFF: Making GET request to", req.URL.String())
	getReq, _ := http.NewRequest(http.MethodGet, req.URL.String(), nil)
	// Don't forget to copy headers
	for k, v := range req.Header {
		for _, val := range v {
			getReq.Header.Add(k, val)
		}
	}
	getResp, err := c.origin.RoundTrip(getReq)
	if err != nil {
		fmt.Println("DIFF: Error making GET request:", err)
		return "", err
	}
	getBody, err := io.ReadAll(getResp.Body)
	if err != nil {
		fmt.Println("DIFF: Error reading GET response body:", err)
		return "", err
	}
	fmt.Println("DIFF: GET response body:", string(getBody))

	diffStr := ""
	dmp := diffmatchpatch.New()
	for _, diff := range dmp.DiffMain(string(originBody), string(getBody), false) {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			diffStr += fmt.Sprintf(`++%s`, diff.Text) // Green text in git
		case diffmatchpatch.DiffDelete:
			diffStr += fmt.Sprintf(`--%s`, diff.Text) // Red text in git
		case diffmatchpatch.DiffEqual:
			diffStr += fmt.Sprintf(` %s`, diff.Text) // No change
		}
	}
	return fmt.Sprintf(
		`DIFF: Original: %s\r\n GOT: %s\r\n, DIFF: %s\r\n`,
		string(originBody),
		string(getBody),
		diffStr), nil
}
