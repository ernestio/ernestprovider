package networkinterface

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                     string            `json:"id" diff:"-"`
	Name                   string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName      string            `json:"resource_group_name" validate:"required" diff:"resource_group_name"`
	Location               string            `json:"location" diff:"-"`
	NetworkSecurityGroup   string            `json:"network_security_group" diff:"network_security_group"`
	NetworkSecurityGroupID string            `json:"network_security_group_id" diff:"-"`
	MacAddress             string            `json:"mac_address" diff:"-"`
	PrivateIPAddress       string            `json:"private_ip_address" diff:"private_ip_address,immutable"`
	VirtualMachineID       string            `json:"virtual_machine_id" diff:"-"`
	IPConfigurations       []IPConfiguration `json:"ip_configuration" structs:"ip_configuration" diff:"ip_configuration"` // validate:"min=1,dive"`
	DNSServers             []string          `json:"dns_servers" validate:"dive,ip" diff:"dns_servers"`
	InternalDNSNameLabel   string            `json:"internal_dns_name_label" diff:"internal_dns_name_label"`
	AppliedDNSServers      []string          `json:"applied_dns_servers" diff:"-"`
	InternalFQDN           string            `json:"internal_fqdn" diff:"-"`
	EnableIPForwarding     bool              `json:"enable_ip_forwarding" diff:"enable_ip_forwarding,immutable"`
	Tags                   map[string]string `json:"tags" diff:"tags"`
	ClientID               string            `json:"azure_client_id" diff:"-"`
	ClientSecret           string            `json:"azure_client_secret" diff:"-"`
	TenantID               string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID         string            `json:"azure_subscription_id" diff:"-"`
	Environment            string            `json:"environment" diff:"-"`
	Components             []json.RawMessage `json:"components" diff:"-"`
}

// IPConfiguration : ...
type IPConfiguration struct {
	Name                              string   `json:"name" validate:"required" structs:"name" diff:"name"`
	Subnet                            string   `json:"subnet" validate:"required" structs:"-" diff:"-"`
	SubnetID                          string   `json:"subnet_id" validate:"required" structs:"subnet_id" diff:"-"`
	PublicIPAddress                   string   `json:"public_ip_address" structs:"-" diff:"public_ip_address"`
	PrivateIPAddress                  string   `json:"private_ip_address" structs:"private_ip_address" diff:"private_ip_address"`
	PrivateIPAddressAllocation        string   `json:"private_ip_address_allocation" validate:"required" structs:"private_ip_address_allocation" diff:"private_ip_address_allocation"`
	PublicIPAddressID                 string   `json:"public_ip_address_id" structs:"public_ip_address_id" diff:"-"`
	LoadBalancerBackendAddressPools   []string `json:"load_balancer_backend_address_pools" structs:"-" diff:"-"`
	LoadBalancerBackendAddressPoolIDs []string `json:"load_balancer_backend_address_pools_ids" structs:"load_balancer_backend_address_pools_ids,omitempty" diff:"-"`
	LoadBalancerInboundNatRules       []string `json:"load_balancer_inbound_nat_rules_ids" structs:"load_balancer_inbound_nat_rules_ids,omitempty" diff:"-"`
}
