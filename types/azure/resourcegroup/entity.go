package resourcegroup

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                string            `json:"id" diff:"-"`
	Name              string            `json:"name" validate:"required" diff:"-"`
	Location          string            `json:"location" validate:"required" diff:"location"`
	ResourceGroupName string            `json:"resource_group_name,omitempty" diff:"-"`
	Tags              map[string]string `json:"tags" diff:"tags"`
	ClientID          string            `json:"azure_client_id" diff:"-"`
	ClientSecret      string            `json:"azure_client_secret" diff:"-"`
	TenantID          string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string            `json:"azure_subscription_id" diff:"-"`
	Environment       string            `json:"environment" diff:"-"`
	Components        []json.RawMessage `json:"components" diff:"-"`
}
