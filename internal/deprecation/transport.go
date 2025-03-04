package deprecation

import (
	"net/http"

	"go.uber.org/zap"
)

type Transport struct {
	delegate http.RoundTripper
	logger   *zap.Logger
}

func NewLoggingTransport(delegate http.RoundTripper, logger *zap.Logger) *Transport {
	return &Transport{
		delegate: delegate,
		logger:   logger.Named("deprecated"),
	}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.delegate.RoundTrip(req)
	if resp != nil {
		javaMethod := resp.Header.Get("X-Java-Method")
		deprecation := resp.Header.Get("Deprecation")
		sunset := resp.Header.Get("Sunset")

		if deprecation != "" {
			t.logger.Warn("deprecation", zap.String("date", deprecation), zap.String("javaMethod", javaMethod), zap.String("path", req.URL.Path), zap.String("method", req.Method))
		}

		if sunset != "" {
			t.logger.Warn("sunset", zap.String("date", sunset), zap.String("javaMethod", javaMethod), zap.String("path", req.URL.Path), zap.String("method", req.Method))
		}
	}
	return resp, err
}
