/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storageaccount

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	event.Base
	ID                     string            `json:"id" diff:"-"`
	Name                   string            `json:"name" validate:"required" diff:"-"`
	ResourceGroupName      string            `json:"resource_group_name" validate:"required" diff:"-"`
	Location               string            `json:"location" validate:"required" diff:"-"`
	AccountKind            string            `json:"account_kind" diff:"account_kind"`
	AccountType            string            `json:"account_type" validate:"required" diff:"account_type"`
	PrimaryLocation        string            `json:"primary_location" diff:"-"`
	SecondaryLocation      string            `json:"secondary_location" diff:"-"`
	PrimaryBlobEndpoint    string            `json:"primary_blob_endpoint" diff:"-"`
	SecondaryBlobEndpoint  string            `json:"secondary_blob_endpoint" diff:"-"`
	PrimaryQueueEndpoint   string            `json:"primary_queue_endpoint" diff:"-"`
	SecondaryQueueEndpoint string            `json:"secondary_queue_endpoint" diff:"-"`
	PrimaryTableEndpoint   string            `json:"primary_table_endpoint" diff:"-"`
	SecondaryTableEndpoint string            `json:"secondary_table_endpoint" diff:"-"`
	PrimaryFileEndpoint    string            `json:"primary_file_endpoint" diff:"-"`
	PrimaryAccessKey       string            `json:"primary_access_key" diff:"-"`
	SecondaryAccessKey     string            `json:"secondary_access_key" diff:"-"`
	EnableBlobEncryption   bool              `json:"enable_blob_encryption" diff:"enable_blob_encryption"`
	Tags                   map[string]string `json:"tags" diff:"tags"`
	ClientID               string            `json:"azure_client_id" diff:"-"`
	ClientSecret           string            `json:"azure_client_secret" diff:"-"`
	TenantID               string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID         string            `json:"azure_subscription_id" diff:"-"`
	Environment            string            `json:"environment" diff:"-"`
	ErrorMessage           string            `json:"error,omitempty" diff:"-"`
	Components             []json.RawMessage `json:"components" diff:"-"`
	CryptoKey              string            `json:"-" diff:"-"`
	Validator              *event.Validator  `json:"-" diff:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev event.Resource
	ev = &Event{CryptoKey: cryptoKey, Validator: val}
	body = []byte(strings.Replace(string(body), `"_component":"storage_accounts"`, `"_component":"storage_account"`, 1))
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_storage_account", body, val, ev)
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
	if parts[6] != "microsoft.storage" {
		return false
	}
	if parts[7] != "storageaccounts" {
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
	ev.ComponentID = "storage_account::" + ev.Name
	ev.Location = d.Get("location").(string)
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.Location = d.Get("location").(string)
	ev.AccountKind = d.Get("account_kind").(string)
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
	fields["account_kind"] = ev.AccountKind
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
