package autonomousdb

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/oci"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Adapter implements the domain.AutonomousDatabaseRepository interface for OCI.
type Adapter struct {
	dbClient database.DatabaseClient
}

// NewAdapter creates a new Adapter instance.
func NewAdapter(provider oci.ClientProvider) (*Adapter, error) {
	dbClient, err := oci.NewDatabaseClient(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}
	return &Adapter{
		dbClient: dbClient,
	}, nil
}

func (a *Adapter) GetAutonomousDatabase(ctx context.Context, ocid string) (*domain.AutonomousDatabase, error) {
	//TODO implement me
	panic("implement me")
}

// ListAutonomousDatabases retrieves a list of autonomous databases from OCI.
func (a *Adapter) ListAutonomousDatabases(ctx context.Context, compartmentID string) ([]domain.AutonomousDatabase, error) {
	var allDatabases []domain.AutonomousDatabase
	page := ""
	for {
		resp, err := a.dbClient.ListAutonomousDatabases(ctx, database.ListAutonomousDatabasesRequest{
			CompartmentId: &compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list autonomous databases: %w", err)
		}
		for _, item := range resp.Items {
			allDatabases = append(allDatabases, mapAutonomousDatabaseSummaryToDomain(item))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}
	return allDatabases, nil
}

// FindAutonomousDatabase finds a specific autonomous database by name.
// This implementation fetches all and filters, which might not be optimal for very large numbers of databases.
// OCI SDK does not provide a direct "FetchPaginatedImages by Name" for autonomous databases.
func (a *Adapter) FindAutonomousDatabase(ctx context.Context, compartmentID, name string) (*domain.AutonomousDatabase, error) {
	dbs, err := a.ListAutonomousDatabases(ctx, compartmentID)
	if err != nil {
		return nil, err
	}
	for _, db := range dbs {
		if db.Name == name {
			return &db, nil
		}
	}
	return nil, domain.NewNotFoundError("autonomous database", name)
}

// mapAutonomousDatabaseSummaryToDomain transforms a database.AutonomousDatabaseSummary instance into a domain.AutonomousDatabase struct.
func mapAutonomousDatabaseSummaryToDomain(db database.AutonomousDatabaseSummary) domain.AutonomousDatabase {
	return domain.AutonomousDatabase{
		Name:              *db.DbName,
		ID:                *db.Id,
		PrivateEndpoint:   *db.PrivateEndpoint,
		PrivateEndpointIp: *db.PrivateEndpointIp,
		ConnectionStrings: db.ConnectionStrings.AllConnectionStrings,
		Profiles:          db.ConnectionStrings.Profiles,
		DatabaseTags:      util.ConvertOciTagsToResourceTags(db.FreeformTags, db.DefinedTags),
	}
}
