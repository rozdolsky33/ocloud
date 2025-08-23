package image

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/domain"
)

// Adapter is an infrastructure-layer adapter that implements the domain.ImageRepository interface.
type Adapter struct {
	client core.ComputeClient
}

// NewAdapter creates a new adapter for interacting with OCI images.
func NewAdapter(client core.ComputeClient) *Adapter {
	return &Adapter{client: client}
}

// GetImage retrieves a single image by its OCID.
func (a *Adapter) GetImage(ctx context.Context, ocid string) (*domain.Image, error) {
	resp, err := a.client.GetImage(ctx, core.GetImageRequest{
		ImageId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("getting image from OCI: %w", err)
	}

	img := a.toDomainModel(resp.Image)
	return &img, nil
}

// ListImages retrieves all images in a given compartment.
func (a *Adapter) ListImages(ctx context.Context, compartmentID string) ([]domain.Image, error) {
	var images []domain.Image
	page := ""

	for {
		resp, err := a.client.ListImages(ctx, core.ListImagesRequest{
			CompartmentId: &compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing images from OCI: %w", err)
		}

		for _, item := range resp.Items {
			images = append(images, a.toDomainModel(item))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	return images, nil
}

// toDomainModel converts an OCI SDK image object to our application's domain model.
func (a *Adapter) toDomainModel(img core.Image) domain.Image {
	return domain.Image{
		OCID:                   *img.Id,
		DisplayName:            *img.DisplayName,
		OperatingSystem:        *img.OperatingSystem,
		OperatingSystemVersion: *img.OperatingSystemVersion,
		LaunchMode:             string(img.LaunchMode),
		TimeCreated:            img.TimeCreated.Time,
	}
}
