package mapping

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/networkloadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/networklb"
)

type NetworkLoadBalancerAttributes struct {
	ID                      *string
	DisplayName             *string
	LifecycleState          networkloadbalancer.LifecycleStateEnum
	IsPrivate               *bool
	TimeCreated             *time.Time
	IpAddresses             []networkloadbalancer.IpAddress
	Listeners               map[string]networkloadbalancer.Listener
	BackendSets             map[string]networkloadbalancer.BackendSet
	SubnetId                *string
	NetworkSecurityGroupIds []string
}

func NewNetworkLoadBalancerAttributesFromOCI(nlb networkloadbalancer.NetworkLoadBalancer) *NetworkLoadBalancerAttributes {
	return &NetworkLoadBalancerAttributes{
		ID:                      nlb.Id,
		DisplayName:             nlb.DisplayName,
		LifecycleState:          nlb.LifecycleState,
		IsPrivate:               nlb.IsPrivate,
		TimeCreated:             &nlb.TimeCreated.Time,
		IpAddresses:             nlb.IpAddresses,
		Listeners:               nlb.Listeners,
		BackendSets:             nlb.BackendSets,
		SubnetId:                nlb.SubnetId,
		NetworkSecurityGroupIds: nlb.NetworkSecurityGroupIds,
	}
}

func NewNetworkLoadBalancerAttributesFromOCISummary(nlb networkloadbalancer.NetworkLoadBalancerSummary) *NetworkLoadBalancerAttributes {
	return &NetworkLoadBalancerAttributes{
		ID:                      nlb.Id,
		DisplayName:             nlb.DisplayName,
		LifecycleState:          nlb.LifecycleState,
		IsPrivate:               nlb.IsPrivate,
		TimeCreated:             &nlb.TimeCreated.Time,
		IpAddresses:             nlb.IpAddresses,
		Listeners:               nlb.Listeners,
		BackendSets:             nlb.BackendSets,
		SubnetId:                nlb.SubnetId,
		NetworkSecurityGroupIds: nlb.NetworkSecurityGroupIds,
	}
}

func NewDomainNetworkLoadBalancerFromAttrs(nlb *NetworkLoadBalancerAttributes) *domain.NetworkLoadBalancer {
	deref := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}

	id := deref(nlb.ID)
	name := deref(nlb.DisplayName)

	typeStr := "Public"
	if nlb.IsPrivate != nil && *nlb.IsPrivate {
		typeStr = "Private"
	}

	var createdTime *time.Time
	if nlb.TimeCreated != nil {
		t := *nlb.TimeCreated
		createdTime = &t
	}

	ips := make([]string, 0, len(nlb.IpAddresses))
	for i := range nlb.IpAddresses {
		ip := nlb.IpAddresses[i]
		addr := deref(ip.IpAddress)
		if addr == "" {
			continue
		}
		if ip.IsPublic != nil {
			if *ip.IsPublic {
				addr += " (public)"
			} else {
				addr += " (private)"
			}
		}
		ips = append(ips, addr)
	}

	// Listeners
	listeners := make(map[string]string, len(nlb.Listeners))
	for lname, l := range nlb.Listeners {
		port := 0
		if l.Port != nil {
			port = *l.Port
		}
		backend := deref(l.DefaultBackendSetName)
		proto := strings.ToLower(string(l.Protocol))
		val := proto + ":" + strconv.Itoa(port) + " → " + backend
		listeners[lname] = val
	}

	// Backend sets
	backendSets := make(map[string]domain.BackendSet, len(nlb.BackendSets))
	for bsName, bs := range nlb.BackendSets {
		policy := strings.ToUpper(string(bs.Policy))

		hc := ""
		if bs.HealthChecker != nil {
			p := strings.ToUpper(string(bs.HealthChecker.Protocol))
			port := 0
			if bs.HealthChecker.Port != nil {
				port = *bs.HealthChecker.Port
			}
			if p != "" {
				hc = p + ":" + strconv.Itoa(port)
			}
		}

		backendSets[bsName] = domain.BackendSet{
			Policy:   policy,
			Health:   hc,
			Backends: []domain.Backend{},
		}
	}

	// Subnet
	subnetID := deref(nlb.SubnetId)
	subnets := make([]string, 0, 1)
	if subnetID != "" {
		subnets = append(subnets, subnetID)
	}

	nsgs := make([]string, len(nlb.NetworkSecurityGroupIds))
	copy(nsgs, nlb.NetworkSecurityGroupIds)
	sort.Strings(nsgs)

	return &domain.NetworkLoadBalancer{
		ID:            id,
		OCID:          id,
		Name:          name,
		State:         string(nlb.LifecycleState),
		Type:          typeStr,
		IPAddresses:   ips,
		Listeners:     listeners,
		BackendHealth: make(map[string]string),
		Subnets:       subnets,
		NSGs:          nsgs,
		Created:       createdTime,
		BackendSets:   backendSets,
		SubnetID:      subnetID,
	}
}
