/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package subnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/builtin/providers/azurerm"
	"github.com/hashicorp/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
)

// Event : This is the Ernest representation of an azure subnet
type Event struct {
	event.Base

	ID                   string   `json:"id"`
	Name                 string   `json:"name" validate:"required"`
	ResourceGroupName    string   `json:"resource_group_name" validate:"required"`
	VirtualNetworkName   string   `json:"virtual_network_name" validate:"required"`
	AddressPrefix        string   `json:"address_prefix"  validate:"required"`
	NetworkSecurityGroup string   `json:"network_security_group_id"`
	RouteTable           string   `json:"route_table_id"`
	IPConfigurations     []string `json:"ip_configurations"`

	ClientID       string `json:"azure_client_id"`
	ClientSecret   string `json:"azure_client_secret"`
	TenantID       string `json:"azure_tenant_id"`
	SubscriptionID string `json:"azure_subscription_id"`
	Environment    string `json:"environment"`

	Provider     *schema.Provider
	Component    *schema.Resource
	ResourceData *schema.ResourceData
	Schema       map[string]*schema.Schema
	ErrorMessage string           `json:"error,omitempty"`
	Subject      string           `json:"-"`
	Body         []byte           `json:"-"`
	CryptoKey    string           `json:"-"`
	Validator    *event.Validator `json:"-"`
}

// New : Constructor
func New(subject string, body []byte, cryptoKey string, val *event.Validator) (event.Event, error) {
	var err error
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey, Validator: val}
	n.Provider = azurerm.Provider().(*schema.Provider)
	n.Component = n.Provider.ResourcesMap["azurerm_subnet"]
	n.Schema = n.schema()
	n.Body = body
	n.Subject = subject
	n.CryptoKey = cryptoKey
	n.Validator = val
	if n.ResourceData, err = n.toResourceData(body); err != nil {
		n.Log("error", err.Error())
		return &n, err
	}

	return &n, nil
}

// Validate checks if all criteria are met
func (ev *Event) Validate() error {
	return ev.Validator.Validate(ev)
}

// Find : Find an object on azure
func (ev *Event) Find() error {
	return errors.New(ev.Subject + " not supported")
}

// Create : Creates a Virtual Network on Azure using terraform
// providers
func (ev *Event) Create() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Create(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error creating the requestd resource : %s", err)
		ev.Log("error", err.Error())
		return err
	}

	return nil
}

// Update : Updates an existing Virtual Network on Azure
// by using azurerm terraform provider resource
func (ev *Event) Update() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Update(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error creating the requestd resource : %s", err)
		ev.Log("error", err.Error())
		return err
	}

	return nil
}

// Get : Requests and loads the resource to Azure through azurerm
// terraform provider
func (ev *Event) Get() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Read(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error getting virtual network : %s", err)
		ev.Log("error", err.Error())
		return err
	}

	ev.toEvent()
	return nil
}

// Delete : Deletes the received resource from azure through
// azurerm terraform provider
func (ev *Event) Delete() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Delete(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error deleting the requested resource : %s", err)
		ev.Log("error", err.Error())
		return err
	}

	return nil
}

// GetBody : Gets the body for this event
func (ev *Event) GetBody() []byte {
	var err error
	if ev.Body, err = json.Marshal(ev); err != nil {
		log.Println(err.Error())
	}
	return ev.Body
}

// GetSubject : Gets the subject for this event
func (ev *Event) GetSubject() string {
	return ev.Subject
}

// Process : starts processing the current message
func (ev *Event) Process() (err error) {
	if err := json.Unmarshal(ev.Body, &ev); err != nil {
		ev.Error(err)
		return err
	}

	return nil
}

// Error : Will respond the current event with an error
func (ev *Event) Error(err error) {
	log.Printf("Error: %s", err.Error())
	ev.ErrorMessage = err.Error()

	ev.Body, err = json.Marshal(ev)
}

// Translates a ResourceData on a valid Ernest Event
func (ev *Event) toEvent() {
	ev.Name = ev.ResourceData.Get("name").(string)
	ev.ResourceGroupName = ev.ResourceData.Get("resource_group_name").(string)
	ev.VirtualNetworkName = ev.ResourceData.Get("virtual_network_name").(string)
	ev.AddressPrefix = ev.ResourceData.Get("address_prefix").(string)
	ev.NetworkSecurityGroup = ev.ResourceData.Get("network_security_group").(string)
	ev.RouteTable = ev.ResourceData.Get("route_table_id").(string)
	ev.IPConfigurations = ev.ResourceData.Get("ip_configurations").([]string)
}

// Translates the current event on a valid ResourceData
func (ev *Event) toResourceData(body []byte) (*schema.ResourceData, error) {
	var d schema.ResourceData
	d.SetSchema(ev.Schema)
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		ev.Log("error", err.Error())
		return nil, err
	}

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
			return nil, err
		}
		if err := d.Set(k, dec); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return nil, err
		}
	}

	fields := make(map[string]interface{})
	fields["name"] = ev.Name
	fields["resource_group_name"] = ev.ResourceGroupName
	fields["virtual_network_name"] = ev.VirtualNetworkName
	fields["address_prefix"] = ev.AddressPrefix
	fields["network_security_group"] = ev.NetworkSecurityGroup
	fields["route_table_id"] = ev.RouteTable
	fields["ip_configurations"] = ev.IPConfigurations
	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return nil, err
		}
	}

	return &d, nil
}

// Based on the Provider and Component schemas it calculates
// the necessary schema to be create a new ResourceData
func (ev *Event) schema() (sch map[string]*schema.Schema) {
	if ev.Schema != nil {
		return ev.Schema
	}
	a := ev.Provider.Schema
	b := ev.Component.Schema
	sch = a
	for k, v := range b {
		sch[k] = v
	}
	return sch
}

// Azure virtual network client
func (ev *Event) client() (*azurerm.ArmClient, error) {
	client, err := ev.Provider.ConfigureFunc(ev.ResourceData)
	if err != nil {
		err := fmt.Errorf("Can't connect to provider : %s", err)
		ev.Log("error", err.Error())
		return nil, err
	}
	c := client.(*azurerm.ArmClient)
	return c, nil
}
