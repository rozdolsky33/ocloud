package util

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/domain"
)

// ShowConstructionAnimation displays a placeholder animation indicating that a feature is under construction.
func ShowConstructionAnimation() {
	fmt.Println("🚧 This feature is not implemented yet. Coming soon!")
}

// ConvertOciTagsToResourceTags converts OCI FreeformTags and DefinedTags to domain.ResourceTags.
func ConvertOciTagsToResourceTags(freeformTags map[string]string, definedTags map[string]map[string]interface{}) domain.ResourceTags {
	resourceTags := make(domain.ResourceTags)
	for k, v := range freeformTags {
		resourceTags[k] = v
	}
	for namespace, tags := range definedTags {
		for k, v := range tags {
			if strVal, ok := v.(string); ok {
				resourceTags[fmt.Sprintf("%s.%s", namespace, k)] = strVal
			}
		}
	}
	return resourceTags
}
