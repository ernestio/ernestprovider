/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storageaccount

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
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
	ClientID               string            `json:"azure_client_id"`
	ClientSecret           string            `json:"azure_client_secret"`
	TenantID               string            `json:"azure_tenant_id"`
	SubscriptionID         string            `json:"azure_subscription_id"`
	Environment            string            `json:"environment"`
	ErrorMessage           string            `json:"error,omitempty"`
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

	return azure.New(subject, "azurerm_storage_account", body, val, ev)
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
	ev.Location = d.Get("location").(string)
	ev.NetworkSecurityGroup = d.Get("network_security_group_id").(string)
	ev.AccountType = d.Get("account_type").(string)
	ev.PrimaryLocation = d.Get("primary_location").(string)
	ev.SecondaryLocation = d.Get("secondary_location").(string)
	ev.PrimaryBlobEndpoint = d.Get("primary_blob_endpoint").(string)
	ev.SecondaryBlobEndpoint = d.Get("secondary_blob_endpoint").(string)
	ev.PrimaryQueueEndpoint = d.Get("primary_queue_endpoint").(string)
	ev.SecondaryQueueEndpoint = d.Get("secondary_queue_endpoint").(string)
	ev.PrimaryTableEndpoint = d.Get("primary_table_endpoint").(string)
	ev.SecondaryTableEndpoint = d.Get("secondary_table_endpoint").(string)
	ev.PrimaryFileEndpoint = d.Get("primary_file_endpoint").(string)
	ev.PrimaryAccessKey = d.Get("primary_access_key").(string)
	ev.SecondaryAccessKey = d.Get("secondary_access_key").(string)
	ev.Tags = d.Get("tags").(map[string]string)

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
			return err
		}
	}

	return nil
}
