package sqlfirewallrule

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                string            `json:"id" diff:"-"`
	Name              string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName string            `json:"resource_group_name" validate:"required" diff:"-"`
	ServerName        string            `json:"server_name" validate:"required" diff:"server_name,immutable"`
	StartIPAddress    string            `json:"start_ip_address" validate:"required" diff:"start_ip_address"`
	EndIPAddress      string            `json:"end_ip_address" validate:"required" diff:"end_ip_address"`
	ClientID          string            `json:"azure_client_id" diff:"-"`
	ClientSecret      string            `json:"azure_client_secret" diff:"-"`
	TenantID          string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string            `json:"azure_subscription_id" diff:"-"`
	Environment       string            `json:"environment" diff:"-"`
	Components        []json.RawMessage `json:"components" diff:"-"`
}
