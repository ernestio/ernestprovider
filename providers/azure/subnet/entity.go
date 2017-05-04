/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package subnet

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
)

// Event : This is the Ernest representation of an azure subnet
type Event struct {
	event.Base
	ID                     string            `json:"id"`
	Name                   string            `json:"name" validate:"required"`
	ResourceGroupName      string            `json:"resource_group_name" validate:"required"`
	VirtualNetworkName     string            `json:"virtual_network_name" validate:"required"`
	AddressPrefix          string            `json:"address_prefix"  validate:"required"`
	NetworkSecurityGroup   string            `json:"network_security_group"`
	NetworkSecurityGroupID string            `json:"network_security_group_id"`
	RouteTable             string            `json:"route_table_id"`
	IPConfigurations       []string          `json:"ip_configurations"`
	ClientID               string            `json:"azure_client_id"`
	ClientSecret           string            `json:"azure_client_secret"`
	TenantID               string            `json:"azure_tenant_id"`
	SubscriptionID         string            `json:"azure_subscription_id"`
	Environment            string            `json:"environment"`
	ErrorMessage           string            `json:"error,omitempty"`
	Components             []json.RawMessage `json:"components"`
	CryptoKey              string            `json:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev azure.Resource
	ev = &Event{CryptoKey: cryptoKey}
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_subnet", body, val, ev)
}

// SetComponents : ....
func (ev *Event) SetComponents(components []event.Event) {
	for _, v := range components {
		ev.Components = append(ev.Components, v.GetBody())
	}
}

// ValidateID : determines if the given id is valid for this resource type
func (ev *Event) ValidateID(id string) bool {
	parts := strings.Split(strings.ToLower(id), "/")
	if len(parts) != 11 {
		return false
	}
	if parts[6] != "microsoft.network" {
		return false
	}
	if parts[9] != "subnets" {
		return false
	}
	return true
}

// SetID : id setter
func (ev *Event) SetID(id string) {
	ev.ID = id
}

// GetID : id getter
func (ev *Event) GetID() string {
	return ev.ID
}

// SetState : state setter
func (ev *Event) SetState(state string) {
	ev.State = state
}

// ResourceDataToEvent : Translates a ResourceData on a valid Ernest Event
func (ev *Event) ResourceDataToEvent(d *schema.ResourceData) error {
	ev.Name = d.Get("name").(string)
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.VirtualNetworkName = d.Get("virtual_network_name").(string)
	ev.AddressPrefix = d.Get("address_prefix").(string)
	ev.NetworkSecurityGroupID = d.Get("network_security_group_id").(string)
	ev.RouteTable = d.Get("route_table_id").(string)

	configs := []string{}
	for _, config := range d.Get("ip_configurations").(*schema.Set).List() {
		configs = append(configs, config.(string))
	}
	ev.IPConfigurations = configs

	return nil
}

// EventToResourceData : Translates the current event on a valid ResourceData
func (ev *Event) EventToResourceData(d *schema.ResourceData) error {
	crypto := aes.New()

	encFields := make(map[string]string)
	encFields["subscription_id"] = ev.SubscriptionID
	encFields["client_id"] = ev.ClientID
	encFields["client_secret"] = ev.ClientSecret
	encFields["tenant_id"] = ev.TenantID
	encFields["environment"] = ev.Environment
	for k, v := range encFields {
		dec, err := crypto.Decrypt(v, ev.CryptoKey)
		if err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
		if err := d.Set(k, dec); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
	}

	fields := make(map[string]interface{})
	fields["name"] = ev.Name
	fields["resource_group_name"] = ev.ResourceGroupName
	fields["virtual_network_name"] = ev.VirtualNetworkName
	fields["address_prefix"] = ev.AddressPrefix
	fields["network_security_group_id"] = ev.NetworkSecurityGroupID
	fields["route_table_id"] = ev.RouteTable
	fields["ip_configurations"] = ev.IPConfigurations
	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
	}

	return nil
}

// Error : will mark the event as errored
func (ev *Event) Error(err error) {
	ev.ErrorMessage = err.Error()
	ev.Body, err = json.Marshal(ev)
}
