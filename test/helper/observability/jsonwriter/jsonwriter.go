package jsonwriter

import (
	"encoding/json"
	"io"
	"time"
)

type JSONWriter struct {
	delegate io.Writer
	level    string
	logger   string
}

func NewJSONWriter(delegate io.Writer, level, logger string) *JSONWriter {
	return &JSONWriter{
		delegate: delegate,
		level:    level,
		logger:   logger,
	}
}

func (j *JSONWriter) Write(b []byte) (int, error) {
	var js json.RawMessage
	if json.Unmarshal(b, &js) == nil {
		return j.delegate.Write(b)
	}

	entry := struct {
		Msg    string `json:"msg,omitempty"`
		Time   string `json:"time,omitempty"`
		Level  string `json:"level,omitempty"`
		Logger string `json:"logger,omitempty"`
	}{
		Time:   time.Now().UTC().Format(time.RFC3339),
		Level:  j.level,
		Msg:    string(b),
		Logger: j.logger,
	}

	enc := json.NewEncoder(j.delegate)
	if err := enc.Encode(&entry); err != nil {
		return 0, err
	}

	return len(b), nil
}
