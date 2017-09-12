/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manageddisk

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
	ID                 string            `json:"id"`
	Name               string            `json:"name" validate:"required"`
	ResourceGroupName  string            `json:"resource_group_name" validate:"required"`
	Location           string            `json:"location"`
	StorageAccountType string            `json:"storage_account_type"`
	CreateOption       string            `json:"create_option"`
	SourceURI          string            `json:"source_uri"`
	SourceResourceID   string            `json:"source_resource_id"`
	OSType             string            `json:"os_type"`
	DiskSizeGB         int32             `json:"disk_size_gb"`
	Tags               map[string]string `json:"tags"`
	ClientID           string            `json:"azure_client_id"`
	ClientSecret       string            `json:"azure_client_secret"`
	TenantID           string            `json:"azure_tenant_id"`
	SubscriptionID     string            `json:"azure_subscription_id"`
	Environment        string            `json:"environment"`
	ErrorMessage       string            `json:"error,omitempty"`
	Components         []json.RawMessage `json:"components"`
	CryptoKey          string            `json:"-"`
	Validator          *event.Validator  `json:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev event.Resource
	ev = &Event{CryptoKey: cryptoKey, Validator: val}
	body = []byte(strings.Replace(string(body), `"_component":"managed_disks"`, `"_component":"managed_disk"`, 1))
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_managed_disk", body, val, ev)
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
	if parts[6] != "microsoft.compute" {
		return false
	}
	if parts[7] != "disks" {
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
	ev.ComponentID = "managed_disk::" + ev.Name
	if ev.ResourceGroupName == "" {
		ev.ResourceGroupName = d.Get("resource_group_name").(string)
	}
	ev.Location = d.Get("location").(string)
	ev.StorageAccountType = d.Get("storage_account_type").(string)
	ev.CreateOption = d.Get("create_option").(string)
	ev.SourceURI = d.Get("source_uri").(string)
	ev.SourceResourceID = d.Get("source_resource_id").(string)
	ev.OSType = d.Get("os_type").(string)
	ev.DiskSizeGB = int32(d.Get("disk_size_gb").(int))
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
	fields["storage_account_type"] = ev.StorageAccountType
	fields["create_option"] = ev.CreateOption
	if ev.SourceURI != "" {
		fields["source_uri"] = ev.SourceURI
	}
	if ev.SourceResourceID != "" {
		fields["source_resource_id"] = ev.SourceResourceID
	}
	fields["os_type"] = ev.OSType
	fields["disk_size_gb"] = ev.DiskSizeGB
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
