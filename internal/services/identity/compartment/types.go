package compartment

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service represents a service that manages operations within a specific tenancy using an identity client and logger.
type Service struct {
	identityClient identity.IdentityClient
	logger         logr.Logger
	TenancyID      string
	TenancyName    string
}

// Compartment represents a logical entity that groups cloud resources within an Oracle Cloud Infrastructure tenancy.
type Compartment struct {
	Name            string
	ID              string
	Description     string
	CompartmentTags util.ResourceTags
}

// IndexableCompartment represents a simplified version of a compartment with fields designed for indexing purposes.
// Name is the lowercased name of the compartment for case-insensitive indexing or searching.
// Description is the lowercased description of the compartment, suitable for indexed queries.
type IndexableCompartment struct {
	Name        string
	Description string
}
