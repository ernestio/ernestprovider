package manageddisk

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                 string            `json:"id" diff:"-"`
	Name               string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName  string            `json:"resource_group_name" validate:"required" diff:"-"`
	Location           string            `json:"location" diff:"-"`
	StorageAccountType string            `json:"storage_account_type" diff:"storage_account_type"`
	CreateOption       string            `json:"create_option" diff:"create_option"`
	SourceURI          string            `json:"source_uri" diff:"source_uri"`
	SourceResourceID   string            `json:"source_resource_id" diff:"-"`
	OSType             string            `json:"os_type" diff:"os_type"`
	DiskSizeGB         int32             `json:"disk_size_gb" diff:"disk_size_gb"`
	Tags               map[string]string `json:"tags" diff:"-"`
	ClientID           string            `json:"azure_client_id" diff:"-"`
	ClientSecret       string            `json:"azure_client_secret" diff:"-"`
	TenantID           string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID     string            `json:"azure_subscription_id" diff:"-"`
	Environment        string            `json:"environment" diff:"-"`
	Components         []json.RawMessage `json:"components" diff:"-"`
}
