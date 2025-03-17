package config

import (
	"net/url"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadOpenAPI(filePath string) (*openapi3.T, error) {
	loader := &openapi3.Loader{
		IsExternalRefsAllowed: true,
	}

	uri, err := url.Parse(filePath)
	if err == nil && uri.Scheme != "" && uri.Host != "" {
		return loader.LoadFromURI(uri)
	}

	return loader.LoadFromFile(filepath.Clean(filePath))
}
