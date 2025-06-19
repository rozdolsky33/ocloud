package images

import (
	"fmt"
	"strings"
)

type IndexableImage struct {
	ID              string
	Name            string
	OperatingSystem string
	ImageOSVersion  string
	Tags            string
}

func flattenTags(tags ImageTags) string {
	var parts []string
	for k, v := range tags.FreeformTags {
		parts = append(parts, fmt.Sprintf("%s:%s", strings.ToLower(k), strings.ToLower(v)))
	}
	for ns, kv := range tags.DefinedTags {
		for k, v := range kv {
			parts = append(parts, fmt.Sprintf("%s.%s:%v", strings.ToLower(ns), strings.ToLower(k), v))
		}
	}
	return strings.Join(parts, " ")
}

func ToIndexableImage(img Image) IndexableImage {
	return IndexableImage{
		ID:              img.ID,
		Name:            strings.ToLower(img.Name),
		OperatingSystem: strings.ToLower(img.OperatingSystem),
		ImageOSVersion:  strings.ToLower(img.ImageOSVersion),
		Tags:            flattenTags(img.ImageTags),
	}
}
