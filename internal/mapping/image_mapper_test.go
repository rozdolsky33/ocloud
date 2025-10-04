package mapping_test

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/require"
)

func TestImage_Mappers_OCI_To_Attrs_To_Domain(t *testing.T) {
	id := "ocid1.image.oc1..abc"
	name := "Oracle-Linux-9"
	os := "Oracle Linux"
	osv := "9"
	launchMode := core.ImageLaunchModeNative
	created := time.Now().UTC().Truncate(time.Second)

	ociImg := core.Image{
		Id:                     &id,
		DisplayName:            &name,
		OperatingSystem:        &os,
		OperatingSystemVersion: &osv,
		LaunchMode:             launchMode,
		TimeCreated:            &common.SDKTime{Time: created},
	}

	attrs := mapping.NewImageAttributesFromOCIImage(ociImg)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.ID)
	require.Equal(t, &name, attrs.DisplayName)
	require.Equal(t, &os, attrs.OperatingSystem)
	require.Equal(t, &osv, attrs.OperatingSystemVersion)
	require.Equal(t, string(launchMode), attrs.LaunchMode)
	require.NotNil(t, attrs.TimeCreated)
	require.True(t, created.Equal(*attrs.TimeCreated))

	// Build domain image from attributes
	img := mapping.NewDomainImageFromAttrs(*attrs)
	require.IsType(t, domain.Image{}, img)
	require.Equal(t, id, img.OCID)
	require.Equal(t, name, img.DisplayName)
	require.Equal(t, os, img.OperatingSystem)
	require.Equal(t, osv, img.OperatingSystemVersion)
	require.Equal(t, string(launchMode), img.LaunchMode)
	require.True(t, created.Equal(img.TimeCreated))
}

func TestNewDomainImageFromAttrs_ZeroValues(t *testing.T) {
	// Using empty attributes should produce zero-value domain Image
	attrs := mapping.ImageAttributes{}
	img := mapping.NewDomainImageFromAttrs(attrs)
	require.Equal(t, "", img.OCID)
	require.Equal(t, "", img.DisplayName)
	require.Equal(t, "", img.OperatingSystem)
	require.Equal(t, "", img.OperatingSystemVersion)
	require.Equal(t, "", img.LaunchMode)
	require.True(t, img.TimeCreated.IsZero())
}
