/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package lb

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
	"github.com/fatih/structs"
)

// Event : This is the Ernest representation of an azure lb
type Event struct {
	event.Base
	ID                       string                    `json:"id"`
	Name                     string                    `json:"name" validate:"required"`
	ResourceGroupName        string                    `json:"resource_group_name" validate:"required"`
	Location                 string                    `json:"location"`
	FrontendIPConfigurations []FrontendIPConfiguration `json:"frontend_ip_configurations" validate:"required"`
	Tags                     map[string]string         `json:"tags"`
	ClientID                 string                    `json:"azure_client_id"`
	ClientSecret             string                    `json:"azure_client_secret"`
	TenantID                 string                    `json:"azure_tenant_id"`
	SubscriptionID           string                    `json:"azure_subscription_id"`
	Environment              string                    `json:"environment"`
	ErrorMessage             string                    `json:"error,omitempty"`
	Components               []json.RawMessage         `json:"components"`
	CryptoKey                string                    `json:"-"`
	Validator                *event.Validator          `json:"-"`
}

// FrontendIPConfiguration ...
type FrontendIPConfiguration struct {
	Name                       string `json:"name" validate:"required" structs:"name"`
	Subnet                     string `json:"subnet" structs:"-"`
	SubnetID                   string `json:"subnet_id" structs:"subnet_id"`
	PrivateIPAddress           string `json:"private_ip_address" structs:"private_ip_address"`
	PrivateIPAddressAllocation string `json:"private_ip_address_allocation" structs:"private_ip_address_allocation"`
	PublicIPAddress            string `json:"public_ip_address" structs:"-"`
	PublicIPAddressID          string `json:"public_ip_address_id" structs:"public_ip_address_id"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev event.Resource
	ev = &Event{CryptoKey: cryptoKey, Validator: val}
	body = []byte(strings.Replace(string(body), `"_component":"lbs"`, `"_component":"lb"`, 1))
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_lb", body, val, ev)
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
	if len(parts) != 9 {
		return false
	}
	if parts[6] != "microsoft.network" {
		return false
	}
	if parts[7] != "loadbalancers" {
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
	idParts := strings.Split(d.Id(), "/")
	ev.ID = d.Id()
	ev.Name = idParts[len(idParts)-1]
	ev.ComponentID = "lb::" + ev.Name
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.Location = d.Get("location").(string)
	if d.Content["frontend_ip_configuration"] != nil {
		ips := []FrontendIPConfiguration{}
		for _, c := range d.Content["frontend_ip_configuration"].([]interface{}) {
			cfg := c.(map[string]interface{})

			ips = append(ips, FrontendIPConfiguration{
				SubnetID:          fmt.Sprintf("%s", cfg["subnet_id"]),
				Name:              fmt.Sprintf("%s", cfg["name"]),
				PrivateIPAddress:  fmt.Sprintf("%s", cfg["private_ip_address"]),
				PublicIPAddressID: fmt.Sprintf("%s", cfg["public_ip_address_id"]),
			})
		}
		ev.FrontendIPConfigurations = ips
	}

	tags := make(map[string]string, 0)
	for k, v := range d.Get("tags").(map[string]interface{}) {
		tags[k] = v.(string)
	}
	ev.Tags = tags

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
	fields["location"] = ev.Location
	fields["resource_group_name"] = ev.ResourceGroupName
	var configs []interface{}
	for _, c := range ev.FrontendIPConfigurations {
		configs = append(configs, structs.Map(c))
	}
	fields["frontend_ip_configuration"] = configs
	fields["tags"] = ev.Tags
	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
	}

	return nil
}

// Clone : will mark the event as errored
func (ev *Event) Clone() (event.Event, error) {
	body, _ := json.Marshal(ev)
	return New(ev.Subject, ev.CryptoKey, body, ev.Validator)
}

// Error : will mark the event as errored
func (ev *Event) Error(err error) {
	ev.ErrorMessage = err.Error()
	ev.Body, err = json.Marshal(ev)
}
