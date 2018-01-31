package publicip

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                        string            `json:"id" diff:"-"`
	Name                      string            `json:"name" validate:"required" diff:"-"`
	Location                  string            `json:"location" validate:"required" diff:"location"`
	ResourceGroupName         string            `json:"resource_group_name" validate:"required" diff:"-"`
	LoadBalancer              string            `json:"lb" diff:"-"`
	PublicIPAddressAllocation string            `json:"public_ip_address_allocation" validate:"required" diff:"public_ip_address_allocation,immutable"`
	IdleTimeoutInMinutes      int               `json:"idle_timeout_in_minutes" diff:"idle_timeout_in_minutes,immutable"`
	DomainNameLabel           string            `json:"domain_name_label" diff:"domain_name_label,immutable"`
	ReverseFQDN               string            `json:"reverse_fqdn" diff:"reverse_fqdn,immutable"`
	FQDN                      string            `json:"fqdn" diff:"fqdn,immutable"`
	IP                        string            `json:"ip_address" diff:"ip_address,immutable"`
	Tags                      map[string]string `json:"tags" diff:"-"`
	ClientID                  string            `json:"azure_client_id" diff:"-"`
	ClientSecret              string            `json:"azure_client_secret" diff:"-"`
	TenantID                  string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID            string            `json:"azure_subscription_id" diff:"-"`
	Environment               string            `json:"environment" diff:"-"`
	Components                []json.RawMessage `json:"components" diff:"-"`
}
