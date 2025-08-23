package domain

import (
	"context"
	"time"
)

// Image represents a bootable image for a compute instance.
// This is our application's internal representation, decoupled from the OCI SDK.
type Image struct {
	OCID                   string
	DisplayName            string
	OperatingSystem        string
	OperatingSystemVersion string
	LaunchMode             string
	TimeCreated            time.Time
}

// ImageRepository defines the port for interacting with image storage.
type ImageRepository interface {
	ListImages(ctx context.Context, compartmentID string) ([]Image, error)
	GetImage(ctx context.Context, ocid string) (*Image, error)
}
