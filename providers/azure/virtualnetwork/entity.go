/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualnetwork

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/builtin/providers/azurerm"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
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
	ClientID          string            `json:"azure_client_id"`
	ClientSecret      string            `json:"azure_client_secret"`
	TenantID          string            `json:"azure_tenant_id"`
	SubscriptionID    string            `json:"azure_subscription_id"`
	Environment       string            `json:"environment"`
	ResourceGroupName string            `json:"resource_group_name"`
	Tags              map[string]string `json:"tags"`

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

type subnet struct {
	Name          string `json:"name"`
	AddressPrefix string `json:"address_prefix"`
	SecurityGroup string `json:"security_group"`
}

// New : Virtual Network constructor returning an Event object
// and an error in case something fails
func New(subject string, body []byte, cryptoKey string, val *event.Validator) (event.Event, error) {
	var err error
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey, Validator: val}
	n.Provider = azurerm.Provider().(*schema.Provider)
	n.Component = n.Provider.ResourcesMap["azurerm_virtual_network"]
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

// Validate : checks if all event fields are valid, it
// responds with an error in case something is wrong
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
	ev.AddressSpace = ev.ResourceData.Get("address_space").([]string)
	ev.DNSServerNames = ev.ResourceData.Get("dns_servers").([]string)
	ev.Location = ev.ResourceData.Get("location").(string)
	ev.Subnets = ev.ResourceData.Get("subnet").([]subnet)
	ev.Tags = ev.ResourceData.Get("tags").(map[string]string)
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
	fields["address_space"] = ev.AddressSpace
	fields["dns_servers"] = ev.DNSServerNames
	fields["location"] = ev.Location
	fields["subnet"] = ev.mapSubnets()
	fields["tags"] = ev.Tags
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

// Maps the event subnets to a valid ResourceData List
func (ev *Event) mapSubnets() *schema.Set {
	list := &schema.Set{
		F: resourceAzureSubnetHash,
	}
	for _, sub := range ev.Subnets {
		s := map[string]interface{}{}
		s["name"] = sub.Name
		s["address_prefix"] = sub.AddressPrefix
		list.Add(s)
	}
	return list
}

func resourceAzureSubnetHash(v interface{}) int {
	m := v.(map[string]interface{})
	subnet := m["name"].(string) + m["address_prefix"].(string)
	if securityGroup, present := m["security_group"]; present {
		subnet = subnet + securityGroup.(string)
	}

	return hashcode.String(subnet)
}
