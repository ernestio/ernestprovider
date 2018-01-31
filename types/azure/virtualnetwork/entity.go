package virtualnetwork

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                string            `json:"id" diff:"-"`
	Name              string            `json:"name" validate:"required" diff:"-"`
	AddressSpace      []string          `json:"address_space" validate:"min=1" diff:"address_space,immutable"`
	DNSServerNames    []string          `json:"dns_server_names" validate:"dive,ip" diff:"dns_server_names"`
	Subnets           []Subnet          `json:"subnets" validate:"min=1" diff:"subnets"`
	Location          string            `json:"location" diff:"-"`
	ResourceGroupName string            `json:"resource_group_name" diff:"-"`
	Tags              map[string]string `json:"tags" diff:"tags"`
	ClientID          string            `json:"azure_client_id" diff:"-"`
	ClientSecret      string            `json:"azure_client_secret" diff:"-"`
	TenantID          string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string            `json:"azure_subscription_id" diff:"-"`
	Environment       string            `json:"environment" diff:"-"`
	Components        []json.RawMessage `json:"components" diff:"-"`
}

// Subnet ..
type Subnet struct {
	Name              string `json:"name" diff:"name"`
	AddressPrefix     string `json:"address_prefix" diff:"address_prefix"`
	SecurityGroupName string `json:"security_group_name" diff:"security_group_name"`
	SecurityGroup     string `json:"security_group" diff:"security_group"`
}
