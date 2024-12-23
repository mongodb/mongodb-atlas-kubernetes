package tag

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

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
