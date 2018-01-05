/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package sqlfirewallrule

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
	ID                string            `json:"id" diff:"-"`
	Name              string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName string            `json:"resource_group_name" validate:"required" diff:"-"`
	ServerName        string            `json:"server_name" validate:"required" diff:"-"`
	StartIPAddress    string            `json:"start_ip_address" validate:"required" diff:"start_ip_address"`
	EndIPAddress      string            `json:"end_ip_address" validate:"required" diff:"end_ip_address"`
	ClientID          string            `json:"azure_client_id" diff:"-"`
	ClientSecret      string            `json:"azure_client_secret" diff:"-"`
	TenantID          string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string            `json:"azure_subscription_id" diff:"-"`
	Environment       string            `json:"environment" diff:"-"`
	ErrorMessage      string            `json:"error,omitempty" diff:"-"`
	Components        []json.RawMessage `json:"components" diff:"-"`
	CryptoKey         string            `json:"-" diff:"-"`
	Validator         *event.Validator  `json:"-" diff:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev event.Resource
	ev = &Event{CryptoKey: cryptoKey, Validator: val}
	body = []byte(strings.Replace(string(body), `"_component":"sql_firewall_rules"`, `"_component":"sql_firewall_rule"`, 1))
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_sql_firewall_rule", body, val, ev)
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
	if parts[6] != "microsoft.sql" {
		return false
	}
	if parts[7] != "servers" {
		return false
	}
	if parts[9] != "firewallrules" {
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
	ev.ID = d.Id()
	ev.Name = d.Get("name").(string)
	ev.ComponentID = "sql_firewall_rule::" + ev.Name
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.ServerName = d.Get("server_name").(string)
	ev.StartIPAddress = d.Get("start_ip_address").(string)
	ev.EndIPAddress = d.Get("end_ip_address").(string)

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
	fields["server_name"] = ev.ServerName
	fields["start_ip_address"] = ev.StartIPAddress
	fields["end_ip_address"] = ev.EndIPAddress
	for k, v := range fields {
		if k != "tags" {
			println(k + " -> " + v.(string))
		}
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
