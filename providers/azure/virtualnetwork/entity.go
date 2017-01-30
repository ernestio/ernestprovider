/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualnetwork

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"

	"github.com/Azure/azure-sdk-for-go/arm/network"
)

// Event : This is the Ernest representation of an azure virtual
//         network
type Event struct {
	event.Base
	ID                string             `json:"id"`
	Name              string             `json:"name" validate:"required"`
	AddressSpace      []string           `json:"address_space" validate:"min=1"`
	DNSServerNames    []string           `json:"dns_server_names" validate:"dive,ip"`
	Subnets           []subnet           `json:"subnets" validate:"min=1,dive"`
	Location          string             `json:"location"`
	ClientID          string             `json:"azure_client_id"`
	ClientSecret      string             `json:"azure_client_secret"`
	TenantID          string             `json:"azure_tenant_id"`
	SubscriptionID    string             `json:"azure_subscription_id"`
	Environment       string             `json:"environment"`
	ResourceGroupName string             `json:"resource_group_name"`
	Tags              map[string]*string `json:"tags"`

	ErrorMessage string           `json:"error,omitempty"`
	Subject      string           `json:"-"`
	Body         []byte           `json:"-"`
	CryptoKey    string           `json:"-"`
	Validator    *event.Validator `json:"-"`
}

type subnet struct {
	Name          string `json:"name" validate:"required"`
	AddressPrefix string `json:"address_prefix" validate:"required,cidr"`
	SecurityGroup string `json:"security_group"`
}

// New : Constructor
func New(subject string, body []byte, cryptoKey string, val *event.Validator) event.Event {
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey, Validator: val}

	return &n
}

// Azure virtual network client
func (ev *Event) client() *network.VirtualNetworksClient {
	client, err := azure.Provider(ev.SubscriptionID, ev.ClientID, ev.ClientSecret, ev.TenantID, ev.Environment, ev.CryptoKey)
	if err != nil {
		panic(err)
	}

	return &client.VnetClient
}

// Validate checks if all criteria are met
func (ev *Event) Validate() error {
	return ev.Validator.Validate(ev)
}

// Find : Find an object on azure
func (ev *Event) Find() error {
	return errors.New(ev.Subject + " not supported")
}

// Update : Updates a nat object on azure
func (ev *Event) Update() error {
	return ev.Create()
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

func (ev *Event) getVirtualNetworkProperties() *network.VirtualNetworkPropertiesFormat {
	// first; get address space prefixes:
	prefixes := ev.AddressSpace

	// then; the dns servers:
	dnses := ev.DNSServerNames

	// then; the subnets:
	subnets := []network.Subnet{}
	for _, subnet := range ev.Subnets {
		name := subnet.Name
		prefix := subnet.AddressPrefix
		secGroup := subnet.SecurityGroup

		var subnetObj network.Subnet
		subnetObj.Name = &name
		subnetObj.SubnetPropertiesFormat = &network.SubnetPropertiesFormat{}
		subnetObj.SubnetPropertiesFormat.AddressPrefix = &prefix

		if secGroup != "" {
			subnetObj.SubnetPropertiesFormat.NetworkSecurityGroup = &network.SecurityGroup{
				ID: &secGroup,
			}
		}

		subnets = append(subnets, subnetObj)
	}

	// finally; return the struct:
	return &network.VirtualNetworkPropertiesFormat{
		AddressSpace: &network.AddressSpace{
			AddressPrefixes: &prefixes,
		},
		DhcpOptions: &network.DhcpOptions{
			DNSServers: &dnses,
		},
		Subnets: &subnets,
	}
}
