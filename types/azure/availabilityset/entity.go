package availabilityset

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                        string            `json:"id" diff:"-"`
	Name                      string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName         string            `json:"resource_group_name" validate:"required" diff:"-"`
	Location                  string            `json:"location" diff:"location,immutable"`
	PlatformUpdateDomainCount int               `json:"platform_update_domain_count" diff:"platform_update_domain_count"`
	PlatformFaultDomainCount  int               `json:"platform_fault_domain_count" diff:"platform_fault_domain_count"`
	Managed                   bool              `json:"managed" diff:"managed,immutable"`
	Tags                      map[string]string `json:"tags" diff:"-"`
	ClientID                  string            `json:"azure_client_id" diff:"-"`
	ClientSecret              string            `json:"azure_client_secret" diff:"-"`
	TenantID                  string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID            string            `json:"azure_subscription_id" diff:"-"`
	Environment               string            `json:"environment" diff:"-"`
	Components                []json.RawMessage `json:"components" diff:"-"`
}
