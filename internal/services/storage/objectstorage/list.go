package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	osadapter "github.com/rozdolsky33/ocloud/internal/oci/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListBuckets retrieves and lists all buckets, allows browsing objects within them,
// and performs actions (view details or download) on selected objects.
func ListBuckets(appCtx *app.ApplicationContext, useJSON bool) error {
	ctx := context.Background()
	client, err := oci.NewObjectStorageClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating object storage client: %w", err)
	}

	bucketAdapter := osadapter.NewAdapter(client)
	service := NewService(bucketAdapter, appCtx.Logger, appCtx.CompartmentID)

	// Get namespace once for all operations
	namespace, err := service.GetNamespace(ctx)
	if err != nil {
		return fmt.Errorf("getting namespace: %w", err)
	}

	// Get region from provider for URL generation
	region, err := appCtx.Provider.Region()
	if err != nil {
		return fmt.Errorf("getting region: %w", err)
	}

	// Outer loop for bucket selection (allows going back from object list)
	for {
		buckets, err := service.ListBuckets(ctx)
		if err != nil {
			return fmt.Errorf("listing buckets: %w", err)
		}

		// Show bucket list TUI
		bucketModel := osadapter.NewBucketListModel(buckets)
		bucketID, err := tui.Run(bucketModel)
		if err != nil {
			if errors.Is(err, tui.ErrCancelled) {
				return nil // User pressed Esc, exit
			}
			return fmt.Errorf("selecting bucket: %w", err)
		}

		// Get bucket name from OCID
		bucketName, err := service.osRepo.GetBucketNameByOCID(ctx, appCtx.CompartmentID, bucketID)
		if err != nil {
			return fmt.Errorf("getting bucket name: %w", err)
		}

		// Inner loop for object browsing (Esc returns to bucket list)
		for {
			objects, err := service.ListObjects(ctx, namespace, bucketName)
			if err != nil {
				return fmt.Errorf("listing objects: %w", err)
			}

			if len(objects) == 0 {
				fmt.Fprintf(appCtx.Stdout, "Bucket %q is empty.\n", bucketName)
				break // Go back to bucket selection
			}

			// Show object list TUI
			objectModel := osadapter.NewObjectListModel(objects, bucketName)
			objectName, err := tui.Run(objectModel)
			if err != nil {
				if errors.Is(err, tui.ErrCancelled) {
					break // Go back to bucket selection
				}
				return fmt.Errorf("selecting object: %w", err)
			}

			// Show action picker (radio button style)
			actionModel := osadapter.NewActionPickerModel(objectName)
			action, err := tui.RunPicker(actionModel)
			if err != nil {
				if errors.Is(err, tui.ErrCancelled) {
					return nil // Exit
				}
				return fmt.Errorf("selecting action: %w", err)
			}

			// Execute selected action
			switch action {
			case "view":
				obj, err := service.GetObjectDetails(ctx, namespace, bucketName, objectName)
				if err != nil {
					return fmt.Errorf("getting object details: %w", err)
				}
				return PrintObjectInfo(obj, appCtx, region, useJSON)

			case "download":
				// Download to current working directory with progress TUI
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("getting current directory: %w", err)
				}

				title := fmt.Sprintf("Downloading %s", objectName)
				progressRunner := tui.NewProgressRunner(title)
				progressRunner.Start()

				// Channel to signal completion
				done := make(chan error, 1)

				// Start download in goroutine
				go func() {
					progressFn := func(p TransferProgress) {
						percent := float64(p.BytesTransferred) / float64(p.TotalBytes)
						bytesInfo := fmt.Sprintf("%s / %s",
							util.HumanizeBytesIEC(p.BytesTransferred),
							util.HumanizeBytesIEC(p.TotalBytes))

						status := "Downloading..."
						if percent >= 1.0 {
							status = "Complete!"
						}

						progressRunner.UpdateProgress(percent, bytesInfo, "", status)
					}

					err := service.DownloadObject(ctx, namespace, bucketName, objectName, cwd, progressFn)
					if err != nil {
						progressRunner.SendError(err)
						done <- err
						return
					}
					progressRunner.SendDone()
					done <- nil
				}()

				// Run the progress TUI (blocks until complete)
				if err := progressRunner.Run(); err != nil {
					return err
				}

				// Wait for download to finish and check result
				downloadErr := <-done
				if downloadErr != nil {
					return fmt.Errorf("downloading object: %w", downloadErr)
				}

				fmt.Fprintf(appCtx.Stdout, "\nDownloaded %q to %s\n", objectName, cwd)
				return nil
			}
		}
	}
}
