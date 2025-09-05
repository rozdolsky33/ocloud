package autonomousdb

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// Adapter implements the domain.AutonomousDatabaseRepository interface for OCI.
type Adapter struct {
	dbClient database.DatabaseClient
}

// NewAdapter creates a new Adapter instance. The compartmentID parameter is accepted for
// backward compatibility with service wiring but is not required by the adapter itself.
func NewAdapter(provider oci.ClientProvider) (*Adapter, error) {
	dbClient, err := oci.NewDatabaseClient(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}
	return &Adapter{
		dbClient: dbClient,
	}, nil
}

// GetAutonomousDatabase retrieves a single Autonomous Database and maps it to the domain model.
func (a *Adapter) GetAutonomousDatabase(ctx context.Context, ocid string) (*domain.AutonomousDatabase, error) {
	response, err := a.dbClient.GetAutonomousDatabase(ctx, database.GetAutonomousDatabaseRequest{
		AutonomousDatabaseId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get autonomous database: %w", err)
	}
	db := a.toDomainAutonomousDB(response.AutonomousDatabase)
	return &db, nil
}

// ListAutonomousDatabases retrieves a list of autonomous databases from OCI.
func (a *Adapter) ListAutonomousDatabases(ctx context.Context, compartmentID string) ([]domain.AutonomousDatabase, error) {
	var allDatabases []domain.AutonomousDatabase
	var page *string
	for {
		resp, err := a.dbClient.ListAutonomousDatabases(ctx, database.ListAutonomousDatabasesRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list autonomous databases: %w", err)
		}
		for _, item := range resp.Items {
			allDatabases = append(allDatabases, a.toDomainAutonomousDB(item))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return allDatabases, nil
}

// toDomainAutonomousDB maps either a full database.AutonomousDatabase (from Get) or a database.AutonomousDatabaseSummary (from List) into the single domain.AutonomousDatabase type.
func (a *Adapter) toDomainAutonomousDB(ociObj interface{}) domain.AutonomousDatabase {
	var (
		name              *string
		id                *string
		privateEndpoint   *string
		privateEndpointIp *string
		connStrings       *database.AutonomousDatabaseConnectionStrings
		freeformTags      map[string]string
		definedTags       map[string]map[string]interface{}
	)

	switch src := ociObj.(type) {
	case database.AutonomousDatabase:
		name = src.DbName
		id = src.Id
		privateEndpoint = src.PrivateEndpoint
		privateEndpointIp = src.PrivateEndpointIp
		connStrings = src.ConnectionStrings
		freeformTags = src.FreeformTags
		definedTags = src.DefinedTags
	case database.AutonomousDatabaseSummary:
		name = src.DbName
		id = src.Id
		privateEndpoint = src.PrivateEndpoint
		privateEndpointIp = src.PrivateEndpointIp
		connStrings = src.ConnectionStrings
		freeformTags = src.FreeformTags
		definedTags = src.DefinedTags
	default:
		return domain.AutonomousDatabase{}
	}

	d := domain.AutonomousDatabase{}
	if name != nil {
		d.Name = *name
	}
	if id != nil {
		d.ID = *id
	}
	if privateEndpoint != nil {
		d.PrivateEndpoint = *privateEndpoint
	}
	if privateEndpointIp != nil {
		d.PrivateEndpointIp = *privateEndpointIp
	}
	if connStrings != nil {
		if connStrings.AllConnectionStrings != nil {
			d.ConnectionStrings = connStrings.AllConnectionStrings
		}
		d.Profiles = connStrings.Profiles
	}
	d.FreeformTags = freeformTags
	d.DefinedTags = definedTags
	return d
}
