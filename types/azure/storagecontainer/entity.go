package storagecontainer

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                  string            `json:"id" diff:"-"`
	Name                string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName   string            `json:"resource_group_name" validate:"required" diff:"-"`
	StorageAccountName  string            `json:"storage_account_name" validate:"required" diff:"storage_account_name,immutable"`
	ContainerAccessType string            `json:"container_access_type" diff:"container_access_type,immutable"`
	Properties          map[string]string `json:"properties" diff:"properties,immutable"`
	ClientID            string            `json:"azure_client_id" diff:"-"`
	ClientSecret        string            `json:"azure_client_secret" diff:"-"`
	TenantID            string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID      string            `json:"azure_subscription_id" diff:"-"`
	Environment         string            `json:"environment" diff:"-"`
	Components          []json.RawMessage `json:"components" diff:"-"`
}
