package instance

import (
	"fmt"
	"strings"
)

type IndexableInstance struct {
	ID              string
	Name            string
	OperatingSystem string
	CreatedAt       string
	Tags            string
}

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

func ToIndexableInstance(instance Instance) IndexableInstance {
	return IndexableInstance{
		ID:              instance.ID,
		Name:            strings.ToLower(instance.Name),
		OperatingSystem: strings.ToLower(instance.ImageOS),
		CreatedAt:       strings.ToLower(instance.CreatedAt.String()),
		Tags:            flattenTags(instance.InstanceTags),
	}
}
