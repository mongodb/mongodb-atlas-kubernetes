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
	assert.NoError(t, err)
	assert.Equal(t, httpClient.Timeout, decorated.Timeout)
	assert.Equal(t, httpClient.Jar, decorated.Jar)
	assert.NotNil(t, decorated.Transport)

	// not going deeper here, just need to confirm that transport was changed
	assert.NotEqual(t, httpClient.Transport, decorated.Transport)
}
