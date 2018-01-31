package lbbackendaddresspool

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                string            `json:"id" diff:"-"`
	Name              string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName string            `json:"resource_group_name" validate:"required" diff:"-"`
	Loadbalancer      string            `json:"loadbalancer" diff:"loadbalancer,immutable"`
	LoadbalancerID    string            `json:"loadbalancer_id" diff:"-"`
	ClientID          string            `json:"azure_client_id" diff:"-"`
	ClientSecret      string            `json:"azure_client_secret" diff:"-"`
	TenantID          string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string            `json:"azure_subscription_id" diff:"-"`
	Environment       string            `json:"environment" diff:"-"`
	Components        []json.RawMessage `json:"components" diff:"-"`
}
