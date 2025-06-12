// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
