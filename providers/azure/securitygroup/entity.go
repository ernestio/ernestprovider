/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package securitygroup

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
	ID                string            `json:"id"`
	Name              string            `json:"name" validate:"required"`
	Location          string            `json:"location" validate:"required"`
	ResourceGroupName string            `json:"resource_group_name" validate:"required"`
	SecurityRules     []SecurityRule    `json:"security_rules"`
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

// SecurityRule ...
type SecurityRule struct {
	Name                     string `json:"name"`
	Description              string `json:"description"`
	Protocol                 string `json:"protocol"`
	SourcePort               string `json:"source_port_range"`
	DestinationPortRange     string `json:"destination_port_range"`
	SourceAddressPrefix      string `json:"source_address_prefix"`
	DestinationAddressPrefix string `json:"destination_address_prefix"`
	Access                   string `json:"access"`
	Priority                 int    `json:"priority"`
	Direction                string `json:"direction"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev azure.Resource
	ev = &Event{CryptoKey: cryptoKey}
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_network_security_group", body, val, ev)
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
	if parts[7] != "networksecuritygroups" {
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
	ev.Location = d.Get("location").(string)
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	rules := []SecurityRule{}
	for _, sec := range d.Get("security_rule").(*schema.Set).List() {
		sg := sec.(map[string]interface{})
		rules = append(rules, SecurityRule{
			Name:                     sg["name"].(string),
			Description:              sg["description"].(string),
			Protocol:                 sg["protocol"].(string),
			SourcePort:               sg["source_port_range"].(string),
			DestinationPortRange:     sg["destination_port_range"].(string),
			SourceAddressPrefix:      sg["source_address_prefix"].(string),
			DestinationAddressPrefix: sg["destination_address_prefix"].(string),
			Access:    sg["access"].(string),
			Priority:  sg["priority"].(int),
			Direction: sg["direction"].(string),
		})
	}
	ev.SecurityRules = rules

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
	fields["location"] = ev.Location

	rules := make([]map[string]interface{}, 0)
	for _, sec := range ev.SecurityRules {
		sr := map[string]interface{}{}
		sr["name"] = sec.Name
		sr["description"] = sec.Description
		sr["protocol"] = sec.Protocol
		sr["source_port_range"] = sec.SourcePort
		sr["destination_port_range"] = sec.DestinationPortRange
		sr["source_address_prefix"] = sec.SourceAddressPrefix
		sr["destination_address_prefix"] = sec.DestinationAddressPrefix
		sr["access"] = sec.Access
		sr["priority"] = sec.Priority
		sr["direction"] = sec.Direction
		rules = append(rules, sr)
	}
	fields["security_rule"] = rules
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
