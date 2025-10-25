package mapping

import (
	"reflect"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
)

// BastionAttributes represents the attributes extracted from OCI SDK Bastion types.
// This intermediate structure handles OCI SDK types and provides safe conversion to domain models.
type BastionAttributes struct {
	OCID                     *string
	DisplayName              *string
	BastionType              *string
	LifecycleState           bastion.BastionLifecycleStateEnum
	CompartmentID            *string
	TargetVcnID              *string
	TargetSubnetID           *string
	MaxSessionTTL            *int
	ClientCidrBlockAllowList []string
	PrivateEndpointIpAddress *string
	FreeformTags             map[string]string
	DefinedTags              map[string]map[string]interface{}
	TimeCreated              interface{}
	TimeUpdated              interface{}
}

// BastionSessionAttributes represents the attributes extracted from OCI SDK Session types.
type BastionSessionAttributes struct {
	OCID                    *string
	DisplayName             *string
	BastionID               *string
	BastionName             *string
	LifecycleState          bastion.SessionLifecycleStateEnum
	SessionType             *string
	SessionTTL              *int
	TargetResourceID        *string
	TargetResourcePort      *int
	TargetResourcePrivateIP *string
	TargetResourceFQDN      *string
	SSHMetadata             map[string]string
	TimeCreated             interface{}
	TimeUpdated             interface{}
}

// NewBastionAttributesFromOCIBastionSummary creates BastionAttributes from OCI BastionSummary.
// Note: BastionSummary doesn't include MaxSessionTtlInSeconds - use full Bastion for that field.
func NewBastionAttributesFromOCIBastionSummary(b bastion.BastionSummary) *BastionAttributes {
	return &BastionAttributes{
		OCID:           b.Id,
		DisplayName:    b.Name,
		BastionType:    b.BastionType,
		LifecycleState: b.LifecycleState,
		CompartmentID:  b.CompartmentId,
		TargetVcnID:    b.TargetVcnId,
		TargetSubnetID: b.TargetSubnetId,
		MaxSessionTTL:  nil, // Not available in BastionSummary
		FreeformTags:   b.FreeformTags,
		DefinedTags:    b.DefinedTags,
		TimeCreated:    b.TimeCreated,
		TimeUpdated:    b.TimeUpdated,
	}
}

// NewBastionAttributesFromOCIBastion creates BastionAttributes from OCI Bastion (full object).
func NewBastionAttributesFromOCIBastion(b bastion.Bastion) *BastionAttributes {
	return &BastionAttributes{
		OCID:                     b.Id,
		DisplayName:              b.Name,
		BastionType:              b.BastionType,
		LifecycleState:           b.LifecycleState,
		CompartmentID:            b.CompartmentId,
		TargetVcnID:              b.TargetVcnId,
		TargetSubnetID:           b.TargetSubnetId,
		MaxSessionTTL:            b.MaxSessionTtlInSeconds,
		ClientCidrBlockAllowList: b.ClientCidrBlockAllowList,
		PrivateEndpointIpAddress: b.PrivateEndpointIpAddress,
		FreeformTags:             b.FreeformTags,
		DefinedTags:              b.DefinedTags,
		TimeCreated:              b.TimeCreated,
		TimeUpdated:              b.TimeUpdated,
	}
}

// NewBastionSessionAttributesFromOCISession creates BastionSessionAttributes from OCI Session.
func NewBastionSessionAttributesFromOCISession(s bastion.Session) *BastionSessionAttributes {
	// Extract target resource details to determine session type
	var targetResourceID, targetResourceFQDN, targetResourcePrivateIP *string
	var targetResourcePort *int
	var sessionType *string

	if s.TargetResourceDetails != nil {
		// Determine session type based on TargetResourceDetails type
		switch details := s.TargetResourceDetails.(type) {
		case bastion.PortForwardingSessionTargetResourceDetails:
			sessionType = stringPtr("PORT_FORWARDING")
			targetResourceID = details.TargetResourceId
			targetResourcePort = details.TargetResourcePort
			targetResourcePrivateIP = details.TargetResourcePrivateIpAddress
		case bastion.ManagedSshSessionTargetResourceDetails:
			sessionType = stringPtr("MANAGED_SSH")
			targetResourceID = details.TargetResourceId
			targetResourcePort = details.TargetResourcePort
		case bastion.DynamicPortForwardingSessionTargetResourceDetails:
			sessionType = stringPtr("DYNAMIC_PORT_FORWARDING")
		}
	}

	return &BastionSessionAttributes{
		OCID:                    s.Id,
		DisplayName:             s.DisplayName,
		BastionID:               s.BastionId,
		BastionName:             s.BastionName,
		LifecycleState:          s.LifecycleState,
		SessionType:             sessionType,
		SessionTTL:              s.SessionTtlInSeconds,
		TargetResourceID:        targetResourceID,
		TargetResourcePort:      targetResourcePort,
		TargetResourcePrivateIP: targetResourcePrivateIP,
		TargetResourceFQDN:      targetResourceFQDN,
		SSHMetadata:             s.SshMetadata,
		TimeCreated:             s.TimeCreated,
		TimeUpdated:             s.TimeUpdated,
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// NewDomainBastionFromAttrs converts BastionAttributes to domain.Bastion.
// This function handles nil-pointer safety and type conversions.
func NewDomainBastionFromAttrs(attrs *BastionAttributes) *domain.Bastion {
	var ocid, displayName, bType, lifecycleState, compartmentID string
	var targetVcnID, targetSubnetID string
	var maxSessionTTL int

	if attrs.OCID != nil {
		ocid = *attrs.OCID
	}
	if attrs.DisplayName != nil {
		displayName = *attrs.DisplayName
	}
	if attrs.BastionType != nil {
		bType = *attrs.BastionType
	}
	if attrs.LifecycleState != "" {
		lifecycleState = string(attrs.LifecycleState)
	}
	if attrs.CompartmentID != nil {
		compartmentID = *attrs.CompartmentID
	}
	if attrs.TargetVcnID != nil {
		targetVcnID = *attrs.TargetVcnID
	}
	if attrs.TargetSubnetID != nil {
		targetSubnetID = *attrs.TargetSubnetID
	}
	if attrs.MaxSessionTTL != nil {
		maxSessionTTL = *attrs.MaxSessionTTL
	}

	var timeCreated, timeUpdated string
	if attrs.TimeCreated != nil {
		timeCreated = formatOCITime(attrs.TimeCreated)
	}
	if attrs.TimeUpdated != nil {
		timeUpdated = formatOCITime(attrs.TimeUpdated)
	}

	var privateEndpointIpAddress string
	if attrs.PrivateEndpointIpAddress != nil {
		privateEndpointIpAddress = *attrs.PrivateEndpointIpAddress
	}

	return &domain.Bastion{
		OCID:                     ocid,
		DisplayName:              displayName,
		BastionType:              bType,
		LifecycleState:           lifecycleState,
		CompartmentID:            compartmentID,
		TargetVcnID:              targetVcnID,
		TargetSubnetID:           targetSubnetID,
		MaxSessionTTL:            maxSessionTTL,
		ClientCidrBlockAllowList: attrs.ClientCidrBlockAllowList,
		PrivateEndpointIpAddress: privateEndpointIpAddress,
		FreeformTags:             attrs.FreeformTags,
		DefinedTags:              attrs.DefinedTags,
		TimeCreated:              timeCreated,
		TimeUpdated:              timeUpdated,
	}
}

// NewDomainBastionSessionFromAttrs converts BastionSessionAttributes to domain.BastionSession.
func NewDomainBastionSessionFromAttrs(attrs *BastionSessionAttributes) *domain.BastionSession {
	var ocid, displayName, bastionID, bastionName string
	var lifecycleState, sessionType string
	var sessionTTL, targetResourcePort int
	var targetResourceID, targetResourcePrivateIP, targetResourceFQDN string

	if attrs.OCID != nil {
		ocid = *attrs.OCID
	}
	if attrs.DisplayName != nil {
		displayName = *attrs.DisplayName
	}
	if attrs.BastionID != nil {
		bastionID = *attrs.BastionID
	}
	if attrs.BastionName != nil {
		bastionName = *attrs.BastionName
	}
	if attrs.LifecycleState != "" {
		lifecycleState = string(attrs.LifecycleState)
	}
	if attrs.SessionType != nil {
		sessionType = *attrs.SessionType
	}
	if attrs.SessionTTL != nil {
		sessionTTL = *attrs.SessionTTL
	}
	if attrs.TargetResourceID != nil {
		targetResourceID = *attrs.TargetResourceID
	}
	if attrs.TargetResourcePort != nil {
		targetResourcePort = *attrs.TargetResourcePort
	}
	if attrs.TargetResourcePrivateIP != nil {
		targetResourcePrivateIP = *attrs.TargetResourcePrivateIP
	}
	if attrs.TargetResourceFQDN != nil {
		targetResourceFQDN = *attrs.TargetResourceFQDN
	}

	var timeCreated, timeUpdated string
	if attrs.TimeCreated != nil {
		timeCreated = formatOCITime(attrs.TimeCreated)
	}
	if attrs.TimeUpdated != nil {
		timeUpdated = formatOCITime(attrs.TimeUpdated)
	}

	return &domain.BastionSession{
		OCID:                    ocid,
		DisplayName:             displayName,
		BastionID:               bastionID,
		BastionName:             bastionName,
		LifecycleState:          lifecycleState,
		SessionType:             sessionType,
		SessionTTL:              sessionTTL,
		TargetResourceID:        targetResourceID,
		TargetResourcePort:      targetResourcePort,
		TargetResourcePrivateIP: targetResourcePrivateIP,
		TargetResourceFQDN:      targetResourceFQDN,
		SSHMetadata:             attrs.SSHMetadata,
		TimeCreated:             timeCreated,
		TimeUpdated:             timeUpdated,
	}
}

// formatOCITime converts OCI SDK time interface to string.
// OCI SDK uses *common.SDKTime which implements String().
func formatOCITime(t interface{}) string {
	if t == nil {
		return ""
	}
	// Use reflection to check if the underlying value is nil
	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return ""
	}
	// Now safe to call String()
	if stringer, ok := t.(interface{ String() string }); ok {
		return stringer.String()
	}
	return ""
}
