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

package searchindex

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"
)

func TestSearchIndexes_DeleteIndex(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		execErr    error
		wantErr    bool
	}{
		{
			name:       "succeeds on 204",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "succeeds on 202",
			statusCode: http.StatusAccepted,
			wantErr:    false,
		},
		{
			name:       "succeeds on 404",
			statusCode: http.StatusNotFound,
			wantErr:    false,
		},
		{
			name:       "fails on non accepted status",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "fails when execute errors",
			statusCode: http.StatusAccepted,
			execErr:    errors.New("atlas API error"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSearchAPI := mockadmin.NewAtlasSearchAPI(t)
			req := admin.DeleteClusterSearchIndexApiRequest{ApiService: mockSearchAPI}

			mockSearchAPI.EXPECT().
				DeleteClusterSearchIndex(context.Background(), "project-id", "cluster-name", "index-id").
				Return(req)
			mockSearchAPI.EXPECT().
				DeleteClusterSearchIndexExecute(req).
				Return(&http.Response{StatusCode: tt.statusCode}, tt.execErr)

			err := NewSearchIndexes(mockSearchAPI).DeleteIndex(context.Background(), "project-id", "cluster-name", "index-id")
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
