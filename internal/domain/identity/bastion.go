package identity

import "context"

// Bastion represents a bastion host in the cloud.
// This is our application's internal representation, decoupled from the OCI SDK.
type Bastion struct {
	OCID                     string
	DisplayName              string
	BastionType              string
	LifecycleState           string
	CompartmentID            string
	TargetVcnID              string
	TargetVcnName            string
	TargetSubnetID           string
	TargetSubnetName         string
	MaxSessionTTL            int
	ClientCidrBlockAllowList []string
	PrivateEndpointIpAddress string
	FreeformTags             map[string]string
	DefinedTags              map[string]map[string]interface{}
	TimeCreated              string
	TimeUpdated              string
}

// BastionSession represents a bastion session.
type BastionSession struct {
	OCID                    string
	DisplayName             string
	BastionID               string
	BastionName             string
	LifecycleState          string
	SessionType             string
	SessionTTL              int
	TargetResourceID        string
	TargetResourceFQDN      string
	TargetResourcePort      int
	TargetResourcePrivateIP string
	SSHMetadata             map[string]string
	TimeCreated             string
	TimeUpdated             string
}

// BastionRepository defines the port for interacting with bastion storage.
// The application layer will use this interface to talk to the infrastructure layer.
type BastionRepository interface {
	// ListBastions lists all bastions in the specified compartment
	ListBastions(ctx context.Context, compartmentID string) ([]Bastion, error)

	// GetBastion retrieves a specific bastion by ID
	GetBastion(ctx context.Context, bastionID string) (*Bastion, error)

	// CreateBastion creates a new bastion
	CreateBastion(ctx context.Context, request CreateBastionRequest) (*Bastion, error)

	// DeleteBastion deletes a bastion
	DeleteBastion(ctx context.Context, bastionID string) error
}

// BastionSessionRepository defines the port for interacting with bastion session storage.
type BastionSessionRepository interface {
	// ListSessions lists all sessions for a bastion
	ListSessions(ctx context.Context, bastionID string) ([]BastionSession, error)

	// GetSession retrieves a specific session by ID
	GetSession(ctx context.Context, sessionID string) (*BastionSession, error)

	// CreateSession creates a new bastion session
	CreateSession(ctx context.Context, request CreateSessionRequest) (*BastionSession, error)

	// DeleteSession deletes a bastion session
	DeleteSession(ctx context.Context, sessionID string) error
}

// CreateBastionRequest represents the request to create a bastion
type CreateBastionRequest struct {
	DisplayName    string
	CompartmentID  string
	TargetSubnetID string
	BastionType    string
	ClientCIDRList []string
	MaxSessionTTL  int
	FreeformTags   map[string]string
	DefinedTags    map[string]map[string]interface{}
}

// CreateSessionRequest represents the request to create a bastion session
type CreateSessionRequest struct {
	BastionID               string
	DisplayName             string
	SessionType             string
	SessionTTL              int
	TargetResourceID        string
	TargetResourcePort      int
	TargetResourcePrivateIP string
	TargetResourceFQDN      string
	SSHPublicKeyFile        string
	SSHPublicKeyContent     string
}
