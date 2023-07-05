package fakerest_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/fakerest"
)

func TestSetBodyDefaults(t *testing.T) {
	rsp := fakerest.NewResponseTo(get(t, "/"))
	fakerest.SetBody(rsp, "[]")
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
}

func testServer1() *fakerest.Server {
	return fakerest.NewServer(fakerest.Script{
		{
			Method:    http.MethodGet,
			URIPrefix: "/some/path",
			Reply:     fakerest.NotFoundJSONReply,
		},
		{
			Method:    http.MethodDelete,
			URIPrefix: "/some/path/item1",
			Reply:     fakerest.RemovedJSONReply,
		},
	})
}

func testServer2() *fakerest.Server {
	return fakerest.NewServer(fakerest.Script{
		{
			Method:    http.MethodGet,
			URIPrefix: "/",
			Reply:     fakerest.NotFoundJSONReply,
		},
		{
			Method:    http.MethodPost,
			URIPrefix: "/item0",
			Reply:     fakerest.NotFoundJSONReply,
		},
		{
			Method:    http.MethodDelete,
			URIPrefix: "/item0",
			Reply:     fakerest.RemovedJSONReply,
		},
	})
}

func get(t *testing.T, path string) *http.Request {
	return do(t, http.MethodGet, path)
}

func delete(t *testing.T, path string) *http.Request {
	return do(t, http.MethodDelete, path)
}

func do(t *testing.T, method, path string) *http.Request {
	t.Helper()

	uri, err := url.Parse(path)
	require.NoError(t, err)
	return &http.Request{Method: method, URL: uri}
}

func bodyAsString(t *testing.T, rsp *http.Response) string {
	t.Helper()

	if rsp.Body == nil {
		return ""
	}
	buf := bytes.NewBufferString("")
	_, err := io.Copy(buf, rsp.Body)
	require.NoError(t, err)
	return buf.String()
}

func TestServer(t *testing.T) {
	svr := testServer1()
	rsp, err := svr.RoundTrip(get(t, "/some/path"))
	require.NoError(t, err)
	assert.Equal(t, "{}", bodyAsString(t, rsp))
}

func TestServerPanics(t *testing.T) {
	svr := testServer1()
	assert.Panics(t, func() {
		svr.RoundTrip(get(t, "/"))
	})
}

func TestCombinedServer(t *testing.T) {
	svr := fakerest.NewCombinedServer(testServer1(), testServer2())

	rsp, err := svr.RoundTrip(get(t, "/"))
	require.NoError(t, err)
	assert.Equal(t, "{}", bodyAsString(t, rsp))

	rsp, err = svr.RoundTrip(delete(t, "/some/path/item1"))
	require.NoError(t, err)
	assert.Equal(t, "{}", bodyAsString(t, rsp))
}

type serverWithState struct {
	*fakerest.Server
	used bool
}

func NewServerWithState() *serverWithState {
	sws := &serverWithState{used: false}
	sws.Server = fakerest.NewServer(fakerest.Script{
		{
			Method:    http.MethodGet,
			URIPrefix: "/",
			Reply: func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				sws.used = true
				return fakerest.NotFoundJSONReply(req, rsp)
			},
		},
	})
	return sws
}

type previousServerWithState struct {
	*fakerest.Server
	alsoUsed bool
}

func NewPreviousServerWithState() *previousServerWithState {
	psws := &previousServerWithState{alsoUsed: false}
	psws.Server = fakerest.NewServer(fakerest.Script{
		{
			Method:    http.MethodGet,
			URIPrefix: "/prev",
			Reply: func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				psws.alsoUsed = true
				return fakerest.NotFoundJSONReply(req, rsp)
			},
		},
	})
	return psws
}

func TestCompositeServerWithStatus(t *testing.T) {
	psws := NewPreviousServerWithState()
	sws := NewServerWithState()
	svr := fakerest.NewCombinedServer(
		psws.Server,
		sws.Server,
	)
	rsp, err := svr.RoundTrip(get(t, "/prev"))
	require.NoError(t, err)
	assert.Equal(t, "{}", bodyAsString(t, rsp))
	assert.True(t, psws.alsoUsed)
}
