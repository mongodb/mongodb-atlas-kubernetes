package httputil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_DecorateClient(t *testing.T) {
	httpClient := &http.Client{Transport: http.DefaultTransport}
	withDigest := Digest("publiApi", "privateApi")
	withLogging := LoggingTransport(zap.S())

	decorated, err := DecorateClient(&http.Client{Transport: http.DefaultTransport}, withDigest, withLogging)
	a := assert.New(t)
	a.NoError(err)
	a.Equal(httpClient.Timeout, decorated.Timeout)
	a.Equal(httpClient.Jar, decorated.Jar)
	a.NotNil(decorated.Transport)

	// not going deeper here, just need to confirm that transport was changed
	a.NotEqual(t, httpClient.Transport, decorated.Transport)
}
