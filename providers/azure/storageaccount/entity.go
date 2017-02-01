/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storageaccount

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"

	"github.com/Azure/azure-sdk-for-go/arm/storage"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	event.Base
	ID                     string             `json:"id"`
	Name                   string             `json:"name" validate:"required"`
	ResourceGroupName      string             `json:"resource_group_name" validate:"required"`
	Location               string             `json:"location" validate:"required"`
	NetworkSecurityGroup   string             `json:"network_security_group_id"`
	AccountType            string             `json:"account_type" validate:"required"`
	PrimaryLocation        string             `json:"primary_location"`
	SecondaryLocation      string             `json:"secondary_location"`
	PrimaryBlobEndpoint    string             `json:"primary_blob_endpoint"`
	SecondaryBlobEndpoint  string             `json:"secondary_blob_endpoint"`
	PrimaryQueueEndpoint   string             `json:"primary_queue_endpoint"`
	SecondaryQueueEndpoint string             `json:"secondary_queue_endpoint"`
	PrimaryTableEndpoint   string             `json:"primary_table_endpoint"`
	SecondaryTableEndpoint string             `json:"secondary_table_endpoint"`
	PrimaryFileEndpoint    string             `json:"primary_file_endpoint"`
	PrimaryAccessKey       string             `json:"primary_access_key"`
	SecondaryAccessKey     string             `json:"secondary_access_key"`
	Tags                   map[string]*string `json:"tags"`

	ClientID       string `json:"azure_client_id"`
	ClientSecret   string `json:"azure_client_secret"`
	TenantID       string `json:"azure_tenant_id"`
	SubscriptionID string `json:"azure_subscription_id"`
	Environment    string `json:"environment"`

	ErrorMessage string           `json:"error,omitempty"`
	Subject      string           `json:"-"`
	Body         []byte           `json:"-"`
	CryptoKey    string           `json:"-"`
	Validator    *event.Validator `json:"-"`
}

// New : Constructor
func New(subject string, body []byte, cryptoKey string, val *event.Validator) event.Event {
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey, Validator: val}

	return &n
}

// Azure client
func (ev *Event) client() *storage.AccountsClient {
	client, err := azure.Provider(ev.SubscriptionID, ev.ClientID, ev.ClientSecret, ev.TenantID, ev.Environment, ev.CryptoKey)
	if err != nil {
		panic(err)
	}

	return &client.StorageServiceClient
}

// Validate checks if all criteria are met
func (ev *Event) Validate() error {
	return ev.Validator.Validate(ev)
}

// Find : Find an object on azure
func (ev *Event) Find() error {
	return errors.New(ev.Subject + " not supported")
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
