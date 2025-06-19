package image

import (
	"fmt"
	"strings"
)

// flattenTags converts ResourceTags into a single space-separated string of normalized key-value pairs.
// FreeformTags are formatted as "key:value", while DefinedTags are formatted as "namespace.key:value".
func flattenTags(tags ResourceTags) string {
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

// ToIndexableImage converts an Image object into an IndexableImage structure optimized for indexing and searching.
func ToIndexableImage(img Image) IndexableImage {
	return IndexableImage{
		ID:              img.ID,
		Name:            strings.ToLower(img.Name),
		OperatingSystem: strings.ToLower(img.OperatingSystem),
		ImageOSVersion:  strings.ToLower(img.ImageOSVersion),
		Tags:            flattenTags(img.ImageTags),
	}
}
