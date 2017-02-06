/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storageaccount

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

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	event.Base

	ID                     string            `json:"id"`
	Name                   string            `json:"name" validate:"required"`
	ResourceGroupName      string            `json:"resource_group_name" validate:"required"`
	Location               string            `json:"location" validate:"required"`
	NetworkSecurityGroup   string            `json:"network_security_group_id"`
	AccountType            string            `json:"account_type" validate:"required"`
	PrimaryLocation        string            `json:"primary_location"`
	SecondaryLocation      string            `json:"secondary_location"`
	PrimaryBlobEndpoint    string            `json:"primary_blob_endpoint"`
	SecondaryBlobEndpoint  string            `json:"secondary_blob_endpoint"`
	PrimaryQueueEndpoint   string            `json:"primary_queue_endpoint"`
	SecondaryQueueEndpoint string            `json:"secondary_queue_endpoint"`
	PrimaryTableEndpoint   string            `json:"primary_table_endpoint"`
	SecondaryTableEndpoint string            `json:"secondary_table_endpoint"`
	PrimaryFileEndpoint    string            `json:"primary_file_endpoint"`
	PrimaryAccessKey       string            `json:"primary_access_key"`
	SecondaryAccessKey     string            `json:"secondary_access_key"`
	Tags                   map[string]string `json:"tags"`

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
	n.Component = n.Provider.ResourcesMap["azurerm_storage_account"]
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
	ev.Location = ev.ResourceData.Get("location").(string)

	ev.ResourceGroupName = ev.ResourceData.Get("resource_group_name").(string)
	ev.Location = ev.ResourceData.Get("location").(string)
	ev.NetworkSecurityGroup = ev.ResourceData.Get("network_security_group_id").(string)
	ev.AccountType = ev.ResourceData.Get("account_type").(string)
	ev.PrimaryLocation = ev.ResourceData.Get("primary_location").(string)
	ev.SecondaryLocation = ev.ResourceData.Get("secondary_location").(string)
	ev.PrimaryBlobEndpoint = ev.ResourceData.Get("primary_blob_endpoint").(string)
	ev.SecondaryBlobEndpoint = ev.ResourceData.Get("secondary_blob_endpoint").(string)
	ev.PrimaryQueueEndpoint = ev.ResourceData.Get("primary_queue_endpoint").(string)
	ev.SecondaryQueueEndpoint = ev.ResourceData.Get("secondary_queue_endpoint").(string)
	ev.PrimaryTableEndpoint = ev.ResourceData.Get("primary_table_endpoint").(string)
	ev.SecondaryTableEndpoint = ev.ResourceData.Get("secondary_table_endpoint").(string)
	ev.PrimaryFileEndpoint = ev.ResourceData.Get("primary_file_endpoint").(string)
	ev.PrimaryAccessKey = ev.ResourceData.Get("primary_access_key").(string)
	ev.SecondaryAccessKey = ev.ResourceData.Get("secondary_access_key").(string)
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
	fields["resource_group_name"] = ev.ResourceGroupName
	fields["location"] = ev.Location
	fields["network_security_group"] = ev.NetworkSecurityGroup
	fields["account_type"] = ev.AccountType
	fields["primary_location"] = ev.PrimaryLocation
	fields["secondary_location"] = ev.SecondaryLocation
	fields["primary_blob_endpoint"] = ev.PrimaryBlobEndpoint
	fields["secondary_blob_endpoint"] = ev.SecondaryBlobEndpoint
	fields["primary_queue_endpoint"] = ev.PrimaryQueueEndpoint
	fields["secondary_queue_endpoint"] = ev.SecondaryQueueEndpoint
	fields["primary_table_endpoint"] = ev.PrimaryTableEndpoint
	fields["secondary_table_endpoint"] = ev.SecondaryTableEndpoint
	fields["primary_file_endpoint"] = ev.PrimaryFileEndpoint
	fields["primary_access_key"] = ev.PrimaryAccessKey
	fields["secondary_access_key"] = ev.SecondaryAccessKey
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
