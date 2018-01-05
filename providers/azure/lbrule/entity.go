/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package lbrule

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
)

// Event : This is the Ernest representation of an azure lb
type Event struct {
	event.Base
	ID                          string            `json:"id" diff:"-"`
	Name                        string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName           string            `json:"resource_group_name" validate:"required" diff:"-"`
	Loadbalancer                string            `json:"loadbalancer" diff:"-"`
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
	ErrorMessage                string            `json:"error,omitempty" diff:"-"`
	Components                  []json.RawMessage `json:"components" diff:"-"`
	CryptoKey                   string            `json:"-" diff:"-"`
	Validator                   *event.Validator  `json:"-" diff:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev event.Resource
	ev = &Event{CryptoKey: cryptoKey, Validator: val}
	body = []byte(strings.Replace(string(body), `"_component":"lb_rules"`, `"_component":"lb_rule"`, 1))
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_lb_rule", body, val, ev)
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
	if parts[7] != "loadbalancers" {
		return false
	}
	if parts[9] != "loadbalancingrules" {
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
	ev.ComponentID = "lb_rule::" + ev.Name
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.FrontendIPConfigurationName = d.Get("frontend_ip_configuration_name").(string)
	ev.Protocol = d.Get("protocol").(string)
	ev.FrontendPort = d.Get("frontend_port").(int)
	ev.BackendPort = d.Get("backend_port").(int)
	ev.BackendAddressPoolID = d.Get("backend_address_pool_id").(string)
	ev.ProbeID = d.Get("probe_id").(string)
	ev.EnableFloatingIP = d.Get("enable_floating_ip").(bool)
	ev.IdleTimeoutInMinutes = d.Get("idle_timeout_in_minutes").(int)
	ev.LoadDistribution = d.Get("load_distribution").(string)
	parts := strings.Split(ev.ID, "/")
	ev.LoadbalancerID = strings.Join(parts[0:9], "/")
	ev.Loadbalancer = parts[8]

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
	fields["loadbalancer_id"] = ev.LoadbalancerID
	fields["frontend_ip_configuration_name"] = ev.FrontendIPConfigurationName
	fields["protocol"] = ev.Protocol
	fields["frontend_port"] = ev.FrontendPort
	fields["backend_port"] = ev.BackendPort
	fields["backend_address_pool_id"] = ev.BackendAddressPoolID
	fields["probe_id"] = ev.ProbeID
	fields["enable_floating_ip"] = ev.EnableFloatingIP
	fields["idle_timeout_in_minutes"] = ev.IdleTimeoutInMinutes
	fields["load_distribution"] = ev.LoadDistribution
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
