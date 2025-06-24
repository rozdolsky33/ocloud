package compartment

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

type Service struct {
	identityClient identity.IdentityClient
	logger         logr.Logger
	TenancyID      string
	TenancyName    string
}

type Compartment struct {
	Name        string
	ID          string
	Description string
}

type IndexableCompartment struct {
	ID          string
	Name        string
	Description string
	Tags        string
	TagValues   string
}
