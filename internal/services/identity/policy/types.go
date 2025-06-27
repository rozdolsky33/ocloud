package policy

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service encapsulates operations for managing and querying policy resources within a specific compartment.
// It utilizes an IdentityClient for API interactions and supports operations like list and search.
// Logger provides event logging functionality, while CompartmentID specifies the compartment for policy retrieval.
type Service struct {
	identityClient identity.IdentityClient
	logger         logr.Logger
	CompartmentID  string
}

// Policy represents a set of rules defining access permissions for resources within a system.
type Policy struct {
	Name        string
	ID          string
	Statement   []string
	Description string
	PolicyTags  util.ResourceTags
}

// IndexablePolicy represents a policy structure optimized for indexing and searching operations.
// It includes fields for policy name, description, statements, and flattened tag information.
type IndexablePolicy struct {
	Name        string
	Description string
	Statement   string
	Tags        string
	TagValues   string
}
