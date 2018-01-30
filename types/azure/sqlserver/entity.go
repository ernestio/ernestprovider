package sqlserver

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                         string            `json:"id" diff:"-"`
	Name                       string            `json:"name" validate:"required" diff:"-"`
	Location                   string            `json:"location" validate:"required" diff:"-"`
	ResourceGroupName          string            `json:"resource_group_name" validate:"required" diff:"-"`
	Version                    string            `json:"version" diff:"version"`
	AdministratorLogin         string            `json:"administrator_login" validate:"required" diff:"administrator_login"`
	AdministratorLoginPassword string            `json:"administrator_login_password" diff:"administrator_login_password"`
	FullyQualifiedDomainName   string            `json:"fully_qualified_domain_name" diff:"fully_qualified_domain_name,immutable"`
	Tags                       map[string]string `json:"tags" diff:"-"`
	ClientID                   string            `json:"azure_client_id" diff:"-"`
	ClientSecret               string            `json:"azure_client_secret" diff:"-"`
	TenantID                   string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID             string            `json:"azure_subscription_id" diff:"-"`
	Environment                string            `json:"environment" diff:"-"`
	Components                 []json.RawMessage `json:"components" diff:"-"`
}
