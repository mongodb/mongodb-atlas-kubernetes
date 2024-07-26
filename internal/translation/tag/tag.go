package tag

import "go.mongodb.org/atlas-sdk/v20231115008/admin"

type Tag struct {
	Key   string
	Value string
}

func FromAtlas(rTags []admin.ResourceTag) []Tag {
	tags := make([]Tag, 0, len(rTags))
	for _, rTag := range rTags {
		tags = append(
			tags,
			Tag{
				Key:   rTag.GetKey(),
				Value: rTag.GetValue(),
			},
		)
	}

	return tags
}

func ToAtlas(tags []Tag) *[]admin.ResourceTag {
	rTags := make([]admin.ResourceTag, 0, len(tags))
	for _, tag := range tags {
		rTags = append(
			rTags,
			*admin.NewResourceTag(tag.Key, tag.Value),
		)
	}

	return &rTags
}
