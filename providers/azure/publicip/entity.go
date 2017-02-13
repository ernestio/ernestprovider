/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package publicip

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
)

// Event : This is the Ernest representation of an azure publicip
type Event struct {
	event.Base
	ID                        string            `json:"id"`
	Name                      string            `json:"name" validate:"required"`
	Location                  string            `json:"location" validate:"required"`
	ResourceGroupName         string            `json:"resource_group_name" validate:"required"`
	PublicIPAddressAllocation string            `json:"public_ip_address_allocation" validate:"required"`
	IdleTimeoutInMinutes      int               `json:"idle_timeout_in_minutes"`
	DomainNameLabel           string            `json:"domain_name_label"`
	ReverseFQDN               string            `json:"reverse_fqdn"`
	FQDN                      string            `json:"fqdn"`
	IP                        string            `json:"ip_address"`
	Tags                      map[string]string `json:"tags"`
	ClientID                  string            `json:"azure_client_id"`
	ClientSecret              string            `json:"azure_client_secret"`
	TenantID                  string            `json:"azure_tenant_id"`
	SubscriptionID            string            `json:"azure_subscription_id"`
	Environment               string            `json:"environment"`
	ErrorMessage              string            `json:"error,omitempty"`
	CryptoKey                 string            `json:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev azure.Resource
	ev = &Event{CryptoKey: cryptoKey}
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_public_ip", body, val, ev)
}

// SetID : id setter
func (ev *Event) SetID(id string) {
	ev.ID = id
}

// GetID : id getter
func (ev *Event) GetID() string {
	return ev.ID
}

// ResourceDataToEvent : Translates a ResourceData on a valid Ernest Event
func (ev *Event) ResourceDataToEvent(d *schema.ResourceData) error {
	ev.Name = d.Get("name").(string)
	ev.Location = d.Get("location").(string)
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.PublicIPAddressAllocation = d.Get("public_ip_address_allocation").(string)
	ev.IdleTimeoutInMinutes = d.Get("idle_timeout_in_minutes").(int)
	ev.DomainNameLabel = d.Get("domain_name_label").(string)
	ev.ReverseFQDN = d.Get("reverse_fqdn").(string)
	ev.FQDN = d.Get("fqdn").(string)
	ev.IP = d.Get("ip_address").(string)
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
	fields["public_ip_address_allocation"] = ev.PublicIPAddressAllocation
	fields["idle_timeout_in_minutes"] = ev.IdleTimeoutInMinutes
	fields["domain_name_label"] = ev.DomainNameLabel
	fields["reverse_fqdn"] = ev.ReverseFQDN
	fields["fqdn"] = ev.FQDN
	fields["ip_address"] = ev.IP
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

// Error : will mark the event as errored
func (ev *Event) Error(err error) {
	ev.ErrorMessage = err.Error()
	ev.Body, err = json.Marshal(ev)
}
