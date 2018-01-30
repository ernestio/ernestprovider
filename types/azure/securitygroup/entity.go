package securitygroup

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                string            `json:"id" diff:"-"`
	Name              string            `json:"name" validate:"required" diff:"-"`
	Location          string            `json:"location" validate:"required" diff:"-"`
	ResourceGroupName string            `json:"resource_group_name" validate:"required" diff:"-"`
	SecurityRules     []SecurityRule    `json:"security_rules" diff:"security_rules"`
	Tags              map[string]string `json:"tags" diff:"tags"`
	ClientID          string            `json:"azure_client_id" diff:"-"`
	ClientSecret      string            `json:"azure_client_secret" diff:"-"`
	TenantID          string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string            `json:"azure_subscription_id" diff:"-"`
	Environment       string            `json:"environment" diff:"-"`
	Components        []json.RawMessage `json:"components" diff:"-"`
}

// SecurityRule ...
type SecurityRule struct {
	Name                     string `json:"name" diff:"name"`
	Description              string `json:"description" diff:"description"`
	Protocol                 string `json:"protocol" diff:"protocol"`
	SourcePort               string `json:"source_port_range" diff:"source_port_range"`
	DestinationPortRange     string `json:"destination_port_range" diff:"destination_port_range"`
	SourceAddressPrefix      string `json:"source_address_prefix" diff:"source_address_prefix"`
	DestinationAddressPrefix string `json:"destination_address_prefix" diff:"destination_address_prefix"`
	Access                   string `json:"access" diff:"access"`
	Priority                 int    `json:"priority" diff:"priority"`
	Direction                string `json:"direction" diff:"direction"`
}
