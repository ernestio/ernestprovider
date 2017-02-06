/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package networkinterface

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
	ID                   string            `json:"id"`
	Name                 string            `json:"name" validate:"required"`
	ResourceGroupName    string            `json:"resource_group_name" validate:"required"`
	Location             string            `json:"location" validate:"required"`
	NetworkSecurityGroup string            `json:"network_security_group_id"`
	MacAddress           string            `json:"mac_address"`
	PrivateIPAddress     string            `json:"private_ip_address"`
	VirtualMachineID     string            `json:"virtual_machine_id"`
	IPConfigurations     []IPConfiguration `json:"ip_configurations" validate:"required,dive"`
	DNSServers           []string          `json:"dns_servers" validate:"dive,ip"`
	InternalDNSNameLabel string            `json:"internal_dns_name_label"`
	AppliedDNSServers    []string          `json:"applied_dns_servers"`
	InternalFQDN         string            `json:"internal_fqdn"`
	EnableIPForwarding   bool              `json:"enable_ip_forwarding"`
	Tags                 map[string]string `json:"tags"`

	Provider       *schema.Provider
	Component      *schema.Resource
	ResourceData   *schema.ResourceData
	Schema         map[string]*schema.Schema
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
func New(subject string, body []byte, cryptoKey string, val *event.Validator) (event.Event, error) {
	var err error
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey, Validator: val}
	n.Provider = azurerm.Provider().(*schema.Provider)
	n.Component = n.Provider.ResourcesMap["azurerm_network_interface"]
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
	ev.Location = ev.ResourceData.Get("location").(string)
	ev.NetworkSecurityGroup = ev.ResourceData.Get("network_security_group_id").(string)
	ev.MacAddress = ev.ResourceData.Get("mac_address").(string)
	ev.PrivateIPAddress = ev.ResourceData.Get("private_ip_address").(string)
	ev.VirtualMachineID = ev.ResourceData.Get("virtual_machine_id").(string)
	ev.IPConfigurations = ev.ResourceData.Get("ip_configurations").([]IPConfiguration)
	ev.DNSServers = ev.ResourceData.Get("dns_servers").([]string)
	ev.InternalDNSNameLabel = ev.ResourceData.Get("internal_dns_name_label").(string)
	ev.AppliedDNSServers = ev.ResourceData.Get("applied_dns_servers").([]string)
	ev.InternalFQDN = ev.ResourceData.Get("internal_fqdn").(string)
	ev.EnableIPForwarding = ev.ResourceData.Get("enable_ip_forwarding").(bool)
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
	fields["network_security_group_id"] = ev.NetworkSecurityGroup
	fields["mac_address"] = ev.MacAddress
	fields["private_ip_address"] = ev.PrivateIPAddress
	fields["virtual_machine_id"] = ev.VirtualMachineID
	fields["ip_configurations"] = ev.IPConfigurations
	fields["dns_servers"] = ev.DNSServers
	fields["internal_dns_name_label"] = ev.InternalDNSNameLabel
	fields["applied_dns_servers"] = ev.AppliedDNSServers
	fields["internal_fqdn"] = ev.InternalFQDN
	fields["enable_ip_forwarding"] = ev.EnableIPForwarding
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
