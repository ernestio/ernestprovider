package lb

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                       string                    `json:"id" diff:"-"`
	Name                     string                    `json:"name" validate:"required" diff:"-"`
	ResourceGroupName        string                    `json:"resource_group_name" validate:"required" diff:"-"`
	Location                 string                    `json:"location" diff:"location"`
	FrontendIPConfigurations []FrontendIPConfiguration `json:"frontend_ip_configurations" validate:"required" diff:"frontend_ip_configurations"`
	Tags                     map[string]string         `json:"tags" diff:"tags"`
	ClientID                 string                    `json:"azure_client_id" diff:"-"`
	ClientSecret             string                    `json:"azure_client_secret" diff:"-"`
	TenantID                 string                    `json:"azure_tenant_id" diff:"-"`
	SubscriptionID           string                    `json:"azure_subscription_id" diff:"-"`
	Environment              string                    `json:"environment" diff:"-"`
	Components               []json.RawMessage         `json:"components" diff:"-"`
}

// FrontendIPConfiguration ...
type FrontendIPConfiguration struct {
	Name                       string `json:"name" validate:"required" structs:"name"diff:"name"`
	Subnet                     string `json:"subnet" structs:"-" diff:"subnet"`
	SubnetID                   string `json:"subnet_id" structs:"subnet_id" diff:"-"`
	PrivateIPAddress           string `json:"private_ip_address" structs:"private_ip_address" diff:"private_ip_address"`
	PrivateIPAddressAllocation string `json:"private_ip_address_allocation" structs:"private_ip_address_allocation" diff:"private_ip_address_allocation"`
	PublicIPAddress            string `json:"public_ip_address" structs:"-" diff:"public_ip_address"`
	PublicIPAddressID          string `json:"public_ip_address_id" structs:"public_ip_address_id" diff:"-"`
}
