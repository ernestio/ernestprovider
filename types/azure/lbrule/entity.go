package lbrule

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

// Event : This is the Ernest representation of an azure lb
type Event struct {
	types.Base
	ID                          string            `json:"id" diff:"-"`
	Name                        string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName           string            `json:"resource_group_name" validate:"required" diff:"-"`
	Loadbalancer                string            `json:"loadbalancer" diff:"loadbalancer,immutable"`
	LoadbalancerID              string            `json:"loadbalancer_id" diff:"-"`
	FrontendIPConfigurationName string            `json:"frontend_ip_configuration_name" diff:"frontend_ip_configuration_name"`
	Protocol                    string            `json:"protocol" diff:"protocol"`
	FrontendPort                int               `json:"frontend_port" diff:"frontend_port"`
	BackendPort                 int               `json:"backend_port" diff:"backend_port"`
	BackendAddressPool          string            `json:"backend_address_pool" diff:"backend_address_pool"`
	BackendAddressPoolID        string            `json:"backend_address_pool_id" diff:"-"`
	Probe                       string            `json:"probe" diff:"probe"`
	ProbeID                     string            `json:"probe_id" diff:"-"`
	EnableFloatingIP            bool              `json:"enable_floating_ip" diff:"enable_floating_ip"`
	IdleTimeoutInMinutes        int               `json:"idle_timeout_in_minutes" diff:"idle_timeout_in_minutes"`
	LoadDistribution            string            `json:"load_distribution" diff:"load_distribution"`
	ClientID                    string            `json:"azure_client_id" diff:"-"`
	ClientSecret                string            `json:"azure_client_secret" diff:"-"`
	TenantID                    string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID              string            `json:"azure_subscription_id" diff:"-"`
	Environment                 string            `json:"environment" diff:"-"`
	Components                  []json.RawMessage `json:"components" diff:"-"`
}
