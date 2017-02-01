/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package networkinterface

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"

	"github.com/Azure/azure-sdk-for-go/arm/network"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	event.Base
	ID                   string             `json:"id"`
	Name                 string             `json:"name" validate:"required"`
	ResourceGroupName    string             `json:"resource_group_name" validate:"required"`
	Location             string             `json:"location" validate:"required"`
	NetworkSecurityGroup string             `json:"network_security_group_id"`
	MacAddress           string             `json:"mac_address"`
	PrivateIPAddress     string             `json:"private_ip_address"`
	VirtualMachineID     string             `json:"virtual_machine_id"`
	IPConfigurations     []IPConfiguration  `json:"ip_configurations" validate:"required,dive"`
	DNSServers           []string           `json:"dns_servers" validate:"dive,ip"`
	InternalDNSNameLabel string             `json:"internal_dns_name_label"`
	AppliedDNSServers    []string           `json:"applied_dns_servers"`
	InternalFQDN         string             `json:"internal_fqdn"`
	EnableIPForwarding   bool               `json:"enable_ip_forwarding"`
	Tags                 map[string]*string `json:"tags"`

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

// IPConfiguration : ...
type IPConfiguration struct {
	Name                            string   `json:"name" validate:"required"`
	Subnet                          string   `json:"subnet_id" validate:"required"`
	PrivateIPAddress                string   `json:"private_ip_address"`
	PrivateIPAddressAllocation      string   `json:"private_ip_address_allocation" validate:"required"`
	PublicIPAddress                 string   `json:"public_ip_address_id"`
	LoadBalancerBackendAddressPools []string `json:"load_balancer_backend_address_pools_ids"`
	LoadBalancerInboundNatRules     []string `json:"load_balancer_inbound_nat_rules_ids"`
}

// New : Constructor
func New(subject string, body []byte, cryptoKey string, val *event.Validator) event.Event {
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey, Validator: val}

	return &n
}

// Azure client
func (ev *Event) client() *network.InterfacesClient {
	client, err := azure.Provider(ev.SubscriptionID, ev.ClientID, ev.ClientSecret, ev.TenantID, ev.Environment, ev.CryptoKey)
	if err != nil {
		panic(err)
	}

	return &client.IfaceClient
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
