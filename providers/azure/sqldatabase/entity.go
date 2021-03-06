/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package sqldatabase

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
	types "github.com/ernestio/ernestprovider/types/azure/sqldatabase"
	"github.com/ernestio/ernestprovider/validator"
)

// Event : This is the Ernest representation of an azure subnet
type Event struct {
	event.Base
	types.Event
	ErrorMessage string               `json:"error,omitempty" diff:"-"`
	CryptoKey    string               `json:"-" diff:"-"`
	Validator    *validator.Validator `json:"-" diff:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *validator.Validator) (event.Event, error) {
	var ev event.Resource
	ev = &Event{CryptoKey: cryptoKey, Validator: val}
	body = []byte(strings.Replace(string(body), `"_component":"sql_databases"`, `"_component":"sql_database"`, 1))
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_sql_database", body, val, ev)
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
	if parts[9] != "databases" {
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

	parts := strings.Split(ev.ID, "/")

	ev.Name = d.Get("name").(string)
	ev.ComponentID = "sql_database::" + ev.Name
	ev.Location = d.Get("location").(string)
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.ServerName = parts[len(parts)-3]
	ev.CreateMode = d.Get("create_mode").(string)
	ev.SourceDatabaseID = d.Get("source_database_id").(string)
	ev.RestorePointInTime = d.Get("restore_point_in_time").(string)
	ev.Edition = d.Get("edition").(string)
	ev.Collation = d.Get("collation").(string)
	ev.MaxSizeBytes = d.Get("max_size_bytes").(string)
	ev.RequestedServiceObjectiveID = d.Get("requested_service_objective_id").(string)
	ev.RequestedServiceObjectiveName = d.Get("requested_service_objective_name").(string)
	ev.SourceDatabaseDeletionData = d.Get("source_database_deletion_date").(string)
	ev.ElasticPoolName = d.Get("elastic_pool_name").(string)
	ev.Encryption = d.Get("encryption").(string)
	ev.CreationDate = d.Get("creation_date").(string)
	ev.DefaultSecondaryLocation = d.Get("default_secondary_location").(string)

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
	fields["server_name"] = ev.ServerName
	fields["create_mode"] = ev.CreateMode
	if fields["create_mode"] == "" {
		fields["create_mode"] = "Default"
	}
	fields["source_database_id"] = ev.SourceDatabaseID
	fields["restore_point_in_time"] = ev.RestorePointInTime
	fields["edition"] = ev.Edition
	fields["collation"] = ev.Collation
	fields["max_size_bytes"] = ev.MaxSizeBytes
	fields["requested_service_objective_id"] = ev.RequestedServiceObjectiveID
	fields["requested_service_objective_name"] = ev.RequestedServiceObjectiveName
	fields["source_database_deletion_date"] = ev.SourceDatabaseDeletionData
	fields["elastic_pool_name"] = ev.ElasticPoolName
	fields["encryption"] = ev.Encryption
	fields["creation_date"] = ev.CreationDate
	fields["default_secondary_location"] = ev.DefaultSecondaryLocation
	fields["tags"] = ev.Tags
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
