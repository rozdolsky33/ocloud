package objectstorage_test

import (
	"context"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/config"
	ocloud_objectstorage "github.com/rozdolsky33/ocloud/internal/oci/storage/objectstorage"
)

// TestIntegrationListBuckets performs an integration test of the ListBuckets function.
// NOTE: This test requires OCI credentials to be configured.
func TestIntegrationListBuckets(t *testing.T) {
	t.Skip("skipping integration test that requires live OCI; run manually when configured")
	provider := config.LoadOCIConfig()

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		t.Fatalf("Failed to create object storage client: %v", err)
	}

	adapter := ocloud_objectstorage.NewAdapter(client)

	// This is the OCID of the root compartment of my tenancy.
	// You will need to replace this with a valid compartment OCID from your tenancy.
	compartmentID := "ocid1.tenancy.oc1..aaaaaaaa3k2ljgq4z4z4z4z4z4z4z4z4z4z4z4z4z4z4"

	buckets, err := adapter.ListBuckets(context.Background(), compartmentID)
	if err != nil {
		t.Fatalf("ListBuckets failed: %v", err)
	}

	if len(buckets) == 0 {
		t.Log("Warning: No buckets found in the compartment. The test passed, but it didn't verify any data.")
	}
}
