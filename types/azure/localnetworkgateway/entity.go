package localnetworkgateway

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
	"github.com/ernestio/ernestprovider/validator"
)

type Event struct {
	types.Base
	ID                string               `json:"id"`
	Name              string               `json:"name" validate:"required" diff:"-"`
	ResourceGroupName string               `json:"resource_group_name" validate:"required" diff:"-"`
	Location          string               `json:"location" validate:"required" diff:"-"`
	GatewayAddress    string               `json:"gateway_address" validate:"required" diff:"gateway_address,immutable"`
	AddressSpace      []string             `json:"address_space" diff:"address_space,immutable"`
	ClientID          string               `json:"azure_client_id" diff:"-"`
	ClientSecret      string               `json:"azure_client_secret" diff:"-"`
	TenantID          string               `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string               `json:"azure_subscription_id" diff:"-"`
	Environment       string               `json:"environment" diff:"-"`
	ErrorMessage      string               `json:"error,omitempty" diff:"-"`
	Components        []json.RawMessage    `json:"components" diff:"-"`
	CryptoKey         string               `json:"-" diff:"-"`
	Validator         *validator.Validator `json:"-" diff:"-"`
}
