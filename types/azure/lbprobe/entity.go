package lbprobe

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                string            `json:"id" diff:"-"`
	Name              string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName string            `json:"resource_group_name" validate:"required" diff:"-"`
	Loadbalancer      string            `json:"loadbalancer" diff:"-"`
	LoadbalancerID    string            `json:"loadbalancer_id" diff:"-"`
	Protocol          string            `json:"protocol" diff:"protocol"`
	Port              int               `json:"port" diff:"port"`
	RequestPath       string            `json:"request_path" diff:"request_path"`
	IntervalInSeconds int               `json:"interval_in_seconds" diff:"interval_in_seconds"`
	NumberOfProbes    int               `json:"number_of_probes" diff:"number_of_probes"`
	ClientID          string            `json:"azure_client_id" diff:"-"`
	ClientSecret      string            `json:"azure_client_secret" diff:"-"`
	TenantID          string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string            `json:"azure_subscription_id" diff:"-"`
	Environment       string            `json:"environment" diff:"-"`
	Components        []json.RawMessage `json:"components" diff:"-"`
}
