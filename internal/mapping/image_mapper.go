package mapping

import (
	compute "github.com/rozdolsky33/ocloud/internal/domain/compute"
)

// NewDomainImageFromAttrs builds a domain Image from provider-agnostic attributes.
func NewDomainImageFromAttrs(attrs ImageAttributes) compute.Image {
	img := compute.Image{
		OCID:                   stringValue(attrs.ID),
		DisplayName:            stringValue(attrs.DisplayName),
		OperatingSystem:        stringValue(attrs.OperatingSystem),
		OperatingSystemVersion: stringValue(attrs.OperatingSystemVersion),
		LaunchMode:             attrs.LaunchMode,
	}
	if attrs.TimeCreated != nil {
		img.TimeCreated = *attrs.TimeCreated
	}
	return img
}
