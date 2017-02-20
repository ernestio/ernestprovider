/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualnetwork

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
)

// Event : This is the Ernest representation of an azure virtual
//         network
type Event struct {
	event.Base
	ID                string            `json:"id"`
	Name              string            `json:"name" validate:"required"`
	AddressSpace      []string          `json:"address_space" validate:"min=1"`
	DNSServerNames    []string          `json:"dns_server_names" validate:"dive,ip"`
	Subnets           []subnet          `json:"subnets" validate:"min=1"`
	Location          string            `json:"location"`
	ResourceGroupName string            `json:"resource_group_name"`
	Tags              map[string]string `json:"tags"`
	ClientID          string            `json:"azure_client_id"`
	ClientSecret      string            `json:"azure_client_secret"`
	TenantID          string            `json:"azure_tenant_id"`
	SubscriptionID    string            `json:"azure_subscription_id"`
	Environment       string            `json:"environment"`
	ErrorMessage      string            `json:"error,omitempty"`
	Components        []json.RawMessage `json:"components"`
	CryptoKey         string            `json:"-"`
}

type subnet struct {
	Name          string `json:"name"`
	AddressPrefix string `json:"address_prefix"`
	SecurityGroup string `json:"security_group"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev azure.Resource
	ev = &Event{CryptoKey: cryptoKey}
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_virtual_network", body, val, ev)
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
	if parts[7] != "virtualnetworks" {
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

// ResourceDataToEvent : Translates a ResourceData on a valid Ernest Event
func (ev *Event) ResourceDataToEvent(d *schema.ResourceData) error {
	ev.Name = d.Get("name").(string)

	prefixes := []string{}
	for _, prefix := range d.Get("address_space").([]interface{}) {
		prefixes = append(prefixes, prefix.(string))
	}
	ev.AddressSpace = prefixes

	dnses := []string{}
	for _, dns := range d.Get("dns_servers").([]interface{}) {
		dnses = append(dnses, dns.(string))
	}
	ev.DNSServerNames = dnses
	ev.Location = d.Get("location").(string)

	subnets := []subnet{}
	for _, sub := range d.Get("subnet").(*schema.Set).List() {
		m := sub.(map[string]interface{})
		s := subnet{
			Name:          m["name"].(string),
			AddressPrefix: m["address_prefix"].(string),
			SecurityGroup: m["security_group"].(string),
		}
		subnets = append(subnets, s)
	}
	ev.Subnets = subnets

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
	fields["resource_group_name"] = ev.ResourceGroupName
	fields["address_space"] = ev.AddressSpace
	fields["dns_servers"] = ev.DNSServerNames
	fields["location"] = ev.Location
	fields["subnet"] = ev.mapSubnets()
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

// Maps the event subnets to a valid ResourceData List
func (ev *Event) mapSubnets() []map[string]interface{} {
	subnets := make([]map[string]interface{}, 0)
	for _, sub := range ev.Subnets {
		s := map[string]interface{}{}
		s["name"] = sub.Name
		s["address_prefix"] = sub.AddressPrefix
		s["security_group"] = sub.SecurityGroup
		subnets = append(subnets, s)
	}
	return subnets
}

// Error : will mark the event as errored
func (ev *Event) Error(err error) {
	ev.ErrorMessage = err.Error()
	ev.Body, err = json.Marshal(ev)
}
