package storageaccount

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                     string            `json:"id" diff:"-"`
	Name                   string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName      string            `json:"resource_group_name" validate:"required" diff:"-"`
	Location               string            `json:"location" validate:"required" diff:"-"`
	AccountKind            string            `json:"account_kind" diff:"account_kind"`
	AccountType            string            `json:"account_type" validate:"required" diff:"account_type"`
	PrimaryLocation        string            `json:"primary_location" diff:"-"`
	SecondaryLocation      string            `json:"secondary_location" diff:"-"`
	PrimaryBlobEndpoint    string            `json:"primary_blob_endpoint" diff:"-"`
	SecondaryBlobEndpoint  string            `json:"secondary_blob_endpoint" diff:"-"`
	PrimaryQueueEndpoint   string            `json:"primary_queue_endpoint" diff:"-"`
	SecondaryQueueEndpoint string            `json:"secondary_queue_endpoint" diff:"-"`
	PrimaryTableEndpoint   string            `json:"primary_table_endpoint" diff:"-"`
	SecondaryTableEndpoint string            `json:"secondary_table_endpoint" diff:"-"`
	PrimaryFileEndpoint    string            `json:"primary_file_endpoint" diff:"-"`
	PrimaryAccessKey       string            `json:"primary_access_key" diff:"-"`
	SecondaryAccessKey     string            `json:"secondary_access_key" diff:"-"`
	EnableBlobEncryption   bool              `json:"enable_blob_encryption" diff:"enable_blob_encryption"`
	Tags                   map[string]string `json:"tags" diff:"tags"`
	ClientID               string            `json:"azure_client_id" diff:"-"`
	ClientSecret           string            `json:"azure_client_secret" diff:"-"`
	TenantID               string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID         string            `json:"azure_subscription_id" diff:"-"`
	Environment            string            `json:"environment" diff:"-"`
	Components             []json.RawMessage `json:"components" diff:"-"`
}
