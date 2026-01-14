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

package tag

import (
	"go.mongodb.org/atlas-sdk/v20250312012/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

type Tag struct {
	*akov2.TagSpec
}

func FromAtlas(rTags []admin.ResourceTag) []*akov2.TagSpec {
	tags := make([]*akov2.TagSpec, 0, len(rTags))
	for _, rTag := range rTags {
		tags = append(
			tags,
			&akov2.TagSpec{
				Key:   rTag.GetKey(),
				Value: rTag.GetValue(),
			},
		)
	}

	return tags
}

func ToAtlas(tags []*akov2.TagSpec) *[]admin.ResourceTag {
	if tags == nil {
		return nil
	}

	rTags := make([]admin.ResourceTag, 0, len(tags))
	for _, tag := range tags {
		rTags = append(
			rTags,
			*admin.NewResourceTag(tag.Key, tag.Value),
		)
	}

	return &rTags
}

func FlexToAtlas(tags []*akov2.TagSpec) *[]admin.ResourceTag {
	if tags == nil {
		return nil
	}

	rTags := make([]admin.ResourceTag, 0, len(tags))
	for _, tag := range tags {
		rTags = append(
			rTags,
			*admin.NewResourceTag(tag.Key, tag.Value),
		)
	}

	return &rTags
}
