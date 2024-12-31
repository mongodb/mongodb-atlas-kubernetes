package validate

import (
	"errors"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func Tags(tags []*akov2.TagSpec) error {
	tagsMap := make(map[string]struct{}, len(tags))

	for _, currTag := range tags {
		if _, ok := tagsMap[currTag.Key]; ok {
			return errors.New("duplicate keys found in tags, this is forbidden")
		}

		tagsMap[currTag.Key] = struct{}{}
	}

	return nil
}
