package models

func MapToTags(tags map[string]string) *[]Tag {
	ret := make([]Tag, 0, len(tags))
	for key, val := range tags {
		ret = append(ret, Tag{
			Key:   key,
			Value: val,
		})
	}
	return &ret
}
