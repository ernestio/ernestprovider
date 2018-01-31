package subnet

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                     string            `json:"id" diff:"-"`
	Name                   string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName      string            `json:"resource_group_name" validate:"required" diff:"-"`
	VirtualNetworkName     string            `json:"virtual_network_name" validate:"required" diff:"virtual_network_name,immutable"`
	AddressPrefix          string            `json:"address_prefix"  validate:"required" diff:"address_prefix,immutable"`
	NetworkSecurityGroup   string            `json:"network_security_group" diff:"network_security_group"`
	NetworkSecurityGroupID string            `json:"network_security_group_id" diff:"-"`
	RouteTable             string            `json:"route_table_id" diff:"-"`
	IPConfigurations       []string          `json:"ip_configurations" diff:"ip_configurations,immutable"`
	ClientID               string            `json:"azure_client_id" diff:"-"`
	ClientSecret           string            `json:"azure_client_secret" diff:"-"`
	TenantID               string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID         string            `json:"azure_subscription_id" diff:"-"`
	Environment            string            `json:"environment" diff:"-"`
	Components             []json.RawMessage `json:"components" diff:"-"`
}
